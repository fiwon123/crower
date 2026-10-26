package executedom

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/helper"
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
func (h *Handler) Execute(key string, params []string) (string, *command.Data, error) {

	command := h.app.AllCommandsByName.Get(key)
	if command == nil {
		command = h.app.AllCommandsByAlias.Get(key)
	}

	if command == nil {
		return "", nil, fmt.Errorf("command not found")
	}

	if len(params) > 0 {
		// If command has placeholders {1}, {2} etc, substitute them
		if strings.Contains(command.Exec, "{") {
			for i, param := range params {
				placeholder := fmt.Sprintf("{%d}", i+1)
				command.Exec = strings.ReplaceAll(command.Exec, placeholder, helper.QuoteFor(runtime.GOOS, param))
			}
			// Also replace {} with first param if used
			if len(params) > 0 {
				command.Exec = strings.ReplaceAll(command.Exec, "{}", helper.QuoteFor(runtime.GOOS, params[0]))
			}
		} else {
			// Append params as quoted arguments
			for _, param := range params {
				command.Exec += " " + helper.QuoteFor(runtime.GOOS, param)
			}
		}
	}

	h.app.Logger.Info(command.Exec)

	// Get steps
	steps := command.Steps
	if len(steps) == 0 {
		steps = []string{command.Exec}
	}

	// Execute steps sequentially
	var allOut strings.Builder
	for i, step := range steps {
		step = strings.TrimSpace(step)
		if step == "" {
			continue
		}

		// Check for nested crower invocation
		if strings.HasPrefix(step, "crower ") || strings.HasPrefix(step, "cr ") {
			out, err := h.executeNested(step)
			allOut.WriteString(out)
			if err != nil {
				return allOut.String(), command, err
			}
			continue
		}

		h.app.Logger.Info(fmt.Sprintf("[%d/%d] %s", i+1, len(steps), step))
		out, err := h.PerformExecute(step)
		allOut.WriteString(out)
		if err != nil {
			return allOut.String(), command, err
		}
	}
	return allOut.String(), command, nil
}

func (h *Handler) executeNested(step string) (string, error) {
	execDepth++
	defer func() { execDepth-- }()

	if execDepth > maxExecDepth {
		return "", fmt.Errorf("maximum recursion depth (%d) exceeded", maxExecDepth)
	}

	// Parse args
	fields := strings.Fields(step)
	if len(fields) == 0 {
		return "", fmt.Errorf("empty command")
	}

	cmdName := fields[1]
	if len(fields) > 1 {
		// Check for cycles
		for _, s := range execStack {
			if s == cmdName {
				return "", fmt.Errorf("circular dependency detected: %s", cmdName)
			}
		}
		execStack = append(execStack, cmdName)
		defer func() { execStack = execStack[:len(execStack)-1] }()
	}

	// For now, execute via shell - full in-process dispatch would require refactoring
	// to avoid import cycles; this is a simplified nested execution
	return h.PerformExecute(step)
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

var (
	execDepth    int
	execStack    []string
	maxExecDepth = 8
)
