package createdom

import (
	"fmt"
	"path/filepath"

	"runtime"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/fiwon123/crower/pkg/utils"
)

type Handler struct {
	app     *app.Config
	execute executedom.Handler
}

func NewHandler(app *app.Config) *Handler {

	execute := executedom.NewHandler(app)

	return &Handler{
		app:     app,
		execute: *execute,
	}
}

// Create command using name, alias and exec parameters
func (h *Handler) CreateCommand(name string, alias []string, exec string) (*command.Data, error) {

	command := command.New(name, alias, exec)

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
func (h *Handler) CreateProcess(name string, args []string) (*command.Data, error) {
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
func (h *Handler) CreateFile(filePath string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute("type nul > " + helper.QuoteFor(runtime.GOOS, filePath))
	case "linux", "darwin":
		out, err = h.execute.PerformExecute("touch " + helper.QuoteFor(runtime.GOOS, filePath))
	default:
		return errors.ErrUnsupportedOS
	}

	if err != nil {
		return fmt.Errorf("create file %s: %w (output: %s)", filePath, err, strings.TrimSpace(out))
	}

	return nil
}

// Create a new folder on folderpath
func (h *Handler) CreateFolder(folderPath string) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = h.execute.PerformExecute("mkdir " + helper.QuoteFor(runtime.GOOS, folderPath))
	case "linux", "darwin":
		out, err = h.execute.PerformExecute("mkdir -p " + helper.QuoteFor(runtime.GOOS, folderPath))
	default:
		return errors.ErrUnsupportedOS
	}

	if err != nil {
		return fmt.Errorf("create folder %s: %w (output: %s)", folderPath, err, strings.TrimSpace(out))
	}

	return nil
}

func (h *Handler) CreateScriptCommand(name string) (string, error) {
	cfgFolderPath := filepath.Dir(h.app.CfgFilePath)
	scriptFolderPath := filepath.Join(cfgFolderPath, "scripts")
	err := utils.CreateFolderIfNotExists(scriptFolderPath)
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		scriptFilePath := filepath.Join(scriptFolderPath, name+".bat")

		if err := utils.CreateFileIfNotExists(scriptFilePath); err != nil {
			return "", err
		}

		return scriptFilePath, nil
	case "linux", "darwin":
		scriptFilePath := filepath.Join(scriptFolderPath, name+".sh")

		if err := utils.CreateFileIfNotExists(scriptFilePath); err != nil {
			return "", err
		}

		return scriptFilePath, nil
	}

	return "", fmt.Errorf("can't find specific OS to create script command")
}
