package extractdom

import (
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

// Extract compressed files in output folder path
func (h *Handler) Extract(paths []string, output string) {

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
			h.app.Logger.Error("failed", "out", string(out), "err", err)
			continue
		}

		h.app.Logger.Info("output: ", "out", string(out))
	}

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

	return "", fmt.Errorf("failed to extract \n")
}
