package helper_test

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/helper"
)

// newInputApp builds an application whose prompts read the given text.
func newInputApp(t *testing.T, input string, noInput bool) *app.Config {
	t.Helper()

	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	cfg, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(input), Out: io.Discard, Err: io.Discard},
		NoInput: noInput,
	})
	if err != nil {
		t.Fatalf("InitApp() error = %v", err)
	}

	return cfg
}

func TestGetUserConfirmationReadsTheAnswer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"yes", "y\n", true},
		{"yes spelled out", "yes\n", true},
		{"no", "n\n", false},
		{"no spelled out", "no\n", false},
		{"invalid, then yes", "maybe\ny\n", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helper.GetUserConfirmation("delete", newInputApp(t, tt.input, false))
			if err != nil {
				t.Fatalf("GetUserConfirmation() error = %v", err)
			}

			if got != tt.want {
				t.Errorf("GetUserConfirmation() = %v, want %v", got, tt.want)
			}
		})
	}
}

// --no-input means no prompt is asked at all, so a run in a script fails with a
// message instead of waiting for a person who is not there.
func TestNoInputFailsBeforeAsking(t *testing.T) {
	_, err := helper.GetUserConfirmation("delete", newInputApp(t, "", true))
	if err == nil {
		t.Fatal("GetUserConfirmation() with --no-input returned no error")
	}

	if !strings.Contains(err.Error(), "input") {
		t.Errorf("error = %v, want it to say that input is required", err)
	}
}

// A closed stream is the same failure as --no-input. It used to be an empty
// answer, and the prompt looped forever.
func TestClosedInputIsAnError(t *testing.T) {
	if _, err := helper.GetUserConfirmation("delete", newInputApp(t, "", false)); err == nil {
		t.Fatal("GetUserConfirmation() on a closed stream returned no error")
	}

	if _, err := helper.GetUserInput("name", helper.IsValidInput, newInputApp(t, "", false)); err == nil {
		t.Fatal("GetUserInput() on a closed stream returned no error")
	}
}

func TestIsValidInputRejectsAnEmptyAnswer(t *testing.T) {
	cfg := newInputApp(t, "", false)

	if _, err := helper.IsValidInput("", cfg); err == nil {
		t.Fatal("IsValidInput(\"\") returned no error")
	}

	if got, err := helper.IsValidInput("value", cfg); err != nil || got != "value" {
		t.Errorf("IsValidInput(\"value\") = %v, %v; want \"value\", nil", got, err)
	}
}

// Row numbers are parsed, not assumed: a letter used to panic the conversion.
func TestIsValidInputKeyRejectsNonNumbers(t *testing.T) {
	cfg := newInputApp(t, "", false)

	if _, err := helper.IsValidInputKey("abc", cfg); err == nil {
		t.Fatal("IsValidInputKey(\"abc\") returned no error")
	}

	if _, err := helper.IsValidInputKey("7", cfg); err == nil {
		t.Fatal("IsValidInputKey(\"7\") on an empty table returned no error")
	}
}

func TestIsValidInputKeyMapsARowToItsCommand(t *testing.T) {
	cfg := newInputApp(t, "", false)
	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo hello"))
	cfg.OrderKeys = []string{"greet"}

	got, err := helper.IsValidInputKey("0", cfg)
	if err != nil {
		t.Fatalf("IsValidInputKey(\"0\") error = %v", err)
	}

	if got != "greet" {
		t.Errorf("IsValidInputKey(\"0\") = %v, want \"greet\"", got)
	}
}
