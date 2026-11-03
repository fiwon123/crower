package executedom

import (
	"runtime"
	"slices"
	"strings"
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/helper"
)

func quoted(s string) string {
	return helper.QuoteFor(runtime.GOOS, s)
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

// Execute must surface the missing placeholder instead of running a line that
// silently dropped the parameters.
func TestExecuteRejectsParametersWithoutAPlaceholder(t *testing.T) {
	t.Chdir(t.TempDir()) // if the guard regresses, mkdir lands here, not in the tree

	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	app.AllCommandsByName.Add("brace", command.New("brace", nil, "mkdir {a,b}"))

	_, _, err = NewHandler(app).Execute("brace", []string{"x"})
	if err == nil {
		t.Fatal("Execute with a parameter and no placeholder returned no error")
	}

	if !strings.Contains(err.Error(), "no parameter placeholder") {
		t.Errorf("error = %v, want it to name the missing placeholder", err)
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

// A Steps command runs its steps, not Exec. An unusable brace in Exec must not
// turn into an error for a line that never runs.
func TestExecuteWithStepsDoesNotValidateExec(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	multi := command.New("multi", nil, "mkdir {a,b}")
	multi.Steps = []string{"echo {1}"}
	app.AllCommandsByName.Add("multi", multi)

	out, _, err := NewHandler(app).Execute("multi", []string{"a"})
	if err != nil {
		t.Fatalf("Execute with steps returned an error: %v", err)
	}

	if !strings.Contains(out, "a") {
		t.Errorf("output = %q, want the step's output", out)
	}
}

func TestExecuteAppliesParamsToSteps(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	multi := command.New("multi", nil, "unused")
	multi.Steps = []string{"echo {1}", "echo {2}"}
	app.AllCommandsByName.Add("multi", multi)

	out, _, err := NewHandler(app).Execute("multi", []string{"a", "b"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	for _, want := range []string{"a", "b"} {
		if !strings.Contains(out, want) {
			t.Errorf("output = %q, want %q from a step", out, want)
		}
	}
}

// A steps set with no placeholder is an error, not a silent drop.
func TestExecuteRejectsParametersWhenNoStepHasAPlaceholder(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	multi := command.New("multi", nil, "unused")
	multi.Steps = []string{"echo a"}
	app.AllCommandsByName.Add("multi", multi)

	_, _, err = NewHandler(app).Execute("multi", []string{"x"})
	if err == nil {
		t.Fatal("Execute with parameters and no step placeholder returned no error")
	}

	if !strings.Contains(err.Error(), "no parameter placeholder") {
		t.Errorf("error = %v, want it to name the missing placeholder", err)
	}
}

// A step with a literal brace is left alone while another step substitutes.
// {x} is literal in every shell; {a,b} would be brace-expanded by bash.
func TestExecuteKeepsALiteralBraceInAStep(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	multi := command.New("multi", nil, "unused")
	multi.Steps = []string{"echo {x}", "echo {1}"}
	app.AllCommandsByName.Add("multi", multi)

	out, _, err := NewHandler(app).Execute("multi", []string{"param"})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if !strings.Contains(out, "{x}") || !strings.Contains(out, "param") {
		t.Errorf("output = %q, want the literal step and the substituted one", out)
	}
}

// Executing a Steps command must not rewrite the steps stored on the map.
func TestExecuteDoesNotRewriteStoredSteps(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	multi := command.New("multi", nil, "unused")
	multi.Steps = []string{"echo {1}"}
	app.AllCommandsByName.Add("multi", multi)

	if _, _, err := NewHandler(app).Execute("multi", []string{"a"}); err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if stored := app.AllCommandsByName.Get("multi"); stored == nil || stored.Steps[0] != "echo {1}" {
		t.Errorf("stored steps = %+v, want them unchanged", stored)
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
