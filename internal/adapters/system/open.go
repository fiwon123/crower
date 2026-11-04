package system

import (
	stderrors "errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// Opener opens paths with the platform's default application.
type Opener struct {
	goos string
	run  func(command string) (string, error)
}

// NewOpener builds an Opener that runs on the host.
func NewOpener() *Opener {
	runner := shell.New()

	return &Opener{goos: runtime.GOOS, run: runner.Exec}
}

// Open opens one path.
//
// It returns the command that ran and its trimmed output so the caller can
// report both; the error wraps the output as well, because the platform's
// message is the only explanation of a failure to open.
func (o *Opener) Open(path string) (string, string, error) {
	commandLine, err := OpenCommand(o.goos, path)
	if err != nil {
		return "", "", err
	}

	out, err := o.run(commandLine)
	trimmed := strings.TrimSpace(out)

	if err != nil {
		return commandLine, trimmed, fmt.Errorf("open %s: %w (output: %s)", path, err, trimmed)
	}

	return commandLine, trimmed, nil
}

// OpenCommand builds the platform's "open with the default application"
// invocation.
//
// goos is a parameter, not a read of runtime.GOOS, so a table test can cover
// every platform the release ships (#129).
func OpenCommand(goos, path string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf(`start "" %s`, shell.QuoteFor(goos, path)), nil
	case "linux":
		return "xdg-open " + shell.QuoteFor(goos, path), nil
	case "darwin":
		return "open " + shell.QuoteFor(goos, path), nil
	default:
		return "", fmt.Errorf("%w: open on %q", errors.ErrUnsupportedOS, goos)
	}
}

// SystemEditor opens the platform's environment-variable editor.
type SystemEditor struct {
	goos  string
	start func(command string) error
	edit  func(file string) error
	home  func() (string, error)
}

// NewSystemEditor builds an editor for the host.
func NewSystemEditor() *SystemEditor {
	runner := shell.New()

	return &SystemEditor{
		goos:  runtime.GOOS,
		start: runner.Start,
		edit:  defaultEditor(runner),
		home:  os.UserHomeDir,
	}
}

// Open opens the system-variable editor.
func (e *SystemEditor) Open() error {
	switch e.goos {
	case "windows":
		return e.start("sysdm.cpl")
	case "linux", "darwin":
		home, err := e.home()
		if err != nil {
			return err
		}

		return e.edit(filepath.Join(home, ShellProfileFileName(e.goos)))
	}

	return fmt.Errorf("opening system variables is not supported on %s", e.goos)
}

// ShellProfileFileName is the file a shell reads on startup.
func ShellProfileFileName(goos string) string {
	if goos == "darwin" {
		return ".zshrc"
	}

	return ".bashrc"
}

// defaultEditor returns an editor that runs $EDITOR, or one of the usual
// editors, in the terminal.
//
// EditorName is resolved when the editor is opened, not when this closure is
// built: a test (or a user) can set $EDITOR after the process started.
func defaultEditor(runner *shell.Runner) func(string) error {
	return func(file string) error {
		name := EditorName()
		if name == "" {
			return stderrors.New("no editor found: set $EDITOR or install nano or vi")
		}

		return runner.Interactive(name, file)
	}
}

// EditorName is $EDITOR when it is set, otherwise the first editor on the path.
func EditorName() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	for _, name := range []string{"nano", "vi", "vim"} {
		if _, err := exec.LookPath(name); err == nil {
			return name
		}
	}

	return ""
}
