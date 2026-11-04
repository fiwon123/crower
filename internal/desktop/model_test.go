package desktop

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// stubExecutor stands in for the command layer's executor: these tests are
// about the model's state, not about running anything.
type stubExecutor struct{}

func (stubExecutor) Execute(string, []string) (string, string, error) {
	return "", "", nil
}

// testModel builds the model over one command and an empty history.
func testModel(t *testing.T) *Model {
	t.Helper()

	commands := catalog.New(nil, []catalog.Command{{Name: "hello", Exec: "echo hello"}})

	return NewModel(commands, timeline.Log{}, stubExecutor{})
}

func modelWithCommands(commands []catalog.Command, history timeline.Log) *Model {
	return NewModel(catalog.New(nil, commands), history, stubExecutor{})
}

func TestNewModel(t *testing.T) {
	m := testModel(t)
	if m.State != ViewList {
		t.Fatalf("want list state")
	}
	if len(m.Keys) != 1 {
		t.Fatalf("want 1 key")
	}
}

func TestMove(t *testing.T) {
	commands := catalog.New(
		[]string{"hello", "b"},
		[]catalog.Command{
			{Name: "hello", Exec: "echo hello"},
			{Name: "b", Exec: "echo b"},
		},
	)

	m := NewModel(commands, timeline.Log{}, stubExecutor{})
	m.Move(1)
	if m.Selected() != "b" {
		t.Fatalf("want b")
	}
	m.Move(1)
	if m.Selected() != "hello" {
		t.Fatalf("want wrap")
	}
}

func TestFilter(t *testing.T) {
	m := modelWithCommands([]catalog.Command{
		{Name: "hello", Exec: "echo hello"},
		{Name: "build", Exec: "make build"},
	}, timeline.Log{})
	m.Filter = "bui"
	keys := m.FilteredKeys()
	if len(keys) != 1 || keys[0] != "build" {
		t.Fatalf("bad filter")
	}
}

func TestReloadReplacesTheSetAndResetsTheSelection(t *testing.T) {
	m := testModel(t)
	m.Filter = "hello"
	m.Index = 3

	commands := catalog.New(nil, []catalog.Command{{Name: "build", Exec: "make"}})
	m.Reload(commands, timeline.Log{})

	if got := m.Selected(); got != "build" {
		t.Errorf("Selected() = %q, want build", got)
	}

	if m.Filter != "" || m.Index != 0 {
		t.Errorf("filter/index = %q/%d, want reset", m.Filter, m.Index)
	}
}

// The history screen names the command an entry touched, and shows a dash for
// an operation that touched no command.
func TestHistoryTextShowsTheCommand(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Create, "greet", "note")
	history.Add(operation.Copy, "", "note")

	m := modelWithCommands([]catalog.Command{{Name: "hello", Exec: "echo hello"}}, history)
	text := m.historyText()

	if !strings.Contains(text, "[create] greet") {
		t.Errorf("history = %q, want the command on the create row", text)
	}

	if !strings.Contains(text, "[copy] -") {
		t.Errorf("history = %q, want a dash on the command-less row", text)
	}
}
