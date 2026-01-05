package handlers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
)

// Execute command based on the user operational system (OS).
// Verify if command exists by name or alias and perform operation
func Execute(name string, args []string, app *app.Data) (string, *command.Data, error) {

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
		return "", nil, fmt.Errorf("command not found")
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

func buildCMD(ex string) (*exec.Cmd, string, []string) {
	fmt.Println()
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

	fmt.Printf("Executing... : %s \n", commandString.String())

	return c, typeCommand, splitCommands
}

// Perform execute operation
func PerformExecute(ex string) (string, error) {

	c, typeCommand, splitCommands := buildCMD(ex)
	c = exec.Command(typeCommand, splitCommands...)
	out, err := c.CombinedOutput()
	return string(out), err
}

func PerformExecuteStart(ex string) error {
	c, typeCommand, splitCommands := buildCMD(ex)
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
