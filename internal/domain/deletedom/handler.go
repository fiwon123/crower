package deletedom

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/fiwon123/crower/internal/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	historydata "github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/executedom"

	"github.com/fiwon123/crower/pkg/crowerutils"
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
func (h *Handler) DeleteCommand(key string) (*commanddata.Data, bool) {
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
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute(fmt.Sprintf("del '%s'", filePath))
	case "linux":
		out, err = h.execute.PerformExecute(fmt.Sprintf("\"rm '%s'\"", filePath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}

// Delete folder from folderpath
func (h *Handler) DeleteFolder(folderPath string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute(fmt.Sprintf("rmdir /s /q '%s'", folderPath))
	case "linux":
		out, err = h.execute.PerformExecute(fmt.Sprintf(`"rm -r '%s'"`, folderPath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}

func (h *Handler) DeleteHistoryContent(content historydata.Content) (string, error) {
	ok := h.app.History.RemoveContent(content)
	if !ok {
		return "", fmt.Errorf("Content not found \n")
	}

	newDataPath := filepath.Join(h.app.HistoryFolderPath, content.File)
	err := crowerutils.DeleteFile(newDataPath)
	if err != nil {
		return "", fmt.Errorf("Content not deleted %v \n", err)
	}

	return "Content deleted", nil
}
