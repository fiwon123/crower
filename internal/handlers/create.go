package handlers

import (
	"fmt"
	"path/filepath"
	"slices"

	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/pkg/utils"
)

// Create command using name, alias and exec parameters
func CreateCommand(name string, alias []string, exec string, app *app.Data) (*command.Data, error) {

	command := command.New(name, alias, exec)

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

func CreateSystemVariable(newVar string, value string, app *app.Data) (string, error) {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("reg", "add", `HKCU\Environment`, "/v", newVar, "/t", "REG_SZ", "/d", value, "/f")
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("Create System Variable error: %v \n", err)
		}

		return "System variable set. You may need to restart or log off for it to take effect.", nil
	case "linux":
		f, err := os.OpenFile(os.Getenv("HOME")+"/.bashrc", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = fmt.Fprintf(f, "\nexport %s=%s\n", newVar, value)
		if err != nil {
			return "", fmt.Errorf("Create System Variable error: %v \n", err)
		}

		return "Added to .bashrc. Restart terminal to take effect.", nil
	}

	return "", fmt.Errorf("Couldn't find OS to excute command")
}

func checkNewVarValue(value string, from string) bool {
	splitted := splitPath(from)
	ok := true
	if slices.Contains(splitted, value) {
		return false
	}

	return ok
}

func splitPath(path string) []string {
	return filepath.SplitList(path)
}

// Create a new file on filepath
func CreateFile(filePath string, app *app.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("type nul > '%s'", filePath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("\"touch '%s'\"", filePath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", out)
	return nil
}

// Create a new folder on folderpath
func CreateFolder(folderPath string, app *app.Data) error {
	var out string
	var err error
	switch runtime.GOOS {
	case "windows":
		out, err = PerformExecute(fmt.Sprintf("mkdir '%s'", folderPath))
	case "linux":
		out, err = PerformExecute(fmt.Sprintf("\"mkdir '%s'\"", folderPath))
	}

	if err != nil {
		return fmt.Errorf("out %s, error %v\n", out, err)
	}

	fmt.Println("result: ", out)
	return nil
}
