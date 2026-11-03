package copydom

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/shell"
)

type Handler struct {
	app   *app.Config
	shell *shell.Runner
}

func NewHandler(app *app.Config) *Handler {

	return &Handler{
		app:   app,
		shell: shell.New(),
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
	commandLine, err := system.CopyCommand(runtime.GOOS, source, destFolder, folder)
	if err != nil {
		return err
	}

	out, err := h.shell.Exec(commandLine)
	if err != nil {
		return fmt.Errorf("copy %s: %w (output: %s)", source, err, strings.TrimSpace(out))
	}

	return nil
}
