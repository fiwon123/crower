package desktop

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// stubExecutor stands in for the command layer's executor: these tests are
// about the model's state, not about running anything.
type stubExecutor struct{}

func (stubExecutor) Execute(string, []string) (string, string, error) {
	return "", "", nil
}

func testApp(t *testing.T) *app.Config {
	t.Helper()

	a, err := app.NewApp(
		filepath.Join(t.TempDir(), "crower.yaml"),
		[]string{},
		command.NewMapData(),
		command.NewMapData(),
		app.Options{Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}},
	)
	if err != nil {
		t.Fatalf("NewApp: %v", err)
	}

	c := command.New("hello", nil, "echo hello")
	a.AllCommandsByName.Add("hello", c)
	a.OrderKeys = append(a.OrderKeys, "hello")

	return a
}

func TestNewModel(t *testing.T) {
	a := testApp(t)
	m := NewModel(a, stubExecutor{})
	if m.State != ViewList {
		t.Fatalf("want list state")
	}
	if len(m.Keys) != 1 {
		t.Fatalf("want 1 key")
	}
}

func TestMove(t *testing.T) {
	a := testApp(t)
	a.AllCommandsByName.Add("b", command.New("b", nil, "echo b"))
	a.OrderKeys = append(a.OrderKeys, "b")
	m := NewModel(a, stubExecutor{})
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
	a := testApp(t)
	a.AllCommandsByName.Add("build", command.New("build", nil, "make build"))
	a.OrderKeys = append(a.OrderKeys, "build")
	m := NewModel(a, stubExecutor{})
	m.Filter = "bui"
	keys := m.FilteredKeys()
	if len(keys) != 1 || keys[0] != "build" {
		t.Fatalf("bad filter")
	}
}

// The history screen names the command an entry touched, and shows a dash for
// an operation that touched no command.
func TestHistoryTextShowsTheCommand(t *testing.T) {
	a := testApp(t)
	a.History.Add(operation.Create, "greet", "note")
	a.History.Add(operation.Copy, "", "note")

	m := NewModel(a, stubExecutor{})
	text := m.historyText()

	if !strings.Contains(text, "[create] greet") {
		t.Errorf("history = %q, want the command on the create row", text)
	}

	if !strings.Contains(text, "[copy] -") {
		t.Errorf("history = %q, want a dash on the command-less row", text)
	}
}
