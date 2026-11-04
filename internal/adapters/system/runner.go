package system

import (
	"context"
	"fmt"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
)

// Target is where and how one line runs: the shell, the working directory, the
// environment the command adds, and the per-step timeout.
type Target struct {
	// Shell names the shell ("sh", "pwsh", "cmd"); empty is the platform's own.
	Shell string

	// Dir is the working directory; empty is the process's own.
	Dir string

	// Env is merged over the process environment. A "${VAR}" value reads the
	// process environment; a literal under a secret-looking key is refused.
	Env map[string]string

	// Timeout is a Go duration ("120s"); empty means no timeout.
	Timeout string
}

// Runner runs shell lines on the host.
//
// It exists so the execute use case never starts a process itself: the shell
// invocation, the working directory, the environment, the timeout, and the
// depth guard around a self-invocation are technology.
type Runner struct {
	shell *shell.Runner
}

// NewRunner builds a Runner that runs on the host.
func NewRunner() *Runner {
	return &Runner{shell: shell.New()}
}

// Run runs one line and returns its combined output.
//
// The target is validated before a process starts, so an unknown shell, a
// missing folder, a secret literal, or a malformed timeout fails without
// running the line.
func (r *Runner) Run(line string, target Target) (string, error) {
	ctx, cancel, err := timeoutContext(target.Timeout)
	if err != nil {
		return "", err
	}
	defer cancel()

	c, err := r.command(ctx, line, target)
	if err != nil {
		return "", err
	}

	out, err := c.CombinedOutput()

	return string(out), timedOut(ctx, err, target.Timeout)
}

// RunNested runs a step that calls crower again.
//
// The step is a new process, so the depth has to travel in the environment: an
// in-process counter cannot see the child that does the recursive call, and a
// command whose step calls itself would spawn processes until the machine ran
// out of them. The counter is a limit on that, and the error names the limit
// rather than claiming to have detected a cycle it cannot see.
func (r *Runner) RunNested(line string, target Target) (string, error) {
	depth := nestedDepth()
	if depth >= maxExecDepth {
		return "", fmt.Errorf("nested crower call %q: maximum depth %d reached", line, maxExecDepth)
	}

	ctx, cancel, err := timeoutContext(target.Timeout)
	if err != nil {
		return "", err
	}
	defer cancel()

	// The depth is an override, not an appended entry: a value the command
	// declares for the same key must not shadow the guard, and the child sees
	// exactly one copy of it.
	overrides := maps.Clone(target.Env)
	if overrides == nil {
		overrides = map[string]string{}
	}

	overrides[nestedDepthEnv] = strconv.Itoa(depth + 1)
	target.Env = overrides

	c, err := r.command(ctx, line, target)
	if err != nil {
		return "", err
	}

	out, err := c.CombinedOutput()

	return string(out), timedOut(ctx, err, target.Timeout)
}

// command builds the process for one line, rejecting a target that cannot run.
func (r *Runner) command(ctx context.Context, line string, target Target) (*exec.Cmd, error) {
	c, err := r.shell.CommandForContext(ctx, target.Shell, line)
	if err != nil {
		return nil, err
	}

	// Only a command that can time out needs its own process group; an
	// ordinary command keeps the terminal's group and its Ctrl+C.
	if target.Timeout != "" {
		shell.Isolate(c)
	}

	resolved, err := resolveDir(target.Dir)
	if err != nil {
		return nil, err
	}

	env, err := buildEnv(target.Env)
	if err != nil {
		return nil, err
	}

	c.Dir = resolved
	c.Env = env

	return c, nil
}

// timeoutContext turns the declared timeout into a context. Empty means none.
func timeoutContext(raw string) (context.Context, context.CancelFunc, error) {
	if raw == "" {
		return context.Background(), func() {}, nil
	}

	duration, err := time.ParseDuration(raw)
	if err != nil || duration <= 0 {
		return nil, nil, fmt.Errorf("%w: %q", errors.ErrInvalidTimeout, raw)
	}

	ctx, cancel := context.WithTimeout(context.Background(), duration)

	return ctx, cancel, nil
}

// timedOut turns the context's deadline into an error that names the timeout.
//
// exec reports a killed process as an exit error ("signal: killed"), which says
// nothing about why; the context is the only place the reason survives.
func timedOut(ctx context.Context, err error, raw string) error {
	if err != nil && ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("%w after %s", errors.ErrTimeout, raw)
	}

	return err
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

// buildEnv merges the command's entries over the process environment.
//
// A map is used so an override replaces the process value instead of appearing
// twice, and the result is sorted so the child environment is deterministic.
func buildEnv(overrides map[string]string) ([]string, error) {
	values := map[string]string{}

	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			values[key] = value
		}
	}

	// Sorted so a config with two bad entries always reports the same one.
	for _, key := range slices.Sorted(maps.Keys(overrides)) {
		if key == "" || strings.ContainsAny(key, "=\x00") {
			return nil, fmt.Errorf("%w: key %q", errors.ErrInvalidEnv, key)
		}

		resolved, err := resolveEnvValue(key, overrides[key])
		if err != nil {
			return nil, err
		}

		setEnvValue(values, key, resolved)
	}

	keys := slices.Sorted(maps.Keys(values))

	env := make([]string, 0, len(keys))
	for _, key := range keys {
		env = append(env, key+"="+values[key])
	}

	return env, nil
}

// setEnvValue writes an override, matching an existing key case-insensitively
// on Windows so the child does not receive both "Path" and "PATH".
func setEnvValue(values map[string]string, key, value string) {
	if runtime.GOOS == "windows" {
		for existing := range values {
			if strings.EqualFold(existing, key) {
				values[existing] = value
				return
			}
		}
	}

	values[key] = value
}

// resolveEnvValue reads a "${VAR}" reference or accepts a literal.
//
// A literal under a key that looks like a secret is refused: the schema says a
// pack carries placeholders, and a pack is a file people commit and share. The
// check is a key-name heuristic, not a secret scanner -- it catches the common
// names, and it is documented as such rather than pretending to be complete.
func resolveEnvValue(key, value string) (string, error) {
	if strings.ContainsRune(value, 0) {
		return "", fmt.Errorf("%w: %s", errors.ErrInvalidEnv, key)
	}

	if name, ok := envReference(value); ok {
		resolved, found := os.LookupEnv(name)
		if !found {
			return "", fmt.Errorf("%w: %s references %s", errors.ErrUnsetEnvVar, key, name)
		}

		return resolved, nil
	}

	if looksSecret(key) {
		return "", fmt.Errorf("%w: %s; use ${VAR}", errors.ErrSecretEnvValue, key)
	}

	return value, nil
}

// envReference reports whether the whole value is a "${VAR}" reference.
//
// Only the exact form is a reference: a value that merely contains "${VAR}"
// stays literal, so a command cannot smuggle a resolved secret into a value the
// secret check would otherwise see as text.
func envReference(value string) (string, bool) {
	if !strings.HasPrefix(value, "${") || !strings.HasSuffix(value, "}") {
		return "", false
	}

	name := value[2 : len(value)-1]
	if name == "" {
		return "", false
	}

	for i, r := range name {
		switch {
		case r == '_', r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return "", false
		}
	}

	return name, true
}

// looksSecret reports whether a key names a credential by convention.
func looksSecret(key string) bool {
	upper := strings.ToUpper(key)

	for _, marker := range []string{"TOKEN", "SECRET", "PASSWORD", "PASSWD", "CREDENTIAL"} {
		if strings.Contains(upper, marker) {
			return true
		}
	}

	return upper == "KEY" || strings.HasSuffix(upper, "_KEY")
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
