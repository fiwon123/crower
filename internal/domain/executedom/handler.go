package executedom

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/fiwon123/crower/internal/shell"
)

type Handler struct {
	app   *app.Config
	shell *shell.Runner
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{
		app:   app,
		shell: shell.New(),
	}
}

// Execute runs the managed command named or aliased by key.
//
// The command's exec line is copied before parameters are substituted into it.
// It used to be rewritten in place on the command map, so the first run with a
// parameter permanently changed the stored command and a second run in the same
// process substituted into the already-substituted text.
func (h *Handler) Execute(key string, params []string) (string, *command.Data, error) {

	cmd := h.app.AllCommandsByName.Get(key)
	if cmd == nil {
		cmd = h.app.AllCommandsByAlias.Get(key)
	}

	if cmd == nil {
		return "", nil, errors.CommandNotFound(key)
	}

	execLine := applyParams(cmd.Exec, params)

	h.app.Logger.Debug("exec: " + execLine)

	// Get steps
	steps := cmd.Steps
	if len(steps) == 0 {
		steps = []string{execLine}
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
			out, err := executeNested(h.shell, step)
			allOut.WriteString(out)
			if err != nil {
				return allOut.String(), cmd, err
			}
			continue
		}

		h.app.Logger.Debug(fmt.Sprintf("[%d/%d] %s", i+1, len(steps), step))

		out, err := h.shell.Exec(step)
		allOut.WriteString(out)
		if err != nil {
			return allOut.String(), cmd, err
		}
	}

	return allOut.String(), cmd, nil
}

// applyParams fills the placeholders of an exec line, or appends the parameters
// when it has none.
func applyParams(execLine string, params []string) string {
	if len(params) == 0 {
		return execLine
	}

	quoted := make([]string, 0, len(params))
	for _, param := range params {
		quoted = append(quoted, helper.QuoteFor(runtime.GOOS, param))
	}

	if !strings.Contains(execLine, "{") {
		return execLine + " " + strings.Join(quoted, " ")
	}

	for i, param := range quoted {
		execLine = strings.ReplaceAll(execLine, fmt.Sprintf("{%d}", i+1), param)
	}

	// {} is the first parameter.
	return strings.ReplaceAll(execLine, "{}", quoted[0])
}

// executeNested runs a step that calls crower again.
//
// The step is a new process, so the depth has to travel in the environment: an
// in-process counter cannot see the child that does the recursive call, and a
// command whose step calls itself would spawn processes until the machine ran
// out of them. The counter is a limit on that, and the error names the limit
// rather than claiming to have detected a cycle it cannot see.
func executeNested(runner *shell.Runner, step string) (string, error) {
	depth := nestedDepth()
	if depth >= maxExecDepth {
		return "", fmt.Errorf("nested crower call %q: maximum depth %d reached", step, maxExecDepth)
	}

	c, err := runner.Command(step)
	if err != nil {
		return "", err
	}

	c.Env = append(os.Environ(), fmt.Sprintf("%s=%d", nestedDepthEnv, depth+1))

	out, err := c.CombinedOutput()
	return string(out), err
}

// nestedDepth is how many nested crower calls led to this process.
func nestedDepth() int {
	depth, err := strconv.Atoi(os.Getenv(nestedDepthEnv))
	if err != nil {
		return 0
	}

	return depth
}

const (
	// nestedDepthEnv carries the nesting level to the child process.
	nestedDepthEnv = "CROWER_EXEC_DEPTH"

	// maxExecDepth is the deepest chain of nested crower calls that is allowed.
	maxExecDepth = 8
)
