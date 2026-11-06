//go:build fyne

// Package fyneui is the Fyne front end. It renders the headless desktop model
// and owns no rules of its own; the default binary does not build it, so the
// cgo-free release is unchanged.
package fyneui

import (
	"fmt"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Editor applies an edit, a creation or a deletion and returns the refreshed
// set and history.
type Editor interface {
	Update(key string, changes catalog.Changes) (catalog.Catalog, timeline.Log, error)
	Create(command catalog.Command) (catalog.Catalog, timeline.Log, error)
	Delete(key string) (catalog.Catalog, timeline.Log, error)
}

// Reverter restores the command set an older history entry recorded and
// returns the refreshed set and history.
//
// It is declared here, at the point of use: the window confirms the restore
// itself, so it needs the state change and not the use case's prompt.
type Reverter interface {
	Revert(steps int) (catalog.Catalog, timeline.Log, error)
}

// Deps is what the window needs from the command layer. Every field is a port
// declared in this package and implemented in cmd, so the UI depends on
// behaviour, not on internal/session.
type Deps struct {
	Exec     desktop.Executor
	Editor   Editor
	Reverter Reverter
	Updater  Updater
	Tools    desktop.Tools
	Groups   desktop.Groups
	Setups   desktop.Setups
	Window   WindowStore
}

// Run opens the Fyne window and blocks until it closes.
func Run(commands catalog.Catalog, history timeline.Log, deps Deps) {
	model := desktop.NewModel(commands, history, deps.Exec)

	a := app.New()
	a.SetIcon(logoResource())

	window := a.NewWindow("crower")
	front := newUI(model, window, deps)
	window.SetContent(front.content())
	window.Resize(initialSize(deps.Window))

	// Fyne offers a window-wide drop callback only, so a drop has one documented
	// destination: the copy/move sources.
	window.SetOnDropped(front.handleDrop)

	// Remember the size the user settled on, so the next run opens the same way.
	window.SetOnClosed(func() {
		if deps.Window == nil {
			return
		}

		if size, ok := recordedSize(window.Canvas().Size()); ok {
			_ = deps.Window.Save(size)
		}
	})

	// The lookup runs after the content exists and never blocks the window:
	// opening crower is not a reason to wait for GitHub.
	front.checkForUpdates()

	window.ShowAndRun()
}

// ui is the widget set the tabs share. Reloading after an edit has to refresh
// both the command list and the history table, so both live here instead of
// being captured by one tab.
type ui struct {
	model    *desktop.Model
	window   fyne.Window
	editor   Editor
	reverter Reverter
	updater  Updater
	tools    desktop.Tools
	groups   desktop.Groups
	setups   desktop.Setups

	activeSetup  string
	setupSelect  *widget.Select
	setupChoices []setupOption

	groupIcons   map[string]string
	groupOptions []groupOption

	// tabs is the AppTabs itself: a drop needs to open the Tools tab, and the
	// Commands tab changes the setup selector's labels, so both need it.
	tabs *container.AppTabs

	// commandsRoot is the Commands tab's border. Showing the details panel needs
	// a layout pass from it: Fyne's layouts skip invisible children, so the panel
	// has no size until then.
	commandsRoot *fyne.Container

	list         *widget.List
	detail       *widget.Label
	detailFields *fyne.Container
	detailBody   *fyne.Container
	detailToggle *widget.Button
	detailShown  bool
	paramsEntry  *widget.Entry
	filterEntry  *widget.Entry
	groupFilter  *widget.Select

	historyScroll   *container.Scroll
	historyPage     *fyne.Container
	historyList     *fyne.Container
	historyEmpty    *widget.Label
	historyLoadMore *widget.Button
	historyCards    []*historyCard
	historyOpen     *historyCard
	historyTotal    int
	historyLimit    int

	// tools tab: one result square per card, plus the sources a drop fills
	sourcesEntry   *widget.Entry
	transferResult *toolResult
	searchResult   *toolResult
	extractResult  *toolResult
	openResult     *toolResult

	// about tab
	updateStatus  *widget.Label
	checkButton   *widget.Button
	installButton *widget.Button
	restoreButton *widget.Button
	restartButton *widget.Button
	aboutActions  *fyne.Container
	aboutBody     *fyne.Container
	aboutUpdates  fyne.CanvasObject
	latest        string
}

func newUI(model *desktop.Model, window fyne.Window, deps Deps) *ui {
	front := &ui{
		model:      model,
		window:     window,
		editor:     deps.Editor,
		reverter:   deps.Reverter,
		updater:    deps.Updater,
		tools:      deps.Tools,
		groups:     deps.Groups,
		setups:     deps.Setups,
		groupIcons: map[string]string{},
	}

	if deps.Setups != nil {
		front.activeSetup = deps.Setups.Current()
	}

	return front
}

// NewContent builds the window content over the headless model.
func NewContent(model *desktop.Model, window fyne.Window, deps Deps) fyne.CanvasObject {
	return newUI(model, window, deps).content()
}

// content builds the tabs. It is network-free: only Run starts the update
// lookup, so a test can build the window without reaching GitHub.
func (u *ui) content() fyne.CanvasObject {
	u.loadGroupIcons()

	u.tabs = container.NewAppTabs(
		container.NewTabItem("Commands", u.commandsTab()),
		container.NewTabItem("History", u.historyTab()),
		container.NewTabItem("Tools", u.toolsTab()),
		container.NewTabItem("About", u.aboutTab()),
	)

	// The table is built from the model, so a command edit changes it. Refresh
	// when the tab is opened instead of refreshing a hidden table on every
	// keystroke.
	u.tabs.OnSelected = func(item *container.TabItem) {
		if item.Text == "History" {
			u.refreshHistory()
		}
	}

	return u.tabs
}

// commandsTab is the command list with its filter, the per-row actions and a
// detail panel that opens below the list.
func (u *ui) commandsTab() fyne.CanvasObject {
	u.detail = widget.NewLabel("")
	u.detail.Wrapping = fyne.TextWrapWord

	u.detailFields = container.NewGridWithColumns(2)

	u.paramsEntry = newScrollFreeEntry()
	u.paramsEntry.SetPlaceHolder("Parameters (space-separated)")
	u.paramsEntry.OnSubmitted = func(string) { u.runFromDetail() }

	runDetail := widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), u.runFromDetail)

	// The body is hidden rather than built on demand: Fyne's layouts skip an
	// invisible child, so collapsing it gives the height back to the list.
	u.detailBody = container.NewVBox(
		u.detail,
		u.detailFields,
		container.NewBorder(nil, nil, nil, runDetail, u.paramsEntry),
	)
	u.detailBody.Hide()

	u.detailToggle = widget.NewButtonWithIcon("Details", theme.MenuDropDownIcon(), func() {
		u.setDetailVisible(!u.detailShown)
	})
	u.detailToggle.Importance = widget.LowImportance

	u.list = widget.NewList(
		func() int { return len(u.model.FilteredKeys()) },
		func() fyne.CanvasObject { return newCommandRow() },
		func(id widget.ListItemID, object fyne.CanvasObject) { u.updateRow(id, object) },
	)

	u.list.OnSelected = func(id widget.ListItemID) { u.selectRow(id) }
	u.list.OnUnselected = func(widget.ListItemID) { u.clearDetail() }

	filter := newScrollFreeEntry()
	u.filterEntry = filter
	filter.SetPlaceHolder("Filter commands")
	filter.OnChanged = func(string) {
		u.model.Filter = filter.Text
		u.model.Index = 0
		u.list.Refresh()
		u.list.UnselectAll()
		u.clearDetail()
	}

	u.groupFilter = widget.NewSelect([]string{allGroupsLabel}, u.selectGroup)
	u.groupFilter.SetSelected(allGroupsLabel)

	groupsButton := widget.NewButtonWithIcon("Groups…", theme.GridIcon(), u.showGroups)

	labels, choices := u.setupChoicesList()
	u.setupChoices = choices

	u.setupSelect = widget.NewSelect(labels, func(label string) {
		name := u.setupName(label)
		if name == "" || name == u.activeSetup {
			return
		}

		u.switchSetup(name)
	})
	u.setupSelect.SetSelected(u.setupLabel(u.activeSetup))

	controls := container.NewHBox(widget.NewLabel("Setup"), u.setupSelect)

	if u.setups != nil {
		removeSetup := widget.NewButtonWithIcon("Delete…", theme.DeleteIcon(), u.showDeleteSetup)
		removeSetup.Importance = widget.DangerImportance

		controls.Add(container.NewHBox(
			widget.NewButtonWithIcon("New…", theme.ContentAddIcon(), u.showNewSetup),
			removeSetup,
		))
	}

	controls.Add(widget.NewLabel("Group"))
	controls.Add(u.groupFilter)
	controls.Add(groupsButton)

	toolbar := container.NewHBox(
		widget.NewButtonWithIcon("New", theme.ContentAddIcon(), u.showCreate),
		widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), u.runSelected),
	)

	detail := container.NewVBox(u.detailToggle, u.detailBody)

	u.commandsRoot = container.NewBorder(
		container.NewVBox(filter, controls),
		container.NewVBox(detail, toolbar),
		nil,
		nil,
		u.list,
	)

	// A left margin keeps the filter, controls and row names off the window
	// edge. commandsRoot stays the inner border so showing the detail panel
	// still re-runs that layout, not the wrapper's.
	return container.NewBorder(nil, nil, gap(theme.Padding()*2, 0), nil, u.commandsRoot)
}

// allGroupsLabel is the group filter's "no filter" entry.
const allGroupsLabel = "All groups"

// groupOption maps a filter label to the group it selects; the label carries
// the group's icon, so it is not the same string as the name.
type groupOption struct {
	label string
	name  string
}

// loadGroupIcons reads the group icons once. A file that cannot be read leaves
// the filter with names only: an icon is decoration, not a reason the window
// fails to open.
func (u *ui) loadGroupIcons() {
	if u.groups == nil {
		return
	}

	icons, err := u.groups.Icons()
	if err != nil {
		return
	}

	u.groupIcons = icons
}

// selectGroup is the filter's change handler: it maps the label back to the
// group name and narrows the list.
func (u *ui) selectGroup(label string) {
	group := ""

	for _, option := range u.groupOptions {
		if option.label == label {
			group = option.name
			break
		}
	}

	u.model.Group = group
	u.model.Index = 0
	u.list.Refresh()
	u.list.UnselectAll()
	u.clearDetail()
}

// refreshGroups rebuilds the filter's options after an edit, a removal or a
// group icon change. A group with no commands left disappears, and the filter
// falls back to all.
func (u *ui) refreshGroups() {
	if u.groupFilter == nil {
		return
	}

	names := u.model.Groups()
	options := make([]string, 0, len(names)+1)
	u.groupOptions = make([]groupOption, 0, len(names)+1)

	options = append(options, allGroupsLabel)
	u.groupOptions = append(u.groupOptions, groupOption{label: allGroupsLabel})

	for _, name := range names {
		label := groupLabel(u.groupIcons[name], name)
		options = append(options, label)
		u.groupOptions = append(u.groupOptions, groupOption{label: label, name: name})
	}

	u.groupFilter.Options = options
	u.groupFilter.Refresh()

	selected := allGroupsLabel

	for _, option := range u.groupOptions {
		if option.name == u.model.Group && option.name != "" {
			selected = option.label
		}
	}

	if u.model.Group != "" && selected == allGroupsLabel {
		u.model.Group = ""
	}

	u.groupFilter.SetSelected(selected)
}

// commandRow is one list row: the command name on the left, its actions on the
// right. widget.List recycles templates, so updateRow re-points every button at
// the row's own key on each update.
//
// It is a widget holding a container, not a container subtype: Fyne walks a
// canvas by concrete type, descending into *fyne.Container and fyne.Widget only,
// so a struct that embeds a container is opaque and its children never paint
// (#406).
type commandRow struct {
	widget.BaseWidget

	name   *widget.Label
	run    *widget.Button
	edit   *widget.Button
	remove *widget.Button
	body   *fyne.Container
}

func newCommandRow() *commandRow {
	row := &commandRow{
		name:   widget.NewLabel(""),
		run:    widget.NewButtonWithIcon("", theme.MediaPlayIcon(), nil),
		edit:   widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), nil),
		remove: widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
	}

	row.run.Importance = widget.LowImportance
	row.edit.Importance = widget.LowImportance
	row.remove.Importance = widget.DangerImportance

	row.body = container.NewBorder(
		nil,
		nil,
		nil,
		container.NewHBox(row.run, row.edit, row.remove, gap(scrollGutter(), 0)),
		row.name,
	)

	row.ExtendBaseWidget(row)

	return row
}

// CreateRenderer hands the row's container to the canvas. BaseWidget.Resize then
// lays the body out at the row's size, so a recycled row stretches with the list.
func (r *commandRow) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(r.body)
}

func (u *ui) updateRow(id widget.ListItemID, object fyne.CanvasObject) {
	row, ok := object.(*commandRow)
	if !ok {
		return
	}

	keys := u.model.FilteredKeys()
	if int(id) >= len(keys) {
		row.name.SetText("")
		return
	}

	key := keys[id]

	label := key
	if command, found := u.model.Commands.Get(key); found {
		label = commandLabel(command)
	}

	row.name.SetText(label)

	row.run.OnTapped = func() { u.runKey(key) }
	row.edit.OnTapped = func() { u.showEditFor(key) }
	row.remove.OnTapped = func() { u.showDeleteFor(key) }
}

// setupOption maps a selector label to the setup it selects; the label carries
// the command count, so it is not the same string as the name.
type setupOption struct {
	label string
	name  string
}

// setupChoicesList is the selector's entries with their command counts, or a
// single entry when the window was built without a setup source (the tests that
// only exercise another tab).
func (u *ui) setupChoicesList() ([]string, []setupOption) {
	if u.setups == nil {
		return []string{u.activeSetup}, []setupOption{{label: u.activeSetup, name: u.activeSetup}}
	}

	names := u.setups.Names()
	labels := make([]string, 0, len(names))
	choices := make([]setupOption, 0, len(names))

	for _, name := range names {
		label := u.setupLabel(name)
		labels = append(labels, label)
		choices = append(choices, setupOption{label: label, name: name})
	}

	return labels, choices
}

// setupLabel is what the selector shows: the name and how many commands it
// holds. A count that cannot be read shows as "?" rather than hiding the setup.
func (u *ui) setupLabel(name string) string {
	if u.setups == nil {
		return name
	}

	count, err := u.setups.Commands(name)
	if err != nil {
		return name + " (?)"
	}

	return fmt.Sprintf("%s (%d)", name, count)
}

// setupName maps a selected label back to the setup's name.
func (u *ui) setupName(label string) string {
	for _, choice := range u.setupChoices {
		if choice.label == label {
			return choice.name
		}
	}

	return ""
}

// refreshSetupSelect rebuilds the switcher's options and reselects the active
// setup, which is what puts the list back after a failed switch.
func (u *ui) refreshSetupSelect() {
	if u.setupSelect == nil || u.setups == nil {
		return
	}

	labels, choices := u.setupChoicesList()
	u.setupChoices = choices

	u.setupSelect.Options = labels
	u.setupSelect.Refresh()
	u.setupSelect.SetSelected(u.setupLabel(u.activeSetup))
}

// switchSetup changes the setup and rebuilds everything the window shows from
// it. The session is reloaded in place, so the ports keep working.
func (u *ui) switchSetup(name string) bool {
	if u.setups == nil {
		return false
	}

	state, err := u.setups.Switch(name)
	if err != nil {
		showError(err.Error(), u.window)
		u.refreshSetupSelect()
		return false
	}

	u.activeSetup = name
	u.model.Reload(state.Catalog, state.History)
	u.loadGroupIcons()

	if u.list != nil {
		u.list.Refresh()
		u.list.UnselectAll()
	}

	u.clearDetail()

	u.refreshSetupSelect()
	u.refreshGroups()
	u.refreshHistory()
	u.window.SetTitle("crower — " + name)

	return true
}

// commandLabel is what a row shows: the icon when one is set, then the name.
func commandLabel(command catalog.Command) string {
	if command.Icon == "" {
		return command.Name
	}

	return command.Icon + " " + command.Name
}

func (u *ui) selectRow(id widget.ListItemID) {
	u.model.Index = int(id)

	selected := u.model.Selected()
	if selected == "" {
		u.clearDetail()
		return
	}

	command, found := u.model.Commands.Get(selected)
	if !found {
		u.detail.SetText(selected)
		setFieldGrid(u.detailFields, nil)
		u.setDetailVisible(true)
		return
	}

	u.detail.SetText(commandLabel(command))
	setFieldGrid(u.detailFields, commandDetailFields(command))
	u.paramsEntry.SetPlaceHolder(paramsHint(command))
	u.setDetailVisible(true)
}

// selectCommand brings a command into view from another tab: it clears the
// filters, opens Commands and selects the row, so the detail panel follows.
func (u *ui) selectCommand(key string) {
	if u.filterEntry != nil {
		u.filterEntry.SetText("")
	}

	u.model.Filter = ""

	if u.groupFilter != nil {
		u.groupFilter.SetSelected(allGroupsLabel)
	}

	u.model.Group = ""

	if u.tabs != nil {
		u.tabs.SelectIndex(0)
	}

	if u.list == nil {
		return
	}

	u.list.Refresh()
	u.list.UnselectAll()
	u.clearDetail()

	index := slices.Index(u.model.FilteredKeys(), key)
	if index < 0 {
		return
	}

	u.model.Index = index
	u.list.Select(index)
}

// setDetailVisible is the Details toggle.
//
// Showing or hiding is not enough on its own: Fyne's layouts skip an invisible
// child, so a panel hidden at build time has no size, and making it visible does
// not re-run the parent's layout. Refreshing the tab root lays the whole chain
// out again, which is what gives the panel a height to draw its text in.
func (u *ui) setDetailVisible(visible bool) {
	u.detailShown = visible

	if u.detailBody == nil {
		return
	}

	if visible {
		u.detailBody.Show()
	} else {
		u.detailBody.Hide()
	}

	if u.commandsRoot != nil {
		u.commandsRoot.Refresh()
	}
}

// clearDetail empties and collapses the panel, so a filter change, a reload or
// an unselect cannot leave another command's text on screen.
func (u *ui) clearDetail() {
	if u.detail != nil {
		u.detail.SetText("")
	}

	setFieldGrid(u.detailFields, nil)

	if u.paramsEntry != nil {
		u.paramsEntry.SetText("")
		u.paramsEntry.SetPlaceHolder("Parameters (space-separated)")
	}

	u.setDetailVisible(false)
}

// paramsHint tells the user what the command's parameters mean without making
// them look like input that is already filled in.
func paramsHint(command catalog.Command) string {
	if len(command.Params) == 0 {
		return "Parameters (space-separated)"
	}

	return "Parameters — declared: " + strings.Join(command.Params, " ")
}

func (u *ui) runSelected() {
	u.runKey(u.model.Selected())
}

// runKey runs one command: the row's key when a row button asked, the
// selection when the toolbar did.
//
// It is smart about parameters: a command that declares them asks for them, one
// that does not runs straight away. The detail panel's Run button is the path
// for supplying them from the window without a prompt.
func (u *ui) runKey(key string) {
	if key == "" {
		return
	}

	command, found := u.model.Commands.Get(key)
	if !found {
		return
	}

	if len(command.Params) > 0 {
		u.askParams(key, command)
		return
	}

	u.model.StartRunKey(key, false)
	u.showOutput()
}

// runFromDetail runs the selected command with the detail panel's entry as its
// parameters. An empty entry runs it without parameters.
func (u *ui) runFromDetail() {
	key := u.model.Selected()
	if key == "" {
		return
	}

	u.model.StartRunKey(key, true)
	u.model.ParamInput = u.paramsEntry.Text
	u.model.SubmitParams()
	u.showOutput()
}

// commandDetailFields is the detail panel's labelled metadata, in the order a
// reader scans it: what it is called, what it runs, where it lives, and the
// limits. Empty values are omitted rather than shown as blank rows.
func commandDetailFields(command catalog.Command) []detailField {
	fields := make([]detailField, 0, 10)

	if command.Name != "" {
		fields = append(fields, detailField{Label: "Name", Value: command.Name})
	}

	if command.Exec != "" {
		fields = append(fields, detailField{Label: "Exec", Value: command.Exec})
	}

	if len(command.Aliases) > 0 {
		fields = append(fields, detailField{Label: "Aliases", Value: strings.Join(command.Aliases, ", ")})
	}

	if command.Group != "" {
		fields = append(fields, detailField{Label: "Group", Value: command.Group})
	}

	if command.Description != "" {
		fields = append(fields, detailField{Label: "Description", Value: command.Description})
	}

	if len(command.Params) > 0 {
		fields = append(fields, detailField{Label: "Params", Value: strings.Join(command.Params, ", ")})
	}

	if command.Cwd != "" {
		fields = append(fields, detailField{Label: "Cwd", Value: command.Cwd})
	}

	if command.Shell != "" {
		fields = append(fields, detailField{Label: "Shell", Value: command.Shell})
	}

	if command.Timeout != "" {
		fields = append(fields, detailField{Label: "Timeout", Value: command.Timeout})
	}

	if len(command.Steps) > 0 {
		fields = append(fields, detailField{Label: "Steps", Value: numberedSteps(command.Steps)})
	}

	return fields
}

// numberedSteps keeps a multi-step command readable: the panel used to show
// only the count, which said nothing about what would run.
func numberedSteps(steps []string) string {
	lines := make([]string, 0, len(steps))

	for i, step := range steps {
		lines = append(lines, fmt.Sprintf("%d. %s", i+1, step))
	}

	return strings.Join(lines, "\n")
}

// reload replaces the model state after an edit and refreshes every view built
// from it. The widgets are optional so a test can exercise the submit paths
// without a window.
func (u *ui) reload(commands catalog.Catalog, history timeline.Log) {
	u.model.Reload(commands, history)

	if u.list != nil {
		u.list.Refresh()
		u.list.UnselectAll()
	}

	u.clearDetail()

	u.refreshGroups()
	u.refreshHistory()
}

func splitAliases(text string) []string {
	return splitList(text)
}

// splitList reads a comma-separated list, dropping the empty entries a trailing
// comma or a stray space would leave. Aliases and params share it.
func splitList(text string) []string {
	items := make([]string, 0)

	for _, part := range strings.Split(text, ",") {
		part = strings.TrimSpace(part)
		if part != "" {
			items = append(items, part)
		}
	}

	return items
}
