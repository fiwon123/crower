package tui

import (
	"bytes"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

type fakeExecutor struct {
	keys   []string
	params [][]string
	out    string
	name   string
	err    error
}

func (e *fakeExecutor) Execute(key string, params []string) (string, string, error) {
	e.keys = append(e.keys, key)
	e.params = append(e.params, params)

	return e.out, e.name, e.err
}

func greeter() catalog.Catalog {
	return catalog.New(nil, []catalog.Command{
		{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"},
	})
}

func runMenu(t *testing.T, input string, commands catalog.Catalog, history timeline.Log, exec Executor) string {
	t.Helper()

	var out bytes.Buffer

	ui := newTUI(commands, history, exec, strings.NewReader(input), &out)

	if err := ui.Run(); err != nil {
		t.Fatalf("Run: %v", err)
	}

	return out.String()
}

func TestExit(t *testing.T) {
	text := runMenu(t, "6\n", greeter(), timeline.Log{}, &fakeExecutor{})

	if !strings.Contains(text, "Goodbye!") {
		t.Errorf("output = %q, want the Goodbye line", text)
	}
}

func TestInvalidOptionThenExit(t *testing.T) {
	text := runMenu(t, "9\n6\n", greeter(), timeline.Log{}, &fakeExecutor{})

	if !strings.Contains(text, "Invalid option") {
		t.Errorf("output = %q, want the invalid-option line", text)
	}

	if !strings.Contains(text, "Goodbye!") {
		t.Errorf("output = %q, want the loop to continue to exit", text)
	}
}

func TestListCommands(t *testing.T) {
	text := runMenu(t, "1\n6\n", greeter(), timeline.Log{}, &fakeExecutor{})

	for _, want := range []string{"=== Commands ===", "1. greet - echo hi", "Aliases: [hi]"} {
		if !strings.Contains(text, want) {
			t.Errorf("output = %q, want %q", text, want)
		}
	}
}

func TestListCommandsEmpty(t *testing.T) {
	text := runMenu(t, "1\n6\n", catalog.New(nil, nil), timeline.Log{}, &fakeExecutor{})

	if !strings.Contains(text, "No commands found") {
		t.Errorf("output = %q, want the empty line", text)
	}
}

func TestRunByNumber(t *testing.T) {
	exec := &fakeExecutor{out: "hello", name: "greet"}

	text := runMenu(t, "2\n1\n\n6\n", greeter(), timeline.Log{}, exec)

	if !slices.Equal(exec.keys, []string{"greet"}) {
		t.Errorf("keys = %v, want greet", exec.keys)
	}

	if len(exec.params) != 1 || len(exec.params[0]) != 0 {
		t.Errorf("params = %v, want none", exec.params)
	}

	for _, want := range []string{"--- Executing ---", "Command: greet", "hello", "--- Done ---"} {
		if !strings.Contains(text, want) {
			t.Errorf("output = %q, want %q", text, want)
		}
	}
}

func TestRunByName(t *testing.T) {
	exec := &fakeExecutor{}

	runMenu(t, "2\ngreet\n\n6\n", greeter(), timeline.Log{}, exec)

	if !slices.Equal(exec.keys, []string{"greet"}) {
		t.Errorf("keys = %v, want greet", exec.keys)
	}
}

func TestRunWithParams(t *testing.T) {
	exec := &fakeExecutor{}

	runMenu(t, "3\n1\none two\n\n6\n", greeter(), timeline.Log{}, exec)

	if !slices.Equal(exec.keys, []string{"greet"}) {
		t.Errorf("keys = %v, want greet", exec.keys)
	}

	if len(exec.params) != 1 || !slices.Equal(exec.params[0], []string{"one", "two"}) {
		t.Errorf("params = %v, want [one two]", exec.params)
	}
}

func TestRunExecutorErrorIsPrinted(t *testing.T) {
	exec := &fakeExecutor{err: errors.New("boom")}

	text := runMenu(t, "2\n1\n\n6\n", greeter(), timeline.Log{}, exec)

	if !strings.Contains(text, "Error: boom") {
		t.Errorf("output = %q, want the executor error", text)
	}

	if !strings.Contains(text, "Goodbye!") {
		t.Errorf("output = %q, want the menu to continue", text)
	}
}

func TestHistory(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Create, "greet", "note")

	text := runMenu(t, "4\n\n6\n", greeter(), history, &fakeExecutor{})

	for _, want := range []string{"=== History ===", "[create] greet"} {
		if !strings.Contains(text, want) {
			t.Errorf("output = %q, want %q", text, want)
		}
	}
}

func TestHistoryEmpty(t *testing.T) {
	text := runMenu(t, "4\n\n6\n", greeter(), timeline.Log{}, &fakeExecutor{})

	if !strings.Contains(text, "No history") {
		t.Errorf("output = %q, want the empty history line", text)
	}
}

func TestExport(t *testing.T) {
	text := runMenu(t, "5\n\n\n6\n", greeter(), timeline.Log{}, &fakeExecutor{})

	if !strings.Contains(text, "Would export to: crower-export.toml") {
		t.Errorf("output = %q, want the default export path", text)
	}
}

func TestExportNamed(t *testing.T) {
	text := runMenu(t, "5\nout.toml\n\n6\n", greeter(), timeline.Log{}, &fakeExecutor{})

	if !strings.Contains(text, "Would export to: out.toml") {
		t.Errorf("output = %q, want the named export path", text)
	}
}
