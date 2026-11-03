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

	// Get steps
	steps := cmd.Steps
	if len(steps) == 0 {
		execLine, err := applyParams(cmd.Exec, params)
		if err != nil {
			return "", cmd, err
		}

		h.app.Logger.Debug("exec: " + execLine)

		steps = []string{execLine}
	} else {
		substituted, err := applyParamsToSteps(steps, params)
		if err != nil {
			return "", cmd, err
		}

		steps = substituted
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

// paramToken is one brace token in a line, with the bounds of its text.
type paramToken struct {
	start int
	end   int
	token string
}

// paramTokens finds the brace tokens in a line. A "{" whose span contains
// another "{" is literal (a shell snippet before a real placeholder), so it is
// skipped rather than swallowing the placeholder that follows.
func paramTokens(line string) []paramToken {
	var tokens []paramToken

	for i := 0; i < len(line); {
		if line[i] != '{' {
			i++
			continue
		}

		relative := strings.IndexByte(line[i+1:], '}')
		if relative < 0 {
			break
		}

		end := i + 1 + relative
		token := line[i+1 : end]
		if strings.ContainsRune(token, '{') {
			i++
			continue
		}

		tokens = append(tokens, paramToken{start: i, end: end, token: token})
		i = end + 1
	}

	return tokens
}

// applyParams fills the placeholders of an exec line, or appends the parameters
// when the line has no braces at all.
//
// The line is scanned once: substituting sequentially re-scanned the text that
// had already been inserted, so a parameter containing "{2}" was replaced again
// by the second parameter. A brace that is not a placeholder (`{a,b}`, `{0}`,
// an index past the parameters) is left as written, and a line with braces but
// no usable placeholder is an error — the old code returned the line unchanged
// and dropped the parameters without a word.
//
// `{}` takes the parameters no `{n}` consumed, as the docs say; with no `{n}`
// that is every parameter. A `find … -exec … {}` template is only affected when
// parameters are passed explicitly.
func applyParams(execLine string, params []string) (string, error) {
	if len(params) == 0 {
		return execLine, nil
	}

	quoted := quoteParams(params)

	if !strings.Contains(execLine, "{") {
		return execLine + " " + strings.Join(quoted, " "), nil
	}

	out, found, outOfRange := substituteParams(execLine, quoted)
	if found {
		return out, nil
	}

	if outOfRange != "" {
		return "", fmt.Errorf("parameter {%s} is out of range: %d given", outOfRange, len(params))
	}

	return "", fmt.Errorf("no parameter placeholder in %q", execLine)
}

// applyParamsToSteps fills the placeholders of every step. A step without
// braces is left alone (appending to every step of a multi-step command would
// repeat the arguments). A set of steps with no usable placeholder, or with an
// index past the parameters, is an error: in a set of one-placeholder steps a
// partial substitution would otherwise reach the shell as a literal {n}.
func applyParamsToSteps(steps []string, params []string) ([]string, error) {
	if len(params) == 0 {
		return steps, nil
	}

	quoted := quoteParams(params)
	out := append([]string(nil), steps...)
	found := false
	outOfRange := ""

	for i, step := range out {
		if !strings.Contains(step, "{") {
			continue
		}

		substituted, ok, rangeError := substituteParams(step, quoted)
		if rangeError != "" && outOfRange == "" {
			outOfRange = rangeError
		}

		out[i] = substituted

		if ok {
			found = true
		}
	}

	if outOfRange != "" {
		return nil, fmt.Errorf("parameter {%s} is out of range: %d given", outOfRange, len(params))
	}

	if found {
		return out, nil
	}

	return nil, fmt.Errorf("no parameter placeholder in this command's steps")
}

func quoteParams(params []string) []string {
	quoted := make([]string, 0, len(params))
	for _, param := range params {
		quoted = append(quoted, helper.QuoteFor(runtime.GOOS, param))
	}

	return quoted
}

// substituteParams fills the placeholders of one line, reporting whether any
// usable placeholder was found and the first out-of-range index (empty when
// there is none). Leftovers stay as written.
func substituteParams(line string, quoted []string) (string, bool, string) {
	tokens := paramTokens(line)
	used := make([]bool, len(quoted))
	found := false
	outOfRange := ""

	for _, token := range tokens {
		switch {
		case token.token == "":
			found = true
		case !isParamIndex(token.token):
		default:
			n, err := strconv.Atoi(token.token)
			if err != nil || n < 1 || n > len(quoted) {
				if outOfRange == "" {
					outOfRange = token.token
				}
				continue
			}

			used[n-1] = true
			found = true
		}
	}

	remaining := make([]string, 0, len(quoted))
	for i, param := range quoted {
		if !used[i] {
			remaining = append(remaining, param)
		}
	}

	remainder := strings.Join(remaining, " ")

	var out strings.Builder
	cursor := 0

	for _, token := range tokens {
		out.WriteString(line[cursor:token.start])

		switch {
		case token.token == "":
			out.WriteString(remainder)
		case !isParamIndex(token.token):
			out.WriteString(line[token.start : token.end+1])
		default:
			n, err := strconv.Atoi(token.token)
			if err != nil || n < 1 || n > len(quoted) {
				out.WriteString(line[token.start : token.end+1])
			} else {
				out.WriteString(quoted[n-1])
			}
		}

		cursor = token.end + 1
	}

	out.WriteString(line[cursor:])

	return out.String(), found, outOfRange
}

// isParamIndex reports whether token is a positive decimal index.
func isParamIndex(token string) bool {
	if token == "" {
		return false
	}

	for _, r := range token {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
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
