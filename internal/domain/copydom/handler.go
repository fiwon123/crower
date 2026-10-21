package copydom

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

// Copy file from a origin filepath to output folder
func (h *Handler) CopyFile(filePath string, destFolder string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute(fmt.Sprintf("copy '%s' '%s'", filePath, destFolder))
	case "linux":
		out, err = h.execute.PerformExecute(fmt.Sprintf("\"cp '%s' '%s'\"", filePath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}

// Copy file from a origin folderpath to output folder
func (h *Handler) CopyFolder(filePath string, destFolder string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute(fmt.Sprintf("xcopy '%s' '%s' /E /I", filePath, destFolder))
	case "linux":
		out, err = h.execute.PerformExecute(fmt.Sprintf("\"cp -r '%s' '%s'\"", filePath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}
