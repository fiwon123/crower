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

// stubSetups switches between two canned command sets.
type stubSetups struct {
	current  string
	names    []string
	next     desktop.State
	switched []string
	err      error
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
