package createhandlers

import (
	"fmt"
	"path/filepath"

	"runtime"
	"strconv"
	"strings"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

// Create command using name, alias and exec parameters
func CreateCommand(name string, alias []string, exec string, app *appdata.Data) (*commanddata.Data, error) {

	command := commanddata.New(name, alias, exec)

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
func CreateProcess(name string, args []string, app *appdata.Data) (*commanddata.Data, error) {
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
		pathStr, err = crowerutils.GetProcessPathByName(processName)
		if err != nil {
			return nil, err
		}
	} else {

		pathStr, err = crowerutils.GetProcessPathByID(int32(pid))
		if err != nil {
			return nil, err
		}

	}

	if strings.Contains(pathStr, "app/") {
		if processName == "" {
			processName, err = crowerutils.GetProcessNameByID(int32(pid))
			if err != nil {
				return nil, err
			}
		}

		var appID string
		appID, err = crowerutils.GetFlatpakAppIDByName(processName)
		if err != nil {
			return nil, err
		}

		execCommand := fmt.Sprintf("flatpak run %s", appID)
		command, err := CreateCommand(name, nil, execCommand, app)

		if err != nil {
			return nil, err
		}

		return command, nil
	} else if pathStr != "" {
		pathStr = fmt.Sprintf("'%s'", pathStr)
		command, err := CreateCommand(name, nil, pathStr, app)

		if err != nil {
			return nil, err
		}

		return command, nil
	}

	return nil, fmt.Errorf("couldn't find the process either by pid or name")
}

// Create a new file on filepath
func CreateFile(filePath string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("type nul > '%s'", filePath), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"touch '%s'\"", filePath), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}

// Create a new folder on folderpath
func CreateFolder(folderPath string, app *appdata.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("mkdir '%s'", folderPath), app)
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"mkdir '%s'\"", folderPath), app)
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	app.Logger.Info("output: ", "out", out)
	return nil
}

func CreateScriptCommand(name string, app *appdata.Data) (string, error) {
	cfgFolderPath := filepath.Dir(app.CfgFilePath)
	scriptFolderPath := filepath.Join(cfgFolderPath, "scripts")
	err := crowerutils.CreateFolderIfNotExists(scriptFolderPath)
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		scriptFilePath := filepath.Join(scriptFolderPath, name+".bat")

		crowerutils.CreateFileIfNotExists(scriptFilePath)

		return scriptFilePath, nil
	case "linux":
		scriptFilePath := filepath.Join(scriptFolderPath, name+".sh")

		crowerutils.CreateFileIfNotExists(scriptFilePath)

		return scriptFilePath, nil
	}

	return "", fmt.Errorf("can't find specific OS to create script command")
}
