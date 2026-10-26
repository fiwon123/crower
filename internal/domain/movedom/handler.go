package movedom

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
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

// Move file from origin path to output folder path
func (h *Handler) MoveFile(filePath string, destFolder string) error {
	return h.move(filePath, destFolder)
}

// Move folder from origin path to output folder path
func (h *Handler) MoveFolder(folderPath string, destFolder string) error {
	return h.move(folderPath, destFolder)
}

// move runs the platform's move command.
//
// The POSIX command line used to be wrapped in double quotes, which the shell
// read as one very long command name.
func (h *Handler) move(source string, destFolder string) error {
	var commandLine string

	switch runtime.GOOS {
	case "windows":
		commandLine = "move " + helper.QuoteFor(runtime.GOOS, source) + " " + helper.QuoteFor(runtime.GOOS, destFolder)
	case "linux", "darwin":
		commandLine = "mv " + helper.QuoteFor(runtime.GOOS, source) + " " + helper.QuoteFor(runtime.GOOS, destFolder)
	default:
		return errors.ErrUnsupportedOS
	}

	out, err := h.execute.PerformExecute(commandLine)
	if err != nil {
		return fmt.Errorf("move %s: %w (output: %s)", source, err, strings.TrimSpace(out))
	}

	return nil
}
