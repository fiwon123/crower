package handlers

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/pkg/utils"
)

// Create command using name, alias and exec parameters
func CreateCommand(name string, alias []string, exec string, args []string, app *app.Data) (*command.Data, error) {

	command := command.New(name, alias, exec)

	if len(args) == 2 {
		command.Name = args[0]
		command.Exec = args[1]
	}

	if command.Name == "" {
		return nil, fmt.Errorf("empty name")
	}

	if command.Exec == "" {
		return nil, fmt.Errorf("empty exec")
	}

	if app.AllCommandsByName.Get(command.Name) != nil {
		return nil, fmt.Errorf("found name, command already added")
	}

	for _, alias := range command.AllAlias {
		if app.AllCommandsByAlias.Get(alias) != nil || app.AllCommandsByName.Get(alias) != nil {
			return nil, fmt.Errorf("found alias, command already added")
		}
	}

	app.AllCommandsByName.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Add(alias, command)
	}

	return command, nil
}

// Create command based on process name or id process
func CreateProcess(name string, args []string, app *app.Data) (*command.Data, error) {
	if len(args) > 0 && name == "" {
		name = args[0]
		args = args[1:]
	}

	process := args[0]
	pathStr := ""
	processName := ""
	pid, err := strconv.Atoi(process)
	if err != nil {
		processName = process
		pathStr, err = utils.GetProcessPathByName(processName)
		if err != nil {
			return nil, err
		}
	} else {

		pathStr, err = utils.GetProcessPathByID(int32(pid))
		if err != nil {
			return nil, err
		}

	}

	if strings.Contains(pathStr, "app/") {
		if processName == "" {
			processName, err = utils.GetProcessNameByID(int32(pid))
			if err != nil {
				return nil, err
			}
		}

		var appID string
		appID, err = utils.GetFlatpakAppIDByName(processName)
		if err != nil {
			return nil, err
		}

		execCommand := fmt.Sprintf("flatpak run %s", appID)
		command, err := CreateCommand(name, nil, execCommand, nil, app)

		if err != nil {
			return nil, err
		}

		return command, nil
	} else if pathStr != "" {
		pathStr = fmt.Sprintf("'%s'", pathStr)
		command, err := CreateCommand(name, nil, pathStr, nil, app)

		if err != nil {
			return nil, err
		}

		return command, nil
	}

	return nil, fmt.Errorf("couldn't find the process either by pid or name")
}

// Create a new file on filepath
func CreateFile(filePath string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("type nul > '%s'", filePath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("touch '%s'", filePath))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))
}

// Create a new folder on folderpath
func CreateFolder(folderPath string, app *app.Data) {
	var out []byte
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("mkdir '%s'", folderPath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("mkdir '%s'", folderPath))
	}

	if err != nil {
		fmt.Printf("out %s, error %v\n", out, err)
		return
	}

	fmt.Println("result: ", string(out))
}
