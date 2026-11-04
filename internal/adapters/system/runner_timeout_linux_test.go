//go:build linux

package system_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/fiwon123/crower/internal/adapters/system"
	crerrors "github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// A timeout has to kill the step, not just the shell that started it: killing
// the shell alone leaves `sleep` holding the output pipe (so the call hangs)
// and running after crower has moved on.
//
// The shell backgrounds sleep, records its pid, and waits; the pid must be gone
// when the timed run returns.
func TestTimeoutKillsTheWholeProcessGroup(t *testing.T) {
	pidFile := filepath.Join(t.TempDir(), "child.pid")
	line := "sleep 5 & echo $! > " + shell.Quote(pidFile) + "; wait"

	_, err := system.NewRunner().Run(line, system.Target{Timeout: "100ms"})
	if !errors.Is(err, crerrors.ErrTimeout) {
		t.Fatalf("error = %v, want ErrTimeout", err)
	}

	body, err := os.ReadFile(pidFile)
	if err != nil {
		t.Fatalf("the shell never recorded its child: %v", err)
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(body)))
	if err != nil {
		t.Fatalf("pid %q: %v", body, err)
	}

	deadline := time.Now().Add(2 * time.Second)

	for {
		if processGone(pid) {
			return
		}

		if time.Now().After(deadline) {
			_ = syscall.Kill(pid, syscall.SIGKILL)
			t.Fatalf("child %d survived the timeout", pid)
		}

		time.Sleep(10 * time.Millisecond)
	}
}

// processGone reports whether the pid has exited. A zombie counts as gone: the
// process was killed, and this container's PID 1 does not reap orphans.
func processGone(pid int) bool {
	body, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return true
	}

	// /proc/pid/stat is "pid (comm) state ..."; comm may contain spaces, so
	// take the fields after the closing parenthesis.
	_, rest, ok := strings.Cut(string(body), ")")
	if !ok {
		return true
	}

	fields := strings.Fields(rest)
	if len(fields) == 0 {
		return true
	}

	return fields[0] == "Z" || fields[0] == "X"
}
