package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
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
		return PerformExecute(fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir), app)
	case "gz":
		return PerformExecute(fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir), app)
	case "tgz":
		return PerformExecute(fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir), app)
	case "bz2":
		return PerformExecute(fmt.Sprintf(`tar -xjf '%s' -C '%s'`, filePath, outDir), app)
	case "xz":
		return PerformExecute(fmt.Sprintf(`tar -xJf '%s' -C '%s'`, filePath, outDir), app)
	case "zip":
		switch runtime.GOOS {
		case "windows":
			return PerformExecute(fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir), app)
		case "linux":
			return PerformExecute(fmt.Sprintf(`unzip '%s' -d  '%s'`, filePath, outDir), app)
		}
	case "7z":
		return PerformExecute(fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir), app)
	case "rar":
		return PerformExecute(fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir), app)
	}

	return "", fmt.Errorf("failed to extract \n")
}
