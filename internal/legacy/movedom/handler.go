package movedom

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
	commandLine, err := system.MoveCommand(runtime.GOOS, source, destFolder)
	if err != nil {
		return err
	}

	out, err := h.shell.Exec(commandLine)
	if err != nil {
		return fmt.Errorf("move %s: %w (output: %s)", source, err, strings.TrimSpace(out))
	}

	return nil
}
