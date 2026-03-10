package executedom

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"

	commanddata "github.com/fiwon123/crower/internal/data/command"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{
		app: app,
	}
}

// Execute command based on the user operational system (OS).
// Verify if command exists by name or alias and perform operation
func (h *Handler) Execute(key string, params []string) (string, *commanddata.Data, error) {

	command := h.app.AllCommandsByName.Get(key)
	if command == nil {
		command = h.app.AllCommandsByAlias.Get(key)
	}

	if command == nil {
		return "", nil, fmt.Errorf("command not found")
	}

	if len(params) > 0 {
		for _, param := range params {
			command.Exec += param
		}
	}

	h.app.Logger.Info(command.Exec)
	out, err := h.PerformExecute(command.Exec)
	return out, command, err
}

func (h *Handler) buildCMD(ex string) (*exec.Cmd, string, []string) {
	h.app.Logger.Info("")
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

	h.app.Logger.Info("Executing... : ", "exec", commandString.String())

	return c, typeCommand, splitCommands
}

// Perform execute operation
func (h *Handler) PerformExecute(ex string) (string, error) {

	c, typeCommand, splitCommands := h.buildCMD(ex)
	c = exec.Command(typeCommand, splitCommands...)
	out, err := c.CombinedOutput()
	return string(out), err
}

func (h *Handler) PerformExecuteStart(ex string) error {
	c, typeCommand, splitCommands := h.buildCMD(ex)
	c = exec.Command(typeCommand, splitCommands...)
	err := c.Start()
	return err
}

// Perform operation that needs another terminal
func (h *Handler) PerformInteractiveTerminal(commandName string, ex string) {
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
