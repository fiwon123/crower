package executecommand

import (
	"runtime"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/shell"
)

func quoted(s string) string {
	return shell.QuoteFor(runtime.GOOS, s)
}

func TestApplyParams(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		params  []string
		want    string
		wantErr bool
	}{
		{"no params leaves the line alone", "echo {1}", nil, "echo {1}", false},
		{"no braces appends", "echo hi", []string{"a b"}, "echo hi " + quoted("a b"), false},
		{"indexed", "echo {1} {2}", []string{"a", "b"}, "echo " + quoted("a") + " " + quoted("b"), false},
		{
			"one pass does not re-scan an inserted parameter",
			"echo {1} {2}",
			[]string{"{2}", "b"},
			"echo " + quoted("{2}") + " " + quoted("b"),
			false,
		},
		{"empty braces take the remaining parameters", "echo {}", []string{"a", "b"}, "echo " + quoted("a") + " " + quoted("b"), false},
		{"empty braces take what an index left", "echo {1} {}", []string{"a", "b", "c"}, "echo " + quoted("a") + " " + quoted("b") + " " + quoted("c"), false},
		{"each empty brace takes the rest", "echo {} {}", []string{"a"}, "echo " + quoted("a") + " " + quoted("a"), false},
		{"a literal brace does not swallow the placeholder", "sed 's/{/X/' {1}", []string{"file.txt"}, "sed 's/{/X/' " + quoted("file.txt"), false},
		{"brace expansion is not a placeholder", "mkdir {a,b}", []string{"x"}, "", true},
		{"an index past the parameters is an error when nothing else matches", "echo {2}", []string{"a"}, "", true},
		{
			"an index past the parameters stays when another placeholder matches",
			"echo {1} {2}",
			[]string{"a"},
			"echo " + quoted("a") + " {2}",
			false,
		},
		{"zero is not an index", "echo {0}", []string{"a"}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyParams(tt.line, tt.params)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("applyParams(%q, %v) returned no error", tt.line, tt.params)
				}

				return
			}

			if err != nil {
				t.Fatalf("applyParams(%q, %v) error = %v", tt.line, tt.params, err)
			}

			if got != tt.want {
				t.Errorf("applyParams(%q, %v) = %q, want %q", tt.line, tt.params, got, tt.want)
			}
		})
	}
}

func TestApplyParamsToSteps(t *testing.T) {
	tests := []struct {
		name    string
		steps   []string
		params  []string
		want    []string
		wantErr bool
	}{
		{"no params leaves the steps", []string{"echo {1}"}, nil, []string{"echo {1}"}, false},
		{"each step substitutes", []string{"echo {1}", "echo {2}"}, []string{"a", "b"}, []string{"echo " + quoted("a"), "echo " + quoted("b")}, false},
		{"a step without braces is untouched", []string{"echo {1}", "echo a"}, []string{"a"}, []string{"echo " + quoted("a"), "echo a"}, false},
		{
			"one pass per step",
			[]string{"echo {1} {2}"},
			[]string{"{2}", "b"},
			[]string{"echo " + quoted("{2}") + " " + quoted("b")},
			false,
		},
		{"an out-of-range index is an error even when another step matches", []string{"echo {1}", "echo {2}"}, []string{"a"}, nil, true},
		{"no placeholder anywhere is an error", []string{"echo a"}, []string{"x"}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyParamsToSteps(tt.steps, tt.params)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("applyParamsToSteps(%v, %v) returned no error", tt.steps, tt.params)
				}

				return
			}

			if err != nil {
				t.Fatalf("applyParamsToSteps(%v, %v) error = %v", tt.steps, tt.params, err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("applyParamsToSteps(%v, %v) = %v, want %v", tt.steps, tt.params, got, tt.want)
			}
		})
	}
}

func TestApplyParamsNamesAnOutOfRangeIndex(t *testing.T) {
	_, err := applyParams("echo {2}", []string{"a"})
	if err == nil {
		t.Fatal("applyParams returned no error")
	}

	if !strings.Contains(err.Error(), "out of range") {
		t.Errorf("error = %v, want it to name the range", err)
	}
}

func TestIsNestedCrower(t *testing.T) {
	tests := []struct {
		name string
		step string
		want bool
	}{
		{"bare crower", "crower list", true},
		{"bare cr", "cr list", true},
		{"leading spaces", "   crower list", true},
		{"single-quoted", "'crower' list", true},
		{"double-quoted", `"crower" list`, true},
		{"quoted alias", `'cr' list`, true},
		{"another command", "echo crower", false},
		{"a prefix, not the word", "crowerx list", false},
		{"crower as an argument", "echo 'crower' list", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNestedCrower(tt.step); got != tt.want {
				t.Errorf("isNestedCrower(%q) = %v, want %v", tt.step, got, tt.want)
			}
		})
	}
}

// A parameter that is exactly "crower" is quoted before substitution; the
// guard has to see through the quoting.
func TestSubstitutedCrowerIsDetected(t *testing.T) {
	line, err := applyParams("{1} list", []string{"crower"})
	if err != nil {
		t.Fatalf("applyParams: %v", err)
	}

	if !isNestedCrower(line) {
		t.Errorf("isNestedCrower(%q) = false, want true", line)
	}
}
