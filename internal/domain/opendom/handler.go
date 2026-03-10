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
}

func NewHandler(app *app.Config) *Handler {

	execute := executedom.NewHandler(app)

	return &Handler{
		app:     app,
		execute: execute,
	}
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

		out, err := h.execute.PerformExecute(commandString)
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
