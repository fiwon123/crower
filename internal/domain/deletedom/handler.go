package deletedom

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/shell"

	"github.com/fiwon123/crower/pkg/utils"
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
	commandLine, err := removeCommand(runtime.GOOS, path, folder)
	if err != nil {
		return err
	}

	out, err := h.shell.Exec(commandLine)
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
