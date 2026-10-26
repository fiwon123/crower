package opendom

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/helper"
)

type Handler struct {
	app     *app.Config
	execute *executedom.Handler
	run     CommandRunner
	edit    Editor
}

// Editor opens a file in the user's editor, which needs the terminal itself.
// It is a field for the same reason run is: a test has no editor and no
// terminal, so a test that used the real one would report on the machine.
type Editor func(file string) error

// CommandRunner executes a command string and returns its combined output.
//
// It is a field so a test can replace the execution. The real runner shells
// out to xdg-open, which needs a desktop environment: on a headless machine
// it exits 3 ("no method available for opening ..."), so a test that used it
// would be reporting on the machine rather than on the code.
type CommandRunner func(command string) (string, error)

func NewHandler(app *app.Config) *Handler {

	execute := executedom.NewHandler(app)

	h := &Handler{
		app:     app,
		execute: execute,
	}

	h.run = execute.PerformExecute
	h.edit = defaultEditor(execute)

	return h
}

// WithEditor replaces the editor OpenSystem uses and returns the handler, so a
// test can record the file that would have been opened.
func (h *Handler) WithEditor(edit Editor) *Handler {
	h.edit = edit
	return h
}

// NewHandlerWithRunner returns a Handler that executes commands with the
// given runner instead of running them.
//
// A test uses it to avoid the desktop dependency. The command is still built
// by the production code, so the test can assert on what would have run.
func NewHandlerWithRunner(app *app.Config, run CommandRunner) *Handler {

	h := NewHandler(app)
	h.run = run

	return h
}

// Open opens every path with the default application.
//
// Each path is opened; the first failure stops the batch and is returned. The
// previous loop returned after the first path, so `crower open a b` silently
// opened only a.
func (h *Handler) Open(paths []string) error {
	for _, f := range paths {
		commandString := ""

		switch runtime.GOOS {
		case "windows":
			commandString = fmt.Sprintf(`start "" %s`, helper.QuoteFor(runtime.GOOS, f))
		case "linux":
			commandString = "xdg-open " + helper.QuoteFor(runtime.GOOS, f)
		case "darwin":
			commandString = "open " + helper.QuoteFor(runtime.GOOS, f)
		}

		if commandString == "" {
			return fmt.Errorf("opening files is not supported on %s", runtime.GOOS)
		}

		h.app.Logger.Info("opening: ", commandString)

		out, err := h.run(commandString)
		if err != nil {
			return fmt.Errorf("open %s: %w (output: %s)", f, err, strings.TrimSpace(out))
		}

		if trimmed := strings.TrimSpace(out); trimmed != "" {
			h.app.Logger.Out(trimmed)
		}
	}

	return nil
}

// OpenSystem opens the system variable editor of the running platform.
func (h *Handler) OpenSystem() error {
	switch runtime.GOOS {
	case "windows":
		return h.execute.PerformExecuteStart("sysdm.cpl")
	case "linux", "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		return h.edit(filepath.Join(home, shellProfileFileName()))
	}

	return fmt.Errorf("opening system variables is not supported on %s", runtime.GOOS)
}

// shellProfileFileName is the file a shell reads on startup.
func shellProfileFileName() string {
	if runtime.GOOS == "darwin" {
		return ".zshrc"
	}

	return ".bashrc"
}

// defaultEditor returns an Editor that runs $EDITOR, or one of the usual
// editors, in the terminal.
func defaultEditor(execute *executedom.Handler) Editor {
	return func(file string) error {
		name := editorName()
		if name == "" {
			return errors.New("no editor found: set $EDITOR or install nano or vi")
		}

		return execute.PerformInteractiveTerminal(name, file)
	}
}

// editorName is $EDITOR when it is set, otherwise the first editor on the path.
func editorName() string {
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
