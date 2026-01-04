package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

// Copy file from a origin filepath to output folder
func CopyFile(filePath string, destFolder string, app *app.Data) error {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("'copy %s %s'", filePath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("'cp %s %s'", filePath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", string(out))
	return nil
}

// Copy file from a origin folderpath to output folder
func CopyFolder(filePath string, destFolder string, app *app.Data) error {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("'xcopy %s %s /E /I'", filePath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("'cp -r %s %s'", filePath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", string(out))
	return nil
}
