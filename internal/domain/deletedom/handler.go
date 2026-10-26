package deletedom

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/fiwon123/crower/pkg/utils"
)

type Handler struct {
	app     *app.Config
	execute executedom.Handler
}

func NewHandler(app *app.Config) *Handler {

	execute := executedom.NewHandler(app)

	return &Handler{
		app:     app,
		execute: *execute,
	}
}

// Delete command using key
func (h *Handler) DeleteCommand(key string) (*command.Data, bool) {
	command := h.app.AllCommandsByName.Get(key)
	if command == nil {
		command = h.app.AllCommandsByAlias.Get(key)
		if command == nil {
			return nil, false
		}
	}

	h.app.AllCommandsByName.Remove(command.Name)

	for _, alias := range command.AllAlias {
		h.app.AllCommandsByAlias.Remove(alias)
	}

	return command, true
}

// Delete file from filepath
func (h *Handler) DeleteFile(filePath string) error {
	return h.remove(filePath, false)
}

// Delete folder from folderpath
func (h *Handler) DeleteFolder(folderPath string) error {
	return h.remove(folderPath, true)
}

// remove runs the platform's delete command.
func (h *Handler) remove(path string, folder bool) error {
	var commandLine string

	switch runtime.GOOS {
	case "windows":
		if folder {
			commandLine = "rmdir /s /q " + helper.QuoteFor(runtime.GOOS, path)
		} else {
			commandLine = "del " + helper.QuoteFor(runtime.GOOS, path)
		}
	case "linux", "darwin":
		commandLine = "rm "
		if folder {
			commandLine = "rm -r "
		}
		commandLine += helper.QuoteFor(runtime.GOOS, path)
	default:
		return errors.ErrUnsupportedOS
	}

	out, err := h.execute.PerformExecute(commandLine)
	if err != nil {
		return fmt.Errorf("delete %s: %w (output: %s)", path, err, strings.TrimSpace(out))
	}

	return nil
}

// DeleteHistoryContent drops one history entry and the backup file it points
// at.
func (h *Handler) DeleteHistoryContent(content history.Content) (string, error) {
	ok := h.app.History.RemoveContent(content)
	if !ok {
		return "", fmt.Errorf("history entry %s not found", content.File)
	}

	newDataPath := filepath.Join(h.app.HistoryFolderPath, content.File)
	if err := utils.DeleteFile(newDataPath); err != nil {
		return "", fmt.Errorf("delete backup %s: %w", newDataPath, err)
	}

	return "history entry deleted", nil
}
