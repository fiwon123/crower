package system

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// Runner runs shell lines on the host.
//
// It exists so the execute use case never starts a process itself: the shell
// invocation, the working directory, and the depth guard around a
// self-invocation are technology.
type Runner struct {
	shell *shell.Runner
}

// NewRunner builds a Runner that runs on the host.
func NewRunner() *Runner {
	return &Runner{shell: shell.New()}
}

// Run runs one line and returns its combined output.
//
// shellName names the shell the command asked for; empty is the platform
// default. dir is the working directory; empty is the process's own. Both are
// validated before a process starts, so an unknown shell or a missing folder
// fails without running the line.
func (r *Runner) Run(line, shellName, dir string) (string, error) {
	c, err := r.command(line, shellName, dir)
	if err != nil {
		return "", err
	}

	out, err := c.CombinedOutput()

	return string(out), err
}

// RunNested runs a step that calls crower again.
//
// The step is a new process, so the depth has to travel in the environment: an
// in-process counter cannot see the child that does the recursive call, and a
// command whose step calls itself would spawn processes until the machine ran
// out of them. The counter is a limit on that, and the error names the limit
// rather than claiming to have detected a cycle it cannot see.
func (r *Runner) RunNested(line, shellName, dir string) (string, error) {
	depth := nestedDepth()
	if depth >= maxExecDepth {
		return "", fmt.Errorf("nested crower call %q: maximum depth %d reached", line, maxExecDepth)
	}

	c, err := r.command(line, shellName, dir)
	if err != nil {
		return "", err
	}

	c.Env = append(os.Environ(), fmt.Sprintf("%s=%d", nestedDepthEnv, depth+1))

	out, err := c.CombinedOutput()

	return string(out), err
}

// command builds the process for one line, rejecting a target that cannot run.
func (r *Runner) command(line, shellName, dir string) (*exec.Cmd, error) {
	c, err := r.shell.CommandFor(shellName, line)
	if err != nil {
		return nil, err
	}

	resolved, err := resolveDir(dir)
	if err != nil {
		return nil, err
	}

	c.Dir = resolved

	return c, nil
}

// resolveDir expands "~" and refuses a directory that does not exist.
//
// The check is done here, before the process starts, so a command with a typo
// in Cwd fails with a message about the folder rather than an opaque chdir
// error from the shell. It proves existence, not access: a directory whose
// permissions deny entry still fails inside exec, after nothing has run.
func resolveDir(dir string) (string, error) {
	if dir == "" {
		return "", nil
	}

	expanded, err := expandHome(dir)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(expanded)
	if err != nil {
		// Keep the cause (permission, symlink loop) reachable with errors.Is
		// while still classifying the failure as an invalid working directory.
		return "", fmt.Errorf("%w: %q: %w", errors.ErrInvalidWorkDir, dir, err)
	}

	if !info.IsDir() {
		return "", fmt.Errorf("%w: %q is not a directory", errors.ErrInvalidWorkDir, dir)
	}

	return expanded, nil
}

// expandHome resolves a leading "~" to the user's home directory.
//
// Only the bare "~" and "~/..." (or "~\..." on Windows) forms are expanded;
// "~user/..." is left as written, because resolving another user's home needs
// platform-specific lookup and no command needs it today.
func expandHome(dir string) (string, error) {
	if dir == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w: %q: %w", errors.ErrInvalidWorkDir, dir, err)
		}

		return home, nil
	}

	if strings.HasPrefix(dir, "~/") || strings.HasPrefix(dir, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%w: %q: %w", errors.ErrInvalidWorkDir, dir, err)
		}

		return filepath.Join(home, dir[2:]), nil
	}

	return dir, nil
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
