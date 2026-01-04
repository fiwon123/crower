package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

// Move file from origin path to output folder path
func MoveFile(filePath string, destFolder string, app *app.Data) error {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("move '%s' '%s'", filePath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("mv '%s' '%s'", filePath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", string(out))
	return nil
}

// Move folder from origin path to output folder path
func MoveFolder(folderPath string, destFolder string, app *app.Data) error {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("move '%s' '%s'", folderPath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("mv '%s' '%s'", folderPath, destFolder))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", string(out))
	return nil
}
