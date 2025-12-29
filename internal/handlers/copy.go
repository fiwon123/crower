package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

func CopyFile(filePath string, destFolder string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("'copy %s %s'", filePath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("'cp %s %s'", filePath, destFolder))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))

}

func CopyFolder(filePath string, destFolder string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("'xcopy %s %s /E /I'", filePath, destFolder))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("'cp -r %s %s'", filePath, destFolder))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))
}
