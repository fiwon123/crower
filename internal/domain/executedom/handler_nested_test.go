package executedom

import (
	"strconv"
	"strings"
	"testing"
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

	_, err := (&Handler{}).executeNested("crower list")
	if err == nil {
		t.Fatal("executeNested() at the maximum depth returned no error")
	}

	if !strings.Contains(err.Error(), "maximum depth") {
		t.Errorf("error = %v, want it to name the limit", err)
	}
}
