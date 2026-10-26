package extractdom

import (
	stderrors "errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

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

// Extract expands every archive in paths.
//
// Every archive is attempted; if any fails, the failures are returned together
// so the command exits non-zero. The previous version logged a failure and
// carried on, which a script read as success.
func (h *Handler) Extract(paths []string, output string) error {
	var failures []error

	for _, f := range paths {

		if output == "" {
			output = filepath.Dir(f)
		}

		base := filepath.Base(f)
		split := strings.Split(base, ".")
		ext := ""

		if len(split) > 0 {
			ext = split[len(split)-1]
		}

		if ext == "" {
			continue
		}

		out, err := h.performExtract(ext, f, output)
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", f, err))
			continue
		}

		h.app.Logger.Out(strings.TrimRight(string(out), "\n"))
	}

	return stderrors.Join(failures...)
}

func (h *Handler) performExtract(ext string, filePath string, outDir string) (string, error) {
	switch ext {
	case "tar":
		return h.execute.PerformExecute(fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir))
	case "gz":
		return h.execute.PerformExecute(fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir))
	case "tgz":
		return h.execute.PerformExecute(fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir))
	case "bz2":
		return h.execute.PerformExecute(fmt.Sprintf(`tar -xjf '%s' -C '%s'`, filePath, outDir))
	case "xz":
		return h.execute.PerformExecute(fmt.Sprintf(`tar -xJf '%s' -C '%s'`, filePath, outDir))
	case "zip":
		switch runtime.GOOS {
		case "windows":
			return h.execute.PerformExecute(fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir))
		case "linux":
			return h.execute.PerformExecute(fmt.Sprintf(`unzip '%s' -d  '%s'`, filePath, outDir))
		}
	case "7z":
		return h.execute.PerformExecute(fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir))
	case "rar":
		return h.execute.PerformExecute(fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir))
	}

	return "", fmt.Errorf("no extractor for .%s archives", ext)
}
