package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/pkg/utils"
)

// Add command from the cfg file.
func AddCommand(name string, alias []string, exec string, args []string, app *data.App) error {

	command := &data.Command{
		Name:     name,
		AllAlias: alias,
		Exec:     exec,
	}

	if len(args) == 2 {
		command.Name = args[0]
		command.Exec = args[1]
	}

	if command.Name == "" {
		return fmt.Errorf("empty name")
	}

	if command.Exec == "" {
		return fmt.Errorf("empty exec")
	}

	if app.AllCommandsByName.Get(command.Name) != nil {
		return fmt.Errorf("found name, command already added")
	}

	for _, alias := range command.AllAlias {
		if app.AllCommandsByAlias.Get(alias) != nil || app.AllCommandsByName.Get(alias) != nil {
			return fmt.Errorf("found alias, command already added")
		}
	}

	app.AllCommandsByName.Add(command.Name, command)

	for _, alias := range command.AllAlias {
		app.AllCommandsByAlias.Add(alias, command)
	}

	return nil
}

func AddProcess(name string, args []string, app *data.App) error {
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
			return err
		}
	} else {
		pathStr, err = utils.GetProcessPathByID(int32(pid))
		if err != nil {
			return err
		}

	}

	if strings.Contains(pathStr, "app/") {
		if processName == "" {
			processName, err = utils.GetProcessNameByID(int32(pid))
			if err != nil {
				return err
			}
		}

		var appID string
		appID, err = utils.GetFlatpakAppIDByName(processName)
		if err != nil {
			return err
		}

		execCommand := fmt.Sprintf("flatpak run %s", appID)
		AddCommand(name, nil, execCommand, nil, app)

		return nil
	}

	return fmt.Errorf("couldn't find the process either by pid or name")
}
