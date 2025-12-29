package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

func Extract(paths []string, outDir string, app *app.Data) {

	for _, f := range paths {

		if outDir == "" {
			outDir = filepath.Dir(f)
		}

		base := filepath.Base(f)
		split := strings.Split(base, ".")
		ex := ""

		if len(split) > 0 {
			ex = split[len(split)-1]
		}

		if ex == "" {
			continue
		}

		out, err := performExtract(ex, f, outDir)
		if err != nil {
			fmt.Printf("out %s , error %v \n", string(out), err)
			continue
		}

		fmt.Printf("result: %s \n", out)
	}

}

func performExtract(ext string, filePath string, outDir string) ([]byte, error) {
	switch ext {
	case "tar":
		return PerformExecute(fmt.Sprintf("'tar -xf %s -C %s'", filePath, outDir))
	case "gz":
		return PerformExecute(fmt.Sprintf("'tar -xzf %s -C %s'", filePath, outDir))
	case "tgz":
		return PerformExecute(fmt.Sprintf("'tar -xzf %s -C %s'", filePath, outDir))
	case "bz2":
		return PerformExecute(fmt.Sprintf("'tar -xjf %s -C %s'", filePath, outDir))
	case "xz":
		return PerformExecute(fmt.Sprintf("'tar -xJf %s -C %s'", filePath, outDir))
	case "zip":
		switch runtime.GOOS {
		case "windows":
			return PerformExecute(fmt.Sprintf("'tar -xf %s -C %s'", filePath, outDir))
		case "linux":
			return PerformExecute(fmt.Sprintf("'unzip %s -d  %s'", filePath, outDir))
		}
	case "7z":
		return PerformExecute(fmt.Sprintf("'7z x %s -o%s'", filePath, outDir))
	case "rar":
		return PerformExecute(fmt.Sprintf("'7z x %s -o%s'", filePath, outDir))
	}

	return nil, fmt.Errorf("failed to extract \n")
}
