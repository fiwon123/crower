//go:build fyne

// Package fyneui is the Fyne front end. It renders the headless desktop model
// and owns no rules of its own; the default binary does not build it, so the
// cgo-free release is unchanged.
package fyneui

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

	// Restart is deferred until the window loop has returned: on Unix the
	// updater replaces this process image, which cannot happen while the Fyne
	// loop is still running. A failure has no window left to report to, so it
	// goes to stderr (the updater also writes restart-error.log).
	if front.restartRequested && front.updater != nil {
		if err := front.updater.Restart(); err != nil {
			fmt.Fprintln(os.Stderr, "restart failed:", err)
		}
	}
}

// ui is the widget set the tabs share. Reloading after an edit has to refresh
// both the command list and the history cards, so both live here instead of
// being captured by one tab.
type ui struct {
	model    *desktop.Model
	window   fyne.Window
	editor   Editor
	reverter Reverter
	updater  Updater

	// restartRequested is set by the About tab's Restart button; Run performs the
	// actual restart after the window loop ends.
	restartRequested bool
	tools            desktop.Tools
	groups           desktop.Groups
	setups           desktop.Setups

	activeSetup  string
	setupSelect  *widget.Select
	setupChoices []setupOption

	groupIcons   map[string]string
	groupOptions []groupOption

	// tabs is the AppTabs itself: a drop needs to open the Tools tab, and the
	// Commands tab changes the setup selector's labels, so both need it.
	tabs *container.AppTabs

	// commandScroll holds the command cards; commandPage is the VBox the page
	// refresh re-lays out when a card opens or closes.
	commandScroll *container.Scroll
	commandPage   *fyne.Container
	commandList   *fyne.Container
	commandEmpty  *widget.Label
	commandCards  []*commandCard
	commandOpen   *commandCard

	filterEntry *widget.Entry
	groupFilter *widget.Select

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

// commandsTab is the command cards with the filter and the setup/group
// controls. Collapsed cards are just a name; opening one reveals its fields,
// parameters and actions.
func (u *ui) commandsTab() fyne.CanvasObject {
	filter := newScrollFreeEntry()
	u.filterEntry = filter
	filter.SetPlaceHolder("Filter commands")
	filter.OnChanged = func(string) {
		u.model.Filter = filter.Text
		u.model.Index = 0
		u.rebuildCommands()
	}

	newCommand := widget.NewButtonWithIcon("New command", theme.ContentAddIcon(), u.showCreate)

	u.groupFilter = widget.NewSelect([]string{allGroupsLabel}, u.selectGroup)
	u.refreshGroups()

	groupsButton := widget.NewButtonWithIcon("Groups", theme.GridIcon(), u.showGroups)

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

	setupRow := container.NewHBox(widget.NewLabel("Setup"), u.setupSelect)

	if u.setups != nil {
		removeSetup := widget.NewButtonWithIcon("Delete setup", theme.DeleteIcon(), u.showDeleteSetup)
		removeSetup.Importance = widget.DangerImportance

		setupRow.Add(widget.NewButtonWithIcon("New setup", theme.ContentAddIcon(), u.showNewSetup))
		setupRow.Add(removeSetup)
	}

	groupRow := container.NewHBox(widget.NewLabel("Group"), u.groupFilter, groupsButton)

	filterRow := container.NewBorder(nil, nil, widget.NewLabel("Filter"), nil, filter)

	// Create sits at the top right, above the list it adds to; the filter reads
	// under it with its own label, and a rule separates the controls from the
	// cards so the two zones never look like one block.
	createRow := container.NewHBox(layout.NewSpacer(), newCommand)

	controls := container.NewVBox(
		gap(0, theme.Padding()),
		createRow,
		filterRow,
		widget.NewSeparator(),
		setupRow,
		gap(0, theme.Padding()/2),
		groupRow,
		widget.NewSeparator(),
	)

	u.commandEmpty = widget.NewLabel("No commands match the filter.")
	u.commandEmpty.Wrapping = fyne.TextWrapWord
	u.commandEmpty.Selectable = true
	u.commandEmpty.Hide()

	u.commandList = container.NewVBox()
	u.commandPage = container.NewVBox(u.commandEmpty, u.commandList)
	u.commandScroll = container.NewVScroll(page(u.commandPage))

	u.rebuildCommands()

	// The filter and setup rows sit outside the scroll, so they get the page's
	// lateral margins: the same left inset and the same right gutter.
	top := container.NewBorder(nil, nil, gap(theme.Padding()*2, 0), gap(scrollGutter(), 0), controls)

	return container.NewBorder(top, nil, nil, nil, u.commandScroll)
}

// rebuildCommands replaces the cards from the current filter, group and setup.
// It drops any open card: its command may no longer be in the list.
func (u *ui) rebuildCommands() {
	if u.commandList == nil {
		return
	}

	u.commandOpen = nil
	u.commandCards = nil

	objects := make([]fyne.CanvasObject, 0)

	for _, key := range u.model.FilteredKeys() {
		command, found := u.model.Commands.Get(key)
		if !found {
			// A key the catalog cannot resolve would otherwise vanish from the
			// tab; show its name so the entry is still visible.
			command = catalog.Command{Name: key}
		}

		card := u.newCommandCard(command)
		u.commandCards = append(u.commandCards, card)
		objects = append(objects, card.expandableCard)
	}

	u.commandList.Objects = objects
	u.commandList.Refresh()

	if len(objects) == 0 {
		text := "No commands match the filter."
		if u.model.Commands.Len() == 0 {
			text = "No commands yet. Use New command to add one."
		}

		u.commandEmpty.SetText(text)
		u.commandEmpty.Show()
	} else {
		u.commandEmpty.Hide()
	}

	u.commandPage.Refresh()

	// A filter can shrink the page under the scroll's offset; go back to the
	// top rather than render into blank space.
	if u.commandScroll != nil {
		u.commandScroll.ScrollToOffset(fyne.NewPos(0, 0))
	}
}

// selectCommand brings a command into view from another tab: it clears the
// filters, opens Commands and opens the command's card.
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

	u.rebuildCommands()

	for index, card := range u.commandCards {
		if card.key != key {
			continue
		}

		u.model.Index = index
		u.toggleCommandCard(card)

		return
	}
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
	u.rebuildCommands()
}

// groupNames is every group the window offers: the ones the commands use and
// the ones stored on their own, sorted and de-duplicated. A group can exist
// before a command joins it, which is what the New group action creates.
func (u *ui) groupNames() []string {
	seen := make(map[string]bool, len(u.groupIcons)+4)
	names := make([]string, 0, len(u.groupIcons)+4)

	for _, name := range u.model.Groups() {
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}

	for name := range u.groupIcons {
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}

	slices.Sort(names)

	return names
}

// refreshGroups rebuilds the filter's options after an edit, a removal or a
// group icon change. A group with no commands and no stored record disappears,
// and the filter falls back to all.
func (u *ui) refreshGroups() {
	if u.groupFilter == nil {
		return
	}

	names := u.groupNames()
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
	u.rebuildCommands()

	u.refreshSetupSelect()
	u.refreshGroups()
	u.refreshHistory()
	u.window.SetTitle("crower — " + name)

	return true
}

// commandLabel is what a card header shows: the icon when one is set, then the name.
func commandLabel(command catalog.Command) string {
	if command.Icon == "" {
		return command.Name
	}

	return command.Icon + " " + command.Name
}

// paramsHint tells the user what the command's parameters mean without making
// them look like input that is already filled in.
func paramsHint(command catalog.Command) string {
	if len(command.Params) == 0 {
		return "Parameters (space-separated)"
	}

	return "Parameters — declared: " + strings.Join(command.Params, " ")
}

// commandDetailFields is a card's labelled metadata, in the order a
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

// numberedSteps keeps a multi-step command readable: the card shows the steps
// themselves, not only a count.
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

	u.rebuildCommands()
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
