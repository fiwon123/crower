package prompt_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/prompt"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// reader answers from a fixed list of lines, so a prompt loop can be tested
// without a terminal.
func reader(lines ...string) func() (string, error) {
	index := 0

	return func() (string, error) {
		if index >= len(lines) {
			return "", errors.New("input exhausted")
		}

		line := lines[index]
		index++

		return line, nil
	}
}

func newPrompt(read func() (string, error), commands ...catalog.Command) *prompt.Prompt {
	return prompt.New(read, func(string, ...any) {}, false, catalog.New(nil, commands), timeline.Log{})
}

func TestConfirm(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want  bool
	}{
		{"yes", []string{"y"}, true},
		{"full yes", []string{"yes"}, true},
		{"no", []string{"n"}, false},
		{"invalid then yes", []string{"maybe", "y"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newPrompt(reader(tt.lines...)).Confirm("Continue")
			if err != nil {
				t.Fatalf("Confirm: %v", err)
			}

			if got != tt.want {
				t.Errorf("Confirm = %v, want %v", got, tt.want)
			}
		})
	}
}

// --no-input fails before printing the question: a script gets an error, not a
// wait.
func TestNoInputFails(t *testing.T) {
	p := prompt.New(func() (string, error) { t.Fatal("read called under --no-input"); return "", nil },
		func(string, ...any) {}, true, catalog.New(nil, nil), timeline.Log{})

	for name, call := range map[string]func() error{
		"confirm": func() error { _, err := p.Confirm("Q"); return err },
		"ask":     func() error { _, err := p.Ask("Q"); return err },
		"row":     func() error { _, err := p.ChooseCommand(); return err },
		"history": func() error { _, err := p.ChooseHistoryEntry(); return err },
	} {
		t.Run(name, func(t *testing.T) {
			if err := call(); err == nil {
				t.Fatal("no error under --no-input")
			}
		})
	}
}

func TestAskRepeatsUntilNonEmpty(t *testing.T) {
	got, err := newPrompt(reader("", "value")).Ask("Name")
	if err != nil {
		t.Fatalf("Ask: %v", err)
	}

	if got != "value" {
		t.Errorf("Ask = %q, want value", got)
	}
}

func TestAliasesStopOnAnEmptyLine(t *testing.T) {
	aliases, err := newPrompt(reader("a1", "a2", "")).Aliases()
	if err != nil {
		t.Fatalf("Aliases: %v", err)
	}

	if !slices.Equal(aliases, []string{"a1", "a2"}) {
		t.Errorf("aliases = %v, want [a1 a2]", aliases)
	}
}

func TestAliasesRejectNonAlphanumerics(t *testing.T) {
	aliases, err := newPrompt(reader("a-b", "ab", "")).Aliases()
	if err != nil {
		t.Fatalf("Aliases: %v", err)
	}

	if !slices.Equal(aliases, []string{"ab"}) {
		t.Errorf("aliases = %v, want [ab]", aliases)
	}
}

func TestChooseCommand(t *testing.T) {
	p := newPrompt(reader("1"),
		catalog.Command{Name: "first", Exec: "echo 1"},
		catalog.Command{Name: "second", Exec: "echo 2"},
	)

	// catalog.New sorts names: first, second.
	got, err := p.ChooseCommand()
	if err != nil {
		t.Fatalf("ChooseCommand: %v", err)
	}

	if got != "second" {
		t.Errorf("row 1 = %q, want second", got)
	}
}

func TestChooseCommandRejectsAnOutOfRangeRow(t *testing.T) {
	p := newPrompt(reader("9", "0"), catalog.Command{Name: "only", Exec: "echo"})

	got, err := p.ChooseCommand()
	if err != nil {
		t.Fatalf("ChooseCommand: %v", err)
	}

	if got != "only" {
		t.Errorf("row = %q, want only", got)
	}
}

// The history table prints newest first, so row 0 is the last entry.
func TestChooseHistoryEntryUsesThePrintedOrder(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Create, "first", "one")
	history.Add(operation.Update, "second", "two")

	p := prompt.New(reader("0"), func(string, ...any) {}, false, catalog.New(nil, nil), history)

	entry, err := p.ChooseHistoryEntry()
	if err != nil {
		t.Fatalf("ChooseHistoryEntry: %v", err)
	}

	if entry.CommandName != "second" {
		t.Errorf("row 0 = %+v, want the newest entry", entry)
	}
}

func TestChooseHistoryEntryRejectsAnEmptyHistory(t *testing.T) {
	p := prompt.New(reader("0"), func(string, ...any) {}, false, catalog.New(nil, nil), timeline.Log{})

	if _, err := p.ChooseHistoryEntry(); err == nil {
		t.Fatal("an empty history returned no error")
	}
}
