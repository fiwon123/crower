//go:build fyne

package fyneui

import (
	"errors"
	"strings"
	"testing"

	"fyne.io/fyne/v2"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// stubGroups stores icons in memory.
type stubGroups struct {
	icons map[string]string
	sets  map[string]string
	err   error
}

func (g *stubGroups) Icons() (map[string]string, error) {
	if g.err != nil {
		return nil, g.err
	}

	if g.icons == nil {
		g.icons = map[string]string{}
	}

	return g.icons, nil
}

func (g *stubGroups) SetIcon(name, icon string) (map[string]string, error) {
	if g.err != nil {
		return nil, g.err
	}

	if g.icons == nil {
		g.icons = map[string]string{}
	}

	g.icons[name] = icon

	if g.sets == nil {
		g.sets = map[string]string{}
	}

	g.sets[name] = icon

	return g.icons, nil
}

func (g *stubGroups) Remove(name string) (map[string]string, error) {
	if g.err != nil {
		return nil, g.err
	}

	delete(g.icons, name)

	if g.sets == nil {
		g.sets = map[string]string{}
	}

	delete(g.sets, name)

	return g.icons, nil
}

func groupModel() *desktop.Model {
	return desktop.NewModel(
		catalog.New(nil, []catalog.Command{
			{Name: "ship", Exec: "make deploy", Group: "deploy"},
			{Name: "backup", Exec: "tar czf", Group: "backup"},
		}),
		timeline.Log{},
		stubExecutor{},
	)
}

func TestGroupLabelShowsTheIcon(t *testing.T) {
	if got := groupLabel("🚀", "deploy"); got != "🚀 deploy" {
		t.Errorf("groupLabel = %q, want the icon before the name", got)
	}

	if got := groupLabel("", "deploy"); got != "deploy" {
		t.Errorf("groupLabel = %q, want the bare name", got)
	}
}

// The filter lists a group with its icon once the icons are loaded.
func TestGroupFilterLabelsCarryTheIcon(t *testing.T) {
	groups := &stubGroups{icons: map[string]string{"deploy": "🚀"}}
	front, _ := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.groups = groups

	front.commandsTab()
	front.loadGroupIcons()
	front.refreshGroups()

	options := front.groupFilter.Options
	if len(options) != 3 {
		t.Fatalf("options = %v, want all + two groups", options)
	}

	if !containsLabel(options, "🚀 deploy") || !containsLabel(options, "backup") {
		t.Errorf("options = %v, want the icon on deploy and a bare backup", options)
	}

	// Selecting the decorated label narrows the list to the group behind it.
	front.selectGroup("🚀 deploy")

	if front.model.Group != "deploy" {
		t.Errorf("model group = %q, want deploy", front.model.Group)
	}

	if got := front.model.FilteredKeys(); len(got) != 1 || got[0] != "ship" {
		t.Errorf("filtered = %v, want [ship]", got)
	}
}

func TestApplyGroupIconsUpdatesTheFilter(t *testing.T) {
	groups := &stubGroups{}
	front, _ := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})

	front.groups = groups
	front.commandsTab()
	front.refreshGroups()

	if !front.applyGroupIcons(map[string]string{"deploy": "🚀", "backup": ""}) {
		t.Fatal("applyGroupIcons returned false")
	}

	if groups.sets["deploy"] != "🚀" {
		t.Errorf("sets = %v, want deploy stored", groups.sets)
	}

	if !containsLabel(front.groupFilter.Options, "🚀 deploy") {
		t.Errorf("options = %v, want the new icon in the filter", front.groupFilter.Options)
	}
}

func TestApplyGroupIconsReportsAFailure(t *testing.T) {
	front, _ := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.groups = &stubGroups{err: errors.New("read-only file")}
	front.window.SetContent(front.commandsTab())

	if front.applyGroupIcons(map[string]string{"deploy": "🚀"}) {
		t.Fatal("applyGroupIcons accepted a failing store")
	}
}

func containsLabel(options []string, want string) bool {
	for _, option := range options {
		if strings.Contains(option, want) {
			return true
		}
	}

	return false
}

// Creating a group stores it even when no command uses it, so the filter offers
// it straight away.
func TestCreateGroupIsOfferedWithoutACommand(t *testing.T) {
	groups := &stubGroups{}
	model := desktop.NewModel(catalog.New(nil, nil), timeline.Log{}, stubExecutor{})

	front, _ := newTestUIWithUpdater(t, model, stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.groups = groups
	front.commandsTab()

	if !front.createGroup("deploy", "🚀") {
		t.Fatal("createGroup returned false")
	}

	if groups.sets["deploy"] != "🚀" {
		t.Errorf("sets = %v, want deploy stored with its icon", groups.sets)
	}

	if !containsLabel(front.groupFilter.Options, "🚀 deploy") {
		t.Errorf("options = %v, want the new group in the filter", front.groupFilter.Options)
	}
}

// A group created without an icon is still stored: the name is the record.
func TestCreateGroupWithoutAnIconKeepsTheName(t *testing.T) {
	groups := &stubGroups{}
	front, _ := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.groups = groups
	front.commandsTab()

	if !front.createGroup(" staging ", "") {
		t.Fatal("createGroup returned false")
	}

	if _, found := groups.icons["staging"]; !found {
		t.Errorf("icons = %v, want the group stored without an icon", groups.icons)
	}

	if !containsLabel(front.groupFilter.Options, "staging") {
		t.Errorf("options = %v, want the bare name", front.groupFilter.Options)
	}
}

// Empty names are refused rather than stored.
func TestCreateGroupRejectsAnEmptyName(t *testing.T) {
	front, _ := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.groups = &stubGroups{}
	front.window.SetContent(front.commandsTab())

	if front.createGroup("   ", "🚀") {
		t.Error("createGroup accepted an empty name")
	}
}

// Removing a group takes it out of the filter.
func TestRemoveGroupDropsItFromTheFilter(t *testing.T) {
	groups := &stubGroups{}
	front, _ := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.groups = groups
	front.commandsTab()
	front.createGroup("staging", "")

	if !front.removeGroup("staging") {
		t.Fatal("removeGroup returned false")
	}

	if containsLabel(front.groupFilter.Options, "staging") {
		t.Errorf("options = %v, want staging gone", front.groupFilter.Options)
	}
}

// The dialog fits the window it opens in: never smaller than its minimum, never
// larger than its maximum.
func TestGroupsDialogSizeFitsTheWindow(t *testing.T) {
	front, window := newTestUIWithUpdater(t, groupModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})

	window.Resize(fyne.NewSize(400, 300))
	small := dialogSizeFor(front.window, 380, 260, 720, 620)

	if small.Width > 400 || small.Height > 300 {
		t.Errorf("small window size = %v, want it to fit", small)
	}

	window.Resize(fyne.NewSize(1920, 1080))
	large := dialogSizeFor(front.window, 380, 260, 720, 620)

	if large.Width > 720 || large.Height > 620 {
		t.Errorf("maximised size = %v, want it bounded", large)
	}
}
