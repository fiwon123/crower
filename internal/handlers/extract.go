package handlers

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

func Extract(folderPath string, fileName string, outDir string, app *app.Data) {

	currentPath := filepath.Join(folderPath, fileName)
	multipleFiles := false
	ext := ""
	split := strings.Split(fileName, ".")
	if len(split) > 0 {
		if split[0] == "*" {
			multipleFiles = true
		}

		ext = split[len(split)-1]
	}

	if multipleFiles {
		files, err := filepath.Glob(currentPath)
		if err != nil {
			fmt.Println("failed invalid filename")
			return
		}

		for _, f := range files {
			performExtract(ext, f, outDir)
		}

		return

	}

	performExtract(ext, currentPath, outDir)

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
