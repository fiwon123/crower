package copydom

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

// Copy file from a origin filepath to output folder
func (h *Handler) CopyFile(filePath string, destFolder string) error {
	return h.copy(filePath, destFolder, false)
}

// Copy folder from a origin folderpath to output folder
func (h *Handler) CopyFolder(folderPath string, destFolder string) error {
	return h.copy(folderPath, destFolder, true)
}

// copy runs the platform's copy command.
//
// The command line is plain shell text. It used to be wrapped in a pair of
// double quotes for a tokenizer that has since been removed, and the shell then
// read the whole line as the name of a single command: `cp` was never run.
func (h *Handler) copy(source string, destFolder string, folder bool) error {
	var commandLine string

	switch runtime.GOOS {
	case "windows":
		if folder {
			commandLine = "xcopy " + helper.QuoteFor(runtime.GOOS, source) + " " + helper.QuoteFor(runtime.GOOS, destFolder) + " /E /I"
		} else {
			commandLine = "copy " + helper.QuoteFor(runtime.GOOS, source) + " " + helper.QuoteFor(runtime.GOOS, destFolder)
		}
	case "linux", "darwin":
		commandLine = "cp "
		if folder {
			commandLine += "-r "
		}
		commandLine += helper.QuoteFor(runtime.GOOS, source) + " " + helper.QuoteFor(runtime.GOOS, destFolder)
	default:
		return errors.ErrUnsupportedOS
	}

	out, err := h.execute.PerformExecute(commandLine)
	if err != nil {
		return fmt.Errorf("copy %s: %w (output: %s)", source, err, strings.TrimSpace(out))
	}

	return nil
}
