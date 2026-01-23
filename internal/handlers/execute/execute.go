package executehandlers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	appdata "github.com/fiwon123/crower/internal/data/app"
	commanddata "github.com/fiwon123/crower/internal/data/command"
)

// Execute command based on the user operational system (OS).
// Verify if command exists by name or alias and perform operation
func Execute(key string, params []string, app *appdata.Data) (string, *commanddata.Data, error) {

	command := app.AllCommandsByName.Get(key)
	if command == nil {
		command = app.AllCommandsByAlias.Get(key)
	}

	if command == nil {
		return "", nil, fmt.Errorf("command not found")
	}

	if len(params) > 0 {
		for _, param := range params {
			command.Exec += param
		}
	}

	app.Logger.Info(command.Exec)
	out, err := PerformExecute(command.Exec, app)
	return out, command, err
}

func buildCMD(ex string, app *appdata.Data) (*exec.Cmd, string, []string) {
	app.Logger.Info("")
	var c *exec.Cmd

	typeCommand := ""
	splitCommands := []string{}
	switch runtime.GOOS {
	case "windows":
		ex = fmt.Sprintf("/c %s", ex)
		splitCommands = append(splitCommands, getSplitCommand(ex)...)
		typeCommand = "cmd"
	case "linux":
		ex = fmt.Sprintf("-c %s", ex)
		splitCommands = append(splitCommands, getSplitCommand(ex)...)
		typeCommand = "sh"
	}

	var commandString strings.Builder
	commandString.WriteString("[")
	commandString.WriteString("\"")
	commandString.WriteString(typeCommand)
	commandString.WriteString("\"")
	for i := range splitCommands {
		commandString.WriteString(", ")
		commandString.WriteString("\"")
		commandString.WriteString(splitCommands[i])
		commandString.WriteString("\"")
	}
	commandString.WriteString("]")

	app.Logger.Info("Executing... : ", "exec", commandString.String())

	return c, typeCommand, splitCommands
}

// Perform execute operation
func PerformExecute(ex string, app *appdata.Data) (string, error) {

	c, typeCommand, splitCommands := buildCMD(ex, app)
	c = exec.Command(typeCommand, splitCommands...)
	out, err := c.CombinedOutput()
	return string(out), err
}

func PerformExecuteStart(ex string, app *appdata.Data) error {
	c, typeCommand, splitCommands := buildCMD(ex, app)
	c = exec.Command(typeCommand, splitCommands...)
	err := c.Start()
	return err
}

// Perform operation that needs another terminal
func PerformInteractiveTerminal(commandName string, ex string) {
	var c *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		c = exec.Command(commandName, ex)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	}

}

func getSplitCommand(ex string) []string {
	tokenRe := regexp.MustCompile(`"([^"]*)"|'([^']*)'|(\S+)`)
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
