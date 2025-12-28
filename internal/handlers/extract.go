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
		base := filepath.Base(f)
		split := strings.Split(base, ".")
		ex := ""

		if len(split) > 0 {
			ex = split[len(split)-1]
		}

		if ex == "" {
			continue
		}

		performExtract(ex, f, outDir)
	}

}

func performExtract(ext string, filePath string, outDir string) {
	switch ext {
	case "tar":
		PerformExecute(fmt.Sprintf("'tar -xf %s -C %s'", filePath, outDir))
	case "gz":
		fmt.Println(filePath)
		fmt.Println(outDir)
		PerformExecute(fmt.Sprintf("'tar -xzf %s -C %s'", filePath, outDir))
	case "tgz":
		PerformExecute(fmt.Sprintf("'tar -xzf %s -C %s'", filePath, outDir))
	case "bz2":
		PerformExecute(fmt.Sprintf("'tar -xjf %s -C %s'", filePath, outDir))
	case "xz":
		PerformExecute(fmt.Sprintf("'tar -xJf %s -C %s'", filePath, outDir))
	case "zip":
		switch runtime.GOOS {
		case "windows":
			PerformExecute(fmt.Sprintf("'tar -xf %s -C %s'", filePath, outDir))
		case "linux":
			fmt.Println("unzip")
			PerformExecute(fmt.Sprintf("'unzip %s -d  %s'", filePath, outDir))
		}
	case "7z":
		PerformExecute(fmt.Sprintf("'7z x %s -o%s'", filePath, outDir))
	case "rar":
		PerformExecute(fmt.Sprintf("'7z x %s -o%s'", filePath, outDir))
	default:
		fmt.Println("failed extract")
		return
	}
}
