package handlers

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/fiwon123/crower/internal/data"
	"github.com/google/shlex"
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
		for _, param := range args {
			command.Exec += param
		}
	}

	fmt.Println(command.Exec)
	out, err := PerformExecute(command.Exec)
	app.LoggerInfo.Info(command.Exec)
	return out, err
}

func PerformExecute(ex string) ([]byte, error) {
	var c *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		ex = fmt.Sprintf("/c %s", ex)

		splitExec, err := getSplitCommand(ex)
		if err != nil {
			return nil, err
		}

		c = exec.Command("cmd", splitExec...)
	case "linux":
		ex = fmt.Sprintf("-c %s", ex)

		splitExec, err := getSplitCommand(ex)
		if err != nil {
			return nil, err
		}

		c = exec.Command("sh", splitExec...)
	}

	out, err := c.CombinedOutput()
	return out, err
}

func getSplitCommand(ex string) ([]string, error) {
	var out []string

	splitExec, err := shlex.Split(ex)
	if err != nil {
		return nil, err
	}
	for i, _ := range splitExec {
		out = append(out, splitExec[i])
	}

	return out, nil
}
