package handlers

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fiwon123/crower/internal/data"
)

func Execute(input data.Command, app *data.App) ([]byte, error) {

	command := app.CommandsMap.Get(input.Name)

	if command == nil && len(input.Alias) > 0 {
		fmt.Println("find command by alias ", input.Alias)
		command = app.AliasMap.Get(input.Alias[0])
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	i := strings.IndexByte(command.Exec, ' ')
	cmdName := ""
	args := ""
	if i == -1 {
		cmdName = command.Exec
	} else {
		cmdName = command.Exec[:i]
		args = command.Exec[i+1:]
	}

	c := exec.Command(cmdName, args)
	out, err := c.Output()

	return out, err
}
