//go:build linux

package system_test

import (
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
)

// The child process must see the incremented depth: that is what lets a nested
// chain stop at the limit.
func TestRunNestedPassesTheDepthToTheChild(t *testing.T) {
	t.Setenv("CROWER_EXEC_DEPTH", "")

	out, err := system.NewRunner().RunNested("echo $CROWER_EXEC_DEPTH", system.Target{})
	if err != nil {
		t.Fatalf("RunNested: %v", err)
	}

	if got := strings.TrimSpace(out); got != "1" {
		t.Errorf("child depth = %q, want 1", got)
	}
}

func TestRunNestedIncrementsAnExistingDepth(t *testing.T) {
	t.Setenv("CROWER_EXEC_DEPTH", "2")

	out, err := system.NewRunner().RunNested("echo $CROWER_EXEC_DEPTH", system.Target{})
	if err != nil {
		t.Fatalf("RunNested: %v", err)
	}

	if got := strings.TrimSpace(out); got != "3" {
		t.Errorf("child depth = %q, want 3", got)
	}
}

// A command that declares the depth key cannot shadow the guard: the nested
// value is written as an override, not appended after the command's own.
func TestRunNestedDepthBeatsACommandOverride(t *testing.T) {
	t.Setenv("CROWER_EXEC_DEPTH", "")

	out, err := system.NewRunner().RunNested("echo $CROWER_EXEC_DEPTH", system.Target{
		Env: map[string]string{"CROWER_EXEC_DEPTH": "0"},
	})
	if err != nil {
		t.Fatalf("RunNested: %v", err)
	}

	if got := strings.TrimSpace(out); got != "1" {
		t.Errorf("child depth = %q, want 1", got)
	}
}
