package handlers

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data"
)

// Execute command based on the user operational system (OS).
func Execute(name string, args []string, app *data.App) ([]byte, error) {

	if name == "" && len(args) > 0 {
		fmt.Println("args", args)
		name = args[0]
		args = args[1:]
	}

	command := app.AllCommandsByName.Get(name)
	if command == nil {
		command = app.AllCommandsByAlias.Get(name)
	}

	if command == nil {
		return nil, fmt.Errorf("command not found")
	}

	if len(args) > 0 {
		command.Exec += "\""
		for _, param := range args {
			command.Exec += param
		}
		command.Exec += "\""
	}

	fmt.Println(command.Exec)
	out, err := PerformExecute(command.Exec)
	app.LoggerInfo.Info(command.Exec)
	return out, err
}

func PerformExecute(ex string) ([]byte, error) {
	ex = strings.TrimSuffix(ex, `\`)
	var c *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		c = exec.Command("cmd", "/c", ex)
	case "linux":
		c = exec.Command("sh", "-c", ex)
	}

	out, err := c.Output()
	return out, err
}
