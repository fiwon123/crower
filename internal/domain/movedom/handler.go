package movedom

import (
	"fmt"
	"runtime"

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

// Move file from origin path to output folder path
func (h *Handler) MoveFile(filePath string, destFolder string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute(fmt.Sprintf("move '%s' '%s'", filePath, destFolder))
	case "linux":
		out, err = h.execute.PerformExecute(fmt.Sprintf("\"mv '%s' '%s'\"", filePath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}

// Move folder from origin path to output folder path
func (h *Handler) MoveFolder(folderPath string, destFolder string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute(fmt.Sprintf("move '%s' '%s'", folderPath, destFolder))
	case "linux":
		out, err = h.execute.PerformExecute(fmt.Sprintf("\"mv '%s' '%s'\"", folderPath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}
