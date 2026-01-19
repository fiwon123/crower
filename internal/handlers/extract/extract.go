package extracthandlers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
)

// Extract compressed files in output folder path
func Extract(paths []string, output string, app *app.Data) {

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

		out, err := performExtract(ext, f, output, app)
		if err != nil {
			app.Logger.Error("failed", "out", string(out), "err", err)
			continue
		}

		app.Logger.Info("output: ", "out", string(out))
	}

}

func performExtract(ext string, filePath string, outDir string, app *app.Data) (string, error) {
	switch ext {
	case "tar":
		return executehandlers.PerformExecute(fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir), app)
	case "gz":
		return executehandlers.PerformExecute(fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir), app)
	case "tgz":
		return executehandlers.PerformExecute(fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir), app)
	case "bz2":
		return executehandlers.PerformExecute(fmt.Sprintf(`tar -xjf '%s' -C '%s'`, filePath, outDir), app)
	case "xz":
		return executehandlers.PerformExecute(fmt.Sprintf(`tar -xJf '%s' -C '%s'`, filePath, outDir), app)
	case "zip":
		switch runtime.GOOS {
		case "windows":
			return executehandlers.PerformExecute(fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir), app)
		case "linux":
			return executehandlers.PerformExecute(fmt.Sprintf(`unzip '%s' -d  '%s'`, filePath, outDir), app)
		}
	case "7z":
		return executehandlers.PerformExecute(fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir), app)
	case "rar":
		return executehandlers.PerformExecute(fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir), app)
	}

	return "", fmt.Errorf("failed to extract \n")
}
