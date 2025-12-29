package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

func MoveFile(filePath string, destFolder string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("move '%s' '%s'", filePath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("mv '%s' '%s'", filePath, destFolder))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))
}

func MoveFolder(folderPath string, destFolder string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("move '%s' '%s'", folderPath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("mv '%s' '%s'", folderPath, destFolder))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))

}
