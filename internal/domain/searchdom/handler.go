package searchdom

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/domain/updatedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/pkg/utils"
)

// unsupportedOSError reports a search that has no implementation for the
// running platform. It is an error rather than an empty result: a search that
// matched nothing and a search that cannot run must not look the same.
func unsupportedOSError(goos string) error {
	return fmt.Errorf("%w: search %q", errors.ErrUnsupportedOS, goos)
}

type CommandRunner interface {
	PerformExecuteStart(ex string) error
	PerformExecute(ex string) (string, error)
}

type Handler struct {
	app            *app.Config
	handler        *updatedom.Handler
	executeHandler *executedom.Handler
	runner         CommandRunner
}

func NewHandler(app *app.Config) *Handler {
	executeHandler := executedom.NewHandler(app)
	handler := updatedom.NewHandler(app)
	return &Handler{
		app:            app,
		handler:        handler,
		executeHandler: executeHandler,
		runner:         executeHandler,
	}
}

func NewHandlerWithRunner(app *app.Config, runner CommandRunner) *Handler {
	executeHandler := executedom.NewHandler(app)
	handler := updatedom.NewHandler(app)
	return &Handler{
		app:            app,
		handler:        handler,
		executeHandler: executeHandler,
		runner:         runner,
	}
}

// SearchBrowser opens a search in the default browser.
func (h *Handler) SearchBrowser(content string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		if utils.IsURL(content) {
			return "", h.runner.PerformExecuteStart(fmt.Sprintf(`start ' ' '%s' `, content))
		}
		return "", h.runner.PerformExecuteStart(fmt.Sprintf(`start ' ' 'https://duckduckgo.com/?q=%s' `, content))
	case "linux":
		if utils.IsURL(content) {
			return "", h.runner.PerformExecuteStart(fmt.Sprintf(`xdg-open "%s" `, content))
		}
		return "", h.runner.PerformExecuteStart(fmt.Sprintf(`xdg-open "https://duckduckgo.com/?q=%s" `, content))
	case "darwin":
		if utils.IsURL(content) {
			return "", h.runner.PerformExecuteStart(fmt.Sprintf(`open "%s"`, content))
		}
		return "", h.runner.PerformExecuteStart(fmt.Sprintf(`open "https://duckduckgo.com/?q=%s"`, content))
	}

	return "", unsupportedOSError(runtime.GOOS)
}

// SearchFile searches files under currentPath.
func (h *Handler) SearchFile(currentPath string, content string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`dir "%s" /s /b /a-d | findstr "%s"`, currentPath, content))
	case "linux", "darwin":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`find "%s" -type f -name "%s"`, currentPath, content))
	}

	return "", unsupportedOSError(runtime.GOOS)
}

// SearchFolder searches folders under currentPath.
func (h *Handler) SearchFolder(currentPath string, content string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`dir "%s" /s /b /ad | findstr "%s"`, currentPath, content))
	case "linux", "darwin":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`find "%s" -type d -name "%s"`, currentPath, content))
	}

	return "", unsupportedOSError(runtime.GOOS)
}

// SearchFileAndFolder searches files and folders under currentPath.
func (h *Handler) SearchFileAndFolder(currentPath string, content string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`dir "%s" /s /b | findstr "%s"`, currentPath, content))
	case "linux", "darwin":
		return h.executeHandler.PerformExecute(fmt.Sprintf(`find "%s" -name "%s"`, currentPath, content))
	}

	return "", unsupportedOSError(runtime.GOOS)
}
