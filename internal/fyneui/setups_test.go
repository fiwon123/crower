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
	createErr error
	deleteErr error
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

	if !slices.Contains(front.setupSelect.Options, "work") {
		t.Errorf("selector = %v, want the new setup in it", front.setupSelect.Options)
	}

	if front.setupSelect.Selected != "default" {
		t.Errorf("selector selected %q, want the active setup kept", front.setupSelect.Selected)
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

	if slices.Contains(front.setupSelect.Options, "work") {
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

	if !slices.Contains(front.setupSelect.Options, "work") {
		t.Errorf("selector = %v, want the failed setup kept", front.setupSelect.Options)
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

	if front.activeSetup != "work" {
		t.Errorf("activeSetup = %q, want work", front.activeSetup)
	}

	if front.setupSelect.Selected != "work" {
		t.Errorf("setup select = %q, want work", front.setupSelect.Selected)
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

	if front.setupSelect.Selected != "default" {
		t.Errorf("setup select = %q, want it put back", front.setupSelect.Selected)
	}
}

func fyneSize() fyne.Size { return fyne.NewSize(960, 640) }
