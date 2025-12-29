package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
)

// Delete command from the cfg file.
func DeleteCommand(key string, app *app.Data) (*command.Data, bool) {
	command := app.AllCommandsByName.Get(key)
	if command == nil {
		command = app.AllCommandsByAlias.Get(key)
		if command == nil {
			return nil, false
		}
	}

	app.AllCommandsByName.Remove(command.Name)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Remove(alias)
	}

	return command, true
}

func DeleteFile(filePath string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("del '%s'", filePath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("rm '%s'", filePath))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))
}

func DeleteFolder(folderPath string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("rmdir /s /q '%s'", folderPath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("rm -r '%s'", folderPath))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))
}
