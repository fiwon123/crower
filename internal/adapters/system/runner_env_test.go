package system_test

import (
	"errors"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/fiwon123/crower/internal/adapters/system"
	crerrors "github.com/fiwon123/crower/internal/errors"
)

// The command's environment reaches the child on top of the process one.
func TestRunMergesTheEnvironment(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	out, err := system.NewRunner().Run(`printf '%s' "$STAGE"`, system.Target{
		Env: map[string]string{"STAGE": "prod"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got := strings.TrimSpace(out); got != "prod" {
		t.Errorf("STAGE = %q, want prod", got)
	}
}

// A nested crower step sees the environment too: it is built from the same
// target, with the depth counter added.
func TestRunNestedMergesTheEnvironment(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	out, err := system.NewRunner().RunNested(`printf '%s' "$STAGE"`, system.Target{
		Env: map[string]string{"STAGE": "prod"},
	})
	if err != nil {
		t.Fatalf("RunNested: %v", err)
	}

	if got := strings.TrimSpace(out); got != "prod" {
		t.Errorf("STAGE = %q, want prod", got)
	}
}

// A reference to an unset variable fails before the line runs.
func TestRunRejectsAnUnsetEnvReference(t *testing.T) {
	_, err := system.NewRunner().Run("echo hi", system.Target{
		Env: map[string]string{"TOKEN": "${CROWER_TEST_UNSET_XYZ}"},
	})
	if !errors.Is(err, crerrors.ErrUnsetEnvVar) {
		t.Fatalf("error = %v, want ErrUnsetEnvVar", err)
	}
}

// A secret-looking literal is refused; the schema carries a reference instead.
func TestRunRejectsASecretLiteral(t *testing.T) {
	_, err := system.NewRunner().Run("echo hi", system.Target{
		Env: map[string]string{"API_TOKEN": "literal"},
	})
	if !errors.Is(err, crerrors.ErrSecretEnvValue) {
		t.Fatalf("error = %v, want ErrSecretEnvValue", err)
	}
}

func TestRunRejectsAnInvalidTimeout(t *testing.T) {
	_, err := system.NewRunner().Run("echo hi", system.Target{Timeout: "soon"})
	if !errors.Is(err, crerrors.ErrInvalidTimeout) {
		t.Fatalf("error = %v, want ErrInvalidTimeout", err)
	}
}

// A step that overruns its timeout is killed and reported as a timeout, not as
// an anonymous "signal: killed".
func TestRunStopsAtTheTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	start := time.Now()

	_, err := system.NewRunner().Run("sleep 5", system.Target{Timeout: "100ms"})
	if !errors.Is(err, crerrors.ErrTimeout) {
		t.Fatalf("error = %v, want ErrTimeout", err)
	}

	if !strings.Contains(err.Error(), "100ms") {
		t.Errorf("error = %v, want it to name the timeout", err)
	}

	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Errorf("the step ran for %s, want it killed near 100ms", elapsed)
	}
}
