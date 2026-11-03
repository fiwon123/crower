package searchdom

import (
	"runtime"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/shell"
)

type CommandRunner interface {
	Start(ex string) error
	Exec(ex string) (string, error)
}

type Handler struct {
	app    *app.Config
	runner CommandRunner
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{
		app:    app,
		runner: shell.New(),
	}
}

func NewHandlerWithRunner(app *app.Config, runner CommandRunner) *Handler {
	return &Handler{
		app:    app,
		runner: runner,
	}
}

// SearchBrowser opens a search in the default browser.
func (h *Handler) SearchBrowser(content string) (string, error) {
	commandLine, err := browserCommand(runtime.GOOS, content)
	if err != nil {
		return "", err
	}

	return "", h.runner.Start(commandLine)
}

// SearchFile searches files under currentPath.
func (h *Handler) SearchFile(currentPath string, content string) (string, error) {
	commandLine, err := searchFileCommand(runtime.GOOS, currentPath, content)
	if err != nil {
		return "", err
	}

	return h.runner.Exec(commandLine)
}

// SearchFolder searches folders under currentPath.
func (h *Handler) SearchFolder(currentPath string, content string) (string, error) {
	commandLine, err := searchFolderCommand(runtime.GOOS, currentPath, content)
	if err != nil {
		return "", err
	}

	return h.runner.Exec(commandLine)
}

// SearchFileAndFolder searches files and folders under currentPath.
func (h *Handler) SearchFileAndFolder(currentPath string, content string) (string, error) {
	commandLine, err := searchFileAndFolderCommand(runtime.GOOS, currentPath, content)
	if err != nil {
		return "", err
	}

	return h.runner.Exec(commandLine)
}
