package executedom

import (
	"strconv"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/shell"
)

// The nesting level comes from the environment, because the recursive call is a
// new crower process: an in-process counter never sees it.
func TestNestedDepthReadsTheEnvironment(t *testing.T) {
	t.Setenv(nestedDepthEnv, "3")

	if got := nestedDepth(); got != 3 {
		t.Errorf("nestedDepth() = %d, want 3", got)
	}
}

func TestNestedDepthIsZeroWithoutTheVariable(t *testing.T) {
	t.Setenv(nestedDepthEnv, "")

	if got := nestedDepth(); got != 0 {
		t.Errorf("nestedDepth() = %d, want 0", got)
	}
}

// A command whose steps call itself would spawn a process per step forever. The
// limit is the only thing standing between that command and the machine running
// out of processes.
func TestExecuteNestedRefusesToGoTooDeep(t *testing.T) {
	t.Setenv(nestedDepthEnv, strconv.Itoa(maxExecDepth))

	_, err := executeNested(shell.New(), "crower list")
	if err == nil {
		t.Fatal("executeNested() at the maximum depth returned no error")
	}

	if !strings.Contains(err.Error(), "maximum depth") {
		t.Errorf("error = %v, want it to name the limit", err)
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
