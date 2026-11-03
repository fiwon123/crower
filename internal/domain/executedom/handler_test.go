package executedom_test

import (
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/executedom"
)

// newTestApp builds an application whose streams are discarded, so a test does
// not write into the terminal running it.
func newTestApp(t *testing.T) *app.Config {
	t.Helper()

	cfg, err := app.InitApp(filepath.Join(t.TempDir(), "crower.yaml"), app.Options{
		Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("InitApp() error = %v", err)
	}

	return cfg
}

// The command line goes to the shell as one argument.
//
// It used to be split on whitespace and handed over as separate arguments after
// -c, and the shell reads everything after the command string as positional
// parameters: `echo hello` ran `echo` with no arguments, and `ls /nonexistent`
// listed the working directory instead of failing.
func TestPerformExecuteKeepsTheArguments(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	out, err := executedom.NewHandler(newTestApp(t)).PerformExecute("echo hello")
	if err != nil {
		t.Fatalf("PerformExecute() error = %v", err)
	}

	if got := strings.TrimSpace(out); got != "hello" {
		t.Errorf("PerformExecute(\"echo hello\") = %q, want \"hello\"", got)
	}
}

// A command that fails has to say so. The failure used to be swallowed, so
// `crower execute` exited 0 while nothing ran.
func TestPerformExecuteReportsAFailingCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	if _, err := executedom.NewHandler(newTestApp(t)).PerformExecute("exit 3"); err == nil {
		t.Fatal("PerformExecute(\"exit 3\") returned no error")
	}

	if _, err := executedom.NewHandler(newTestApp(t)).PerformExecute("ls /nonexistent-dir-xyz"); err == nil {
		t.Fatal("PerformExecute of a command that lists a missing folder returned no error")
	}
}

func TestExecuteRunsTheStoredCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	cfg := newTestApp(t)
	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo hello"))

	out, executed, err := executedom.NewHandler(cfg).Execute("greet", nil)
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got := strings.TrimSpace(out); got != "hello" {
		t.Errorf("Execute() output = %q, want \"hello\"", got)
	}

	if executed == nil || executed.Name != "greet" {
		t.Errorf("Execute() returned command %v, want greet", executed)
	}
}

func TestExecuteFindsACommandByAlias(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	cfg := newTestApp(t)
	greet := command.New("greet", []string{"hi"}, "echo hello")
	cfg.AllCommandsByName.Add("greet", greet)
	cfg.AllCommandsByAlias.Add("hi", greet)

	if _, _, err := executedom.NewHandler(cfg).Execute("hi", nil); err != nil {
		t.Fatalf("Execute() by alias error = %v", err)
	}
}

// Parameters are substituted into a copy of the command.
//
// The substitution used to rewrite the command stored on the map, so the first
// run with a parameter permanently changed the command every later run would
// see.
func TestExecuteDoesNotRewriteTheStoredCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test drives a POSIX shell")
	}

	cfg := newTestApp(t)
	cfg.AllCommandsByName.Add("greet", command.New("greet", nil, "echo {1}"))

	handler := executedom.NewHandler(cfg)

	out, _, err := handler.Execute("greet", []string{"first"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got := strings.TrimSpace(out); got != "first" {
		t.Errorf("Execute() output = %q, want \"first\"", got)
	}

	out, _, err = handler.Execute("greet", []string{"second"})
	if err != nil {
		t.Fatalf("second Execute() error = %v", err)
	}

	if got := strings.TrimSpace(out); got != "second" {
		t.Errorf("second Execute() output = %q, want \"second\"", got)
	}

	if stored := cfg.AllCommandsByName.Get("greet"); stored.Exec != "echo {1}" {
		t.Errorf("stored command = %q, want it unchanged", stored.Exec)
	}
}

// A command that does not exist is an error, not a crash. The lookup used to
// return a nil command and the note builder dereferenced it, which the shell
// reported as signal: segmentation fault.
func TestExecuteReportsAnUnknownCommand(t *testing.T) {
	_, _, err := executedom.NewHandler(newTestApp(t)).Execute("no-such-command", nil)
	if err == nil {
		t.Fatal("Execute() of an unknown command returned no error")
	}

	if !strings.Contains(err.Error(), "no-such-command") {
		t.Errorf("error %v does not name the command", err)
	}
}
