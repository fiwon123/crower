//go:build fyne

package fyneui

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// stubSetups switches between two canned command sets.
type stubSetups struct {
	current   string
	names     []string
	next      desktop.State
	switched  []string
	err       error
	created   []string
	deleted   []string
	counts    map[string]int
	createErr error
	deleteErr error
	countErr  error
}

func (s *stubSetups) Current() string { return s.current }

func (s *stubSetups) Names() []string { return s.names }

func (s *stubSetups) Switch(name string) (desktop.State, error) {
	if s.err != nil {
		return desktop.State{}, s.err
	}

	s.switched = append(s.switched, name)
	s.current = name

	return s.next, nil
}

func (s *stubSetups) Commands(name string) (int, error) {
	if s.countErr != nil {
		return 0, s.countErr
	}

	return s.counts[name], nil
}

func (s *stubSetups) Create(name string) error {
	if s.createErr != nil {
		return s.createErr
	}

	s.created = append(s.created, name)
	s.names = append(s.names, name)

	return nil
}

func (s *stubSetups) Delete(name string) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}

	s.deleted = append(s.deleted, name)

	kept := make([]string, 0, len(s.names))

	for _, existing := range s.names {
		if existing != name {
			kept = append(kept, existing)
		}
	}

	s.names = kept

	return nil
}

// New… creates the setup and puts it in the selector without switching to it.
func TestCreateSetupAddsItToTheSelector(t *testing.T) {
	setups := &stubSetups{current: "default", names: []string{"default"}}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	if !front.createSetup(" work ") {
		t.Fatal("createSetup returned false")
	}

	if len(setups.created) != 1 || setups.created[0] != "work" {
		t.Fatalf("created = %v, want [work]", setups.created)
	}

	if !slices.Contains(front.setupSelect.Options, front.setupLabel("work")) {
		t.Errorf("selector = %v, want the new setup in it", front.setupSelect.Options)
	}

	if want := front.setupLabel("default"); front.setupSelect.Selected != want {
		t.Errorf("selector selected %q, want %q kept", front.setupSelect.Selected, want)
	}
}

func TestCreateSetupRejectsAnEmptyName(t *testing.T) {
	setups := &stubSetups{current: "default", names: []string{"default"}}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	if front.createSetup("   ") {
		t.Fatal("createSetup accepted an empty name")
	}

	if len(setups.created) != 0 {
		t.Errorf("created = %v, want nothing", setups.created)
	}
}

// Delete… removes the selected setup from the selector.
func TestDeleteSetupRemovesItFromTheSelector(t *testing.T) {
	setups := &stubSetups{current: "default", names: []string{"default", "work"}}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	front.setupSelect.SetSelected("work")

	if !front.deleteSetup("work") {
		t.Fatal("deleteSetup returned false")
	}

	if len(setups.deleted) != 1 || setups.deleted[0] != "work" {
		t.Fatalf("deleted = %v, want [work]", setups.deleted)
	}

	if slices.Contains(front.setupSelect.Options, front.setupLabel("work")) {
		t.Errorf("selector = %v, want the deleted setup gone", front.setupSelect.Options)
	}
}

// A refusal (the active setup, the default, a missing one) leaves the selector
// alone.
func TestDeleteSetupReportsAFailure(t *testing.T) {
	setups := &stubSetups{
		current:   "work",
		names:     []string{"default", "work"},
		deleteErr: errors.New(`setup "work" is active; switch first`),
	}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	if front.deleteSetup("work") {
		t.Fatal("deleteSetup accepted a failing delete")
	}

	if !slices.Contains(front.setupSelect.Options, front.setupLabel("work")) {
		t.Errorf("selector = %v, want the failed setup kept", front.setupSelect.Options)
	}
}

// The selector's options say how many commands each setup holds, and choosing
// one still switches by name.
func TestSetupOptionsShowTheCommandCount(t *testing.T) {
	setups := &stubSetups{
		current: "default",
		names:   []string{"default", "work"},
		counts:  map[string]int{"default": 3, "work": 12},
	}

	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	want := []string{"default (3)", "work (12)"}
	if !slices.Equal(front.setupSelect.Options, want) {
		t.Fatalf("options = %v, want %v", front.setupSelect.Options, want)
	}

	if front.setupSelect.Selected != "default (3)" {
		t.Errorf("selected = %q, want the active setup's label", front.setupSelect.Selected)
	}

	// Choosing a label, not a bare name, still switches to the setup.
	front.setupSelect.SetSelected("work (12)")

	if len(setups.switched) != 1 || setups.switched[0] != "work" {
		t.Fatalf("switched = %v, want [work]", setups.switched)
	}
}

// A count that cannot be read shows as "?" instead of hiding the setup.
func TestSetupLabelShowsUnknownCounts(t *testing.T) {
	setups := &stubSetups{current: "default", names: []string{"default"}, countErr: errors.New("unreadable")}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups

	if got := front.setupLabel("default"); got != "default (?)" {
		t.Errorf("setupLabel = %q, want the unknown marker", got)
	}
}

func TestSwitchSetupRebuildsTheWindow(t *testing.T) {
	setups := &stubSetups{
		current: "default",
		names:   []string{"default", "work"},
		next: desktop.State{
			Catalog: catalog.New(nil, []catalog.Command{{Name: "deploy", Exec: "make deploy", Group: "ship"}}),
			History: timeline.Log{},
		},
	}

	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	if !front.switchSetup("work") {
		t.Fatal("switchSetup returned false")
	}

	if len(setups.switched) != 1 || setups.switched[0] != "work" {
		t.Fatalf("switched = %v, want [work]", setups.switched)
	}

	if got := front.model.Selected(); got != "deploy" {
		t.Errorf("Selected() = %q, want the new setup's command", got)
	}

	if len(front.commandCards) != 1 || front.commandCards[0].key != "deploy" {
		t.Errorf("cards = %v, want the new setup's command rendered", front.commandCards)
	}

	if front.activeSetup != "work" {
		t.Errorf("activeSetup = %q, want work", front.activeSetup)
	}

	if want := front.setupLabel("work"); front.setupSelect.Selected != want {
		t.Errorf("setup select = %q, want %q", front.setupSelect.Selected, want)
	}

	if title := window.Title(); !strings.Contains(title, "work") {
		t.Errorf("window title = %q, want the setup named", title)
	}
}

func TestSwitchSetupReportsAFailureAndKeepsTheSetup(t *testing.T) {
	setups := &stubSetups{current: "default", names: []string{"default", "work"}, err: errors.New("cannot read setup")}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1"})
	front.setups = setups
	window.SetContent(front.commandsTab())
	window.Resize(fyneSize())

	if front.switchSetup("work") {
		t.Fatal("switchSetup accepted a failing switch")
	}

	if front.activeSetup != "default" {
		t.Errorf("activeSetup = %q, want the previous setup kept", front.activeSetup)
	}

	if want := front.setupLabel("default"); front.setupSelect.Selected != want {
		t.Errorf("setup select = %q, want %q", front.setupSelect.Selected, want)
	}
}

func fyneSize() fyne.Size { return fyne.NewSize(960, 640) }
