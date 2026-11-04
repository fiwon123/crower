package system

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fiwon123/crower/internal/shell"
)

// Runner runs shell lines on the host.
//
// It exists so the execute use case never starts a process itself: the shell
// invocation, and the depth guard around a self-invocation, are technology.
type Runner struct {
	shell *shell.Runner
}

// NewRunner builds a Runner that runs on the host.
func NewRunner() *Runner {
	return &Runner{shell: shell.New()}
}

// Run runs one line and returns its combined output.
func (r *Runner) Run(line string) (string, error) {
	return r.shell.Exec(line)
}

// RunNested runs a step that calls crower again.
//
// The step is a new process, so the depth has to travel in the environment: an
// in-process counter cannot see the child that does the recursive call, and a
// command whose step calls itself would spawn processes until the machine ran
// out of them. The counter is a limit on that, and the error names the limit
// rather than claiming to have detected a cycle it cannot see.
func (r *Runner) RunNested(line string) (string, error) {
	depth := nestedDepth()
	if depth >= maxExecDepth {
		return "", fmt.Errorf("nested crower call %q: maximum depth %d reached", line, maxExecDepth)
	}

	c, err := r.shell.Command(line)
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
