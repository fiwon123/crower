package shell

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	crerrors "github.com/fiwon123/crower/internal/errors"
)

func TestCommandForPlatform(t *testing.T) {
	tests := []struct {
		name string
		goos string
		want []string
	}{
		{"windows", "windows", []string{"cmd", "/c", "echo hello"}},
		{"linux", "linux", []string{"sh", "-c", "echo hello"}},
		{"darwin", "darwin", []string{"sh", "-c", "echo hello"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := commandFor(tt.goos, "echo hello")
			if err != nil {
				t.Fatalf("commandFor(%q) error = %v", tt.goos, err)
			}

			if got, want := strings.Join(cmd.Args, "|"), strings.Join(tt.want, "|"); got != want {
				t.Errorf("args = %q, want %q", got, want)
			}
		})
	}
}

// A command may name its shell; the name is portable, so the same value picks
// the same interpreter on every platform.
func TestCommandForNamedShell(t *testing.T) {
	tests := []struct {
		name  string
		goos  string
		shell string
		want  []string
	}{
		{"sh on linux", "linux", "sh", []string{"sh", "-c", "echo hello"}},
		{"sh on windows", "windows", "sh", []string{"sh", "-c", "echo hello"}},
		{"pwsh on linux", "linux", "pwsh", []string{"pwsh", "-NoProfile", "-Command", "echo hello"}},
		{"pwsh on darwin", "darwin", "pwsh", []string{"pwsh", "-NoProfile", "-Command", "echo hello"}},
		{"cmd on linux", "linux", "cmd", []string{"cmd", "/c", "echo hello"}},
		{"cmd on darwin", "darwin", "cmd", []string{"cmd", "/c", "echo hello"}},
		{"empty on windows is the platform default", "windows", "", []string{"cmd", "/c", "echo hello"}},
		{"empty on linux is the platform default", "linux", "", []string{"sh", "-c", "echo hello"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := commandForShell(tt.goos, tt.shell, "echo hello")
			if err != nil {
				t.Fatalf("commandForShell(%q, %q) error = %v", tt.goos, tt.shell, err)
			}

			if got, want := strings.Join(cmd.Args, "|"), strings.Join(tt.want, "|"); got != want {
				t.Errorf("args = %q, want %q", got, want)
			}
		})
	}
}

// A shell name the schema does not define is an error, not a fallback: running
// through a different interpreter than the command asked for would hide the
// mistake.
func TestCommandForUnknownShell(t *testing.T) {
	if _, err := commandForShell("linux", "fish", "echo hello"); err == nil {
		t.Fatal("commandForShell accepted an unknown shell")
	} else if !errors.Is(err, crerrors.ErrUnknownShell) {
		t.Errorf("error = %v, want ErrUnknownShell", err)
	}
}

// An unsupported platform must be an error: a shell invocation that does not
// exist cannot look like a command that succeeded.
func TestCommandForUnsupportedPlatform(t *testing.T) {
	if _, err := commandFor("plan9", "echo hello"); err == nil {
		t.Fatal("commandFor accepted an unsupported platform")
	} else if !errors.Is(err, crerrors.ErrUnsupportedOS) {
		t.Errorf("error = %v, want ErrUnsupportedOS", err)
	}
}

// The command line goes to the shell as one argument.
//
// It used to be split on whitespace and handed over as separate arguments after
// -c, and the shell reads everything after the command string as positional
// parameters: `echo hello` ran `echo` with no arguments, and `ls /nonexistent`
// listed the working directory instead of failing.
func TestExecKeepsTheArguments(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	out, err := New().Exec("echo hello")
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}

	if got := strings.TrimSpace(out); got != "hello" {
		t.Errorf("Exec(\"echo hello\") = %q, want \"hello\"", got)
	}
}

// Quoting survives: a parameter with a space is one argument, not two.
func TestExecKeepsQuotedValues(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	out, err := New().Exec(`printf '%s' "a b"`)
	if err != nil {
		t.Fatalf("Exec() error = %v", err)
	}

	if got := strings.TrimSpace(out); got != "a b" {
		t.Errorf("Exec() = %q, want \"a b\"", got)
	}
}

// A command that fails has to say so. The failure used to be swallowed, so
// `crower execute` exited 0 while nothing ran.
func TestExecReportsAFailingCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	if _, err := New().Exec("exit 3"); err == nil {
		t.Fatal("Exec(\"exit 3\") returned no error")
	}

	if _, err := New().Exec("ls /nonexistent-dir-xyz"); err == nil {
		t.Fatal("Exec of a command that lists a missing folder returned no error")
	}
}

// Start has to build the same command Exec does, and must not wait for it: a
// browser or a file manager outlives the command that opened it.
func TestStartDoesNotWait(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	if err := New().Start("true"); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
}
