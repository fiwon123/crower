package opendom

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
)

type Handler struct {
	app     *app.Config
	execute *executedom.Handler
	run     CommandRunner
}

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

// Open filepath based on user operational system(OS).
func (h *Handler) Open(paths []string) error {

	for _, f := range paths {
		commandString := ""

		var fstring strings.Builder
		fstring.WriteString(f)

		switch runtime.GOOS {
		case "windows":
			commandString = fmt.Sprintf(`start ' ' '%s'`, fstring.String())
		case "linux":
			commandString = fmt.Sprintf(`"xdg-open '%s'"`, fstring.String())
		}

		if commandString == "" {
			continue
		}

		h.app.Logger.Info("performing execute...: ", "exec", commandString)

		out, err := h.run(commandString)
		if err != nil {
			return fmt.Errorf("error %v out %v", err, string(out))
		}

		h.app.Logger.Info(string(out))
		return nil
	}

	return nil
}

// Try to open system UI based on operational system (OS)
func (h *Handler) OpenSystem() error {
	switch runtime.GOOS {
	case "windows":
		return h.execute.PerformExecuteStart("sysdm.cpl")
	case "linux":
		h.execute.PerformInteractiveTerminal("nano", "~/.bashrc")
	}

	return nil
}
