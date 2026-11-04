//go:build windows

package shell

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The real runner has to reach cmd and come back. Every local gate is
// string-level, so this is the only proof that `cmd /c` sees the line the
// quoting helpers built.
func TestWindowsExecRunsThroughCmd(t *testing.T) {
	out, err := New().Exec("echo hello")
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	if got := strings.TrimSpace(out); got != "hello" {
		t.Errorf("Exec(\"echo hello\") = %q, want \"hello\"", got)
	}
}

// A path with a space survives both quoting layers: QuoteFor wraps it, and cmd
// has to hand the one argument to `type`.
func TestWindowsExecKeepsAPathWithASpace(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a folder")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	file := filepath.Join(dir, "note.txt")

	const body = "hello from a path with a space"
	if err := os.WriteFile(file, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	out, err := New().Exec("type " + QuoteFor("windows", file))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	if !strings.Contains(out, body) {
		t.Errorf("output = %q, want %q", out, body)
	}
}

// A quoted "&" is data, not a command separator. If the quoting is lost, cmd
// runs the second command and creates the marker.
func TestWindowsExecKeepsAQuotedAmpersand(t *testing.T) {
	marker := filepath.Join(t.TempDir(), "marker.txt")
	value := "a & type nul > " + marker

	out, err := New().Exec("echo " + QuoteFor("windows", value))
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}

	if _, err := os.Stat(marker); err == nil {
		t.Fatalf("the unquoted & ran a second command and created %s", marker)
	}

	if !strings.Contains(out, "a & type nul >") {
		t.Errorf("output = %q, want the value intact", out)
	}
}

// TestWindowsArgHelper is the target of the embedded-quote test. An ordinary
// run has no "--" marker and does nothing; the runner invocation passes the
// value after it.
func TestWindowsArgHelper(t *testing.T) {
	for i, arg := range os.Args {
		if arg != "--" {
			continue
		}

		_, _ = os.Stdout.WriteString(strings.Join(os.Args[i+1:], "\x00"))
		os.Exit(0)
	}
}

// A value containing a quote and an ampersand reaches the program intact. The
// test binary is the program, so the check is on its argv, not on echo output.
func TestWindowsExecKeepsAnEmbeddedQuote(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("os.Executable: %v", err)
	}

	const value = `he said "hi" & more`

	line := QuoteFor("windows", exe) +
		" -test.run=TestWindowsArgHelper -- " +
		QuoteFor("windows", value)

	out, err := New().Exec(line)
	if err != nil {
		t.Fatalf("Exec(%q): %v", line, err)
	}

	if out != value {
		t.Errorf("arg = %q, want %q (line %q)", out, value, line)
	}
}
