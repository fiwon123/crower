package createdom

import (
	"fmt"
	"path/filepath"

	"runtime"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/app"

	commanddata "github.com/fiwon123/crower/internal/data/command"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
	"github.com/fiwon123/crower/pkg/crowerutils"
)

type CreateHandler struct {
	app *app.Config
}

func NewCreateHandler(app *app.Config) *CreateHandler {
	return &CreateHandler{
		app: app,
	}
}

// Create command using name, alias and exec parameters
func (h *CreateHandler) CreateCommand(name string, alias []string, exec string) (*commanddata.Data, error) {

	command := commanddata.New(name, alias, exec)

	if command.Name == "" {
		return nil, fmt.Errorf("empty name")
	}

	if command.Exec == "" {
		return nil, fmt.Errorf("empty exec")
	}

	if h.app.AllCommandsByName.Get(command.Name) != nil {
		return nil, fmt.Errorf("found name, command already added")
	}

	for _, alias := range command.AllAlias {
		if h.app.AllCommandsByAlias.Get(alias) != nil || h.app.AllCommandsByName.Get(alias) != nil {
			return nil, fmt.Errorf("found alias, command already added")
		}
	}

	h.app.AllCommandsByName.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		h.app.AllCommandsByAlias.Add(alias, command)
	}

	return command, nil
}

// Create command based on process name or id process
func (h *CreateHandler) CreateProcess(name string, args []string) (*commanddata.Data, error) {
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
		command, err := h.CreateCommand(name, nil, execCommand)

		if err != nil {
			return nil, err
		}

		return command, nil
	} else if pathStr != "" {
		pathStr = fmt.Sprintf("'%s'", pathStr)
		command, err := h.CreateCommand(name, nil, pathStr)

		if err != nil {
			return nil, err
		}

		return command, nil
	}

	return nil, fmt.Errorf("couldn't find the process either by pid or name")
}

// Create a new file on filepath
func (h *CreateHandler) CreateFile(filePath string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("type nul > '%s'", filePath))
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"touch '%s'\"", filePath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}

// Create a new folder on folderpath
func (h *CreateHandler) CreateFolder(folderPath string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("mkdir '%s'", folderPath))
	case "linux":
		out, err = executehandlers.PerformExecute(fmt.Sprintf("\"mkdir '%s'\"", folderPath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	h.app.Logger.Info("output: ", "out", out)
	return nil
}

func (h *CreateHandler) CreateScriptCommand(name string) (string, error) {
	cfgFolderPath := filepath.Dir(h.app.CfgFilePath)
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
