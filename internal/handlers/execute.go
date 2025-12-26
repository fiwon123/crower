package handlers

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
)

// Execute command based on the user operational system (OS).
func Execute(name string, args []string, app *app.Data) ([]byte, *command.Data, error) {

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
		return nil, nil, fmt.Errorf("command not found")
	}

	if len(args) > 0 {
		for _, param := range args {
			command.Exec += param
		}
	}

	fmt.Println(command.Exec)
	out, err := PerformExecute(command.Exec)
	app.LoggerInfo.Info(command.Exec)
	return out, command, err
}

func PerformExecute(ex string) ([]byte, error) {
	var c *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		ex = fmt.Sprintf("/c %s", ex)
		c = exec.Command("cmd", getSplitCommand(ex)...)
	case "linux":
		ex = fmt.Sprintf("-c %s", ex)
		c = exec.Command("sh", getSplitCommand(ex)...)
	}

	out, err := c.CombinedOutput()
	return out, err
}

func getSplitCommand(ex string) []string {
	tokenRe := regexp.MustCompile(`"([^"]+)"|'([^']+)'|([^\s]+)`)
	matches := tokenRe.FindAllStringSubmatch(ex, -1)

	var args []string
	for _, m := range matches {
		for i := 1; i <= 3; i++ {
			if m[i] != "" {
				args = append(args, m[i])
				break
			}
		}
	}

	return args
}
