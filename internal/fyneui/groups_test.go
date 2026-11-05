//go:build fyne

package fyneui

import (
	"errors"
	"strings"
	"testing"

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

	if icon == "" {
		delete(g.icons, name)
	} else {
		g.icons[name] = icon
	}

	if g.sets == nil {
		g.sets = map[string]string{}
	}

	g.sets[name] = icon

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
