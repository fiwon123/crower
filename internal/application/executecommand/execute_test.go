package executecommand_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/executecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// fakeState records what the use case asked it to store.
type fakeState struct {
	commands  catalog.Catalog
	last      map[operation.Operation]string
	records   []record
	recordErr error
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) Commands() catalog.Catalog {
	return s.commands
}

func (s *fakeState) LastCommand(op operation.Operation) (string, bool) {
	name, ok := s.last[op]
	return name, ok
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

// fakeOutput keeps every stream apart: diagnostics, results and debug lines
// are three different things and a test must be able to tell them apart.
type fakeOutput struct {
	infos   []string
	results []string
	debugs  []string
	table   string
	rows    int
}

func (o *fakeOutput) Info(message string, data ...any) {
	o.infos = append(o.infos, fmt.Sprint(append([]any{message}, data...)...))
}

func (o *fakeOutput) Table(commands catalog.Catalog) string {
	o.rows = commands.Len()
	return o.table
}

func (o *fakeOutput) Result(text string) {
	o.results = append(o.results, text)
}

func (o *fakeOutput) Debug(message string, data ...any) {
	o.debugs = append(o.debugs, fmt.Sprint(append([]any{message}, data...)...))
}

func (o *fakeOutput) infoText() string {
	return strings.Join(o.infos, "\n")
}

type fakePrompt struct {
	rows       []string
	chooseErr  error
	chooseCall int
	answer     bool
	confirmErr error
	asked      []string
}

func (p *fakePrompt) ChooseCommand() (string, error) {
	if p.chooseErr != nil {
		return "", p.chooseErr
	}

	p.chooseCall++

	if len(p.rows) == 0 {
		return "", nil
	}

	row := p.rows[0]
	p.rows = p.rows[1:]

	return row, nil
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)
	return p.answer, p.confirmErr
}

type fakeRunner struct {
	calls       []string
	nestedCalls []string
	shells      []string
	dirs        []string
	envs        []map[string]string
	timeouts    []string
	out         map[string]string
	err         error
	errOnLine   string
}

func (r *fakeRunner) Run(line string, target executecommand.Target) (string, error) {
	r.record(line, target)
	r.calls = append(r.calls, line)

	if r.errOnLine != "" && line == r.errOnLine {
		return r.out[line], r.err
	}

	return r.out[line], nil
}

func (r *fakeRunner) RunNested(line string, target executecommand.Target) (string, error) {
	r.record(line, target)
	r.nestedCalls = append(r.nestedCalls, line)

	if r.errOnLine != "" && line == r.errOnLine {
		return r.out[line], r.err
	}

	return r.out[line], nil
}

// record keeps each target field separately, so a test can assert on the one
// it cares about without rebuilding the struct.
func (r *fakeRunner) record(_ string, target executecommand.Target) {
	r.shells = append(r.shells, target.Shell)
	r.dirs = append(r.dirs, target.Dir)
	r.envs = append(r.envs, target.Env)
	r.timeouts = append(r.timeouts, target.Timeout)
}

type fakeNotes struct{}

func (fakeNotes) ExecuteCommandNote(command catalog.Command, params []string) string {
	return "note:command:" + command.Name + ":" + fmt.Sprint(params)
}

func (fakeNotes) ExecuteLastNote(op operation.Operation, command catalog.Command, _ []string) string {
	return "note:last:" + string(op) + ":" + command.Name
}

type harness struct {
	state  *fakeState
	output *fakeOutput
	prompt *fakePrompt
	runner *fakeRunner
}

func newHarness(commands ...catalog.Command) *harness {
	return &harness{
		state: &fakeState{
			commands: catalog.New(nil, commands),
			last:     map[operation.Operation]string{},
		},
		output: &fakeOutput{table: "TABLE"},
		prompt: &fakePrompt{answer: true},
		runner: &fakeRunner{out: map[string]string{}},
	}
}

func (h *harness) useCase() *executecommand.UseCase {
	return executecommand.New(h.state, h.output, h.prompt, h.runner, fakeNotes{})
}

func (h *harness) run(request executecommand.Request) error {
	return h.useCase().Run(request)
}

func twoCommands() []catalog.Command {
	return []catalog.Command{
		{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"},
		{Name: "build", Exec: "go build ./..."},
	}
}

func TestRunByName(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.runner.out["echo hi"] = "hi\n"

	if err := h.run(executecommand.Request{Key: "greet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.runner.calls, []string{"echo hi"}; !slices.Equal(got, want) {
		t.Errorf("ran %v, want %v", got, want)
	}

	// The output is a result, trimmed of the trailing newline, not a diagnostic.
	if got, want := h.output.results, []string{"hi"}; !slices.Equal(got, want) {
		t.Errorf("results = %v, want %v", got, want)
	}

	if len(h.state.records) != 1 {
		t.Fatalf("recorded %d entries, want 1", len(h.state.records))
	}

	got := h.state.records[0]
	if got.op != operation.Execute || got.commandName != "greet" || got.note != "note:command:greet:[]" {
		t.Errorf("recorded %+v, want execute/greet/note", got)
	}

	// A keyed run neither shows the table nor asks anything.
	if h.output.rows != 0 || len(h.prompt.asked) != 0 {
		t.Errorf("a keyed run prompted: rows=%d asked=%v", h.output.rows, h.prompt.asked)
	}
}

func TestRunByAlias(t *testing.T) {
	h := newHarness(twoCommands()...)

	if err := h.run(executecommand.Request{Key: "hi"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.state.records[0].commandName != "greet" {
		t.Errorf("recorded command %q, want greet", h.state.records[0].commandName)
	}
}

func TestRunUnknownCommand(t *testing.T) {
	h := newHarness(twoCommands()...)

	err := h.run(executecommand.Request{Key: "nope"})
	if err == nil || !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("error = %v, want command not found", err)
	}

	if len(h.runner.calls) != 0 || len(h.state.records) != 0 {
		t.Error("an unknown command still ran or recorded something")
	}
}

// With no key create asks which command to run, shows it, and confirms before
// running it.
func TestRunAsksAndConfirms(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"build"}

	if err := h.run(executecommand.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.output.rows != 2 {
		t.Errorf("table shows %d rows, want 2", h.output.rows)
	}

	text := h.output.infoText()

	for _, want := range []string{"Name:    build", "Exec:    go build ./..."} {
		if !strings.Contains(text, want) {
			t.Errorf("diagnostics = %q, want %q", text, want)
		}
	}

	if got, want := h.prompt.asked, []string{"Continue to execute"}; !slices.Equal(got, want) {
		t.Errorf("questions = %v, want %v", got, want)
	}

	if h.state.records[0].commandName != "build" {
		t.Errorf("recorded %q, want build", h.state.records[0].commandName)
	}
}

func TestRunUnknownSelectionAsksAgain(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"nope", "greet"}

	if err := h.run(executecommand.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(h.output.infoText(), "Command not found, try to select one.") {
		t.Errorf("diagnostics = %q, want the retry notice", h.output.infoText())
	}

	if h.prompt.chooseCall != 2 {
		t.Errorf("ChooseCommand called %d time(s), want 2", h.prompt.chooseCall)
	}
}

func TestRunUnknownSelectionTwiceFails(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"nope", "still-nope"}

	err := h.run(executecommand.Request{})
	if err == nil || !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("error = %v, want command not found", err)
	}

	if len(h.runner.calls) != 0 {
		t.Error("a failed selection still ran a command")
	}
}

func TestRunDeclined(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.prompt.rows = []string{"greet"}
	h.prompt.answer = false

	if err := h.run(executecommand.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.runner.calls) != 0 || len(h.state.records) != 0 {
		t.Error("a declined run still ran or recorded something")
	}

	if !strings.Contains(h.output.infoText(), "Cancelling execute...") {
		t.Errorf("diagnostics = %q, want the cancel line", h.output.infoText())
	}
}

func TestRunPromptErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("choosing", func(t *testing.T) {
		h := newHarness(twoCommands()...)
		h.prompt.chooseErr = boom

		if err := h.run(executecommand.Request{}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("confirming", func(t *testing.T) {
		h := newHarness(twoCommands()...)
		h.prompt.rows = []string{"greet"}
		h.prompt.confirmErr = boom

		if err := h.run(executecommand.Request{}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

func TestRunEmptySelectionIsInvalid(t *testing.T) {
	h := newHarness(twoCommands()...)

	err := h.run(executecommand.Request{})
	if err == nil || !strings.Contains(err.Error(), "invalid row") {
		t.Fatalf("error = %v, want invalid row", err)
	}
}

func TestRunAppliesParams(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Exec: "echo {1}"})

	if err := h.run(executecommand.Request{Key: "greet", Params: []string{"hi"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got := len(h.runner.calls); got != 1 {
		t.Fatalf("ran %d lines, want 1", got)
	}

	if !strings.Contains(h.runner.calls[0], "hi") {
		t.Errorf("ran %q, want the parameter substituted", h.runner.calls[0])
	}

	if !strings.Contains(h.state.records[0].note, "[hi]") {
		t.Errorf("note = %q, want the params recorded", h.state.records[0].note)
	}
}

// A Steps command runs its steps, not Exec. An unusable brace in Exec must not
// turn into an error for a line that never runs.
func TestRunWithStepsDoesNotValidateExec(t *testing.T) {
	h := newHarness(catalog.Command{Name: "multi", Exec: "mkdir {a,b}", Steps: []string{"echo {1}"}})

	if err := h.run(executecommand.Request{Key: "multi", Params: []string{"a"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.runner.calls) != 1 {
		t.Fatalf("ran %v, want the one step", h.runner.calls)
	}
}

func TestRunAppliesParamsToSteps(t *testing.T) {
	h := newHarness(catalog.Command{Name: "multi", Exec: "unused", Steps: []string{"echo {1}", "echo {2}"}})

	if err := h.run(executecommand.Request{Key: "multi", Params: []string{"a", "b"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.runner.calls) != 2 {
		t.Fatalf("ran %v, want two steps", h.runner.calls)
	}
}

func TestRunRejectsParametersWithoutAPlaceholder(t *testing.T) {
	h := newHarness(catalog.Command{Name: "brace", Exec: "mkdir {a,b}"})

	err := h.run(executecommand.Request{Key: "brace", Params: []string{"x"}})
	if err == nil || !strings.Contains(err.Error(), "no parameter placeholder") {
		t.Fatalf("error = %v, want it to name the missing placeholder", err)
	}

	if len(h.runner.calls) != 0 {
		t.Error("a refused run still executed something")
	}
}

// A steps set with no placeholder is an error, not a silent drop.
func TestRunRejectsParametersWhenNoStepHasAPlaceholder(t *testing.T) {
	h := newHarness(catalog.Command{Name: "multi", Exec: "unused", Steps: []string{"echo a"}})

	err := h.run(executecommand.Request{Key: "multi", Params: []string{"x"}})
	if err == nil || !strings.Contains(err.Error(), "no parameter placeholder") {
		t.Fatalf("error = %v, want it to name the missing placeholder", err)
	}
}

// A step with a literal brace is left alone while another step substitutes.
func TestRunKeepsALiteralBraceInAStep(t *testing.T) {
	h := newHarness(catalog.Command{Name: "multi", Exec: "unused", Steps: []string{"echo {x}", "echo {1}"}})

	if err := h.run(executecommand.Request{Key: "multi", Params: []string{"param"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	want := []string{"echo {x}", "echo " + shell.Quote("param")}
	if !slices.Equal(h.runner.calls, want) {
		t.Errorf("ran %v, want %v", h.runner.calls, want)
	}
}

// The catalog is a value: running a Steps command cannot write the substituted
// line back over the stored steps.
func TestRunDoesNotRewriteStoredSteps(t *testing.T) {
	h := newHarness(catalog.Command{Name: "multi", Exec: "unused", Steps: []string{"echo {1}"}})

	if err := h.run(executecommand.Request{Key: "multi", Params: []string{"a"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	stored, found := h.state.commands.Get("multi")
	if !found || !slices.Equal(stored.Steps, []string{"echo {1}"}) {
		t.Errorf("stored steps = %v, want them unchanged", stored.Steps)
	}
}

func TestRunNestedStepUsesTheGuardedRunner(t *testing.T) {
	h := newHarness(
		catalog.Command{Name: "outer", Exec: "unused", Steps: []string{"crower list", "echo done"}},
	)

	if err := h.run(executecommand.Request{Key: "outer"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.runner.nestedCalls, []string{"crower list"}; !slices.Equal(got, want) {
		t.Errorf("nested calls = %v, want %v", got, want)
	}

	if got, want := h.runner.calls, []string{"echo done"}; !slices.Equal(got, want) {
		t.Errorf("plain calls = %v, want %v", got, want)
	}
}

// A command's shell and folder reach every step: a step inherits them, which is
// why the schema needs no per-step syntax.
func TestRunPassesShellAndCwdToEveryStep(t *testing.T) {
	h := newHarness(catalog.Command{
		Name:  "deploy",
		Shell: "pwsh",
		Cwd:   "~/projects/api",
		Steps: []string{"echo a", "echo b"},
	})

	if err := h.run(executecommand.Request{Key: "deploy"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.runner.shells, []string{"pwsh", "pwsh"}; !slices.Equal(got, want) {
		t.Errorf("shells = %v, want %v", got, want)
	}

	if got, want := h.runner.dirs, []string{"~/projects/api", "~/projects/api"}; !slices.Equal(got, want) {
		t.Errorf("dirs = %v, want %v", got, want)
	}
}

// A nested crower step runs in the same shell and folder.
func TestRunPassesShellAndCwdToANestedStep(t *testing.T) {
	h := newHarness(catalog.Command{Name: "outer", Shell: "sh", Cwd: "/tmp", Steps: []string{"crower list"}})

	if err := h.run(executecommand.Request{Key: "outer"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.runner.shells, []string{"sh"}; !slices.Equal(got, want) {
		t.Errorf("shells = %v, want %v", got, want)
	}

	if got, want := h.runner.dirs, []string{"/tmp"}; !slices.Equal(got, want) {
		t.Errorf("dirs = %v, want %v", got, want)
	}
}

// The environment and timeout reach every step too, and they are the command's
// own values rather than a copy the use case invented.
func TestRunPassesEnvAndTimeoutToEveryStep(t *testing.T) {
	h := newHarness(catalog.Command{
		Name:    "deploy",
		Env:     map[string]string{"STAGE": "prod"},
		Timeout: "30s",
		Steps:   []string{"echo a", "echo b"},
	})

	if err := h.run(executecommand.Request{Key: "deploy"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.runner.envs) != 2 {
		t.Fatalf("recorded %d envs, want 2", len(h.runner.envs))
	}

	for i, env := range h.runner.envs {
		if env["STAGE"] != "prod" {
			t.Errorf("step %d env = %v, want STAGE=prod", i, env)
		}
	}

	if got, want := h.runner.timeouts, []string{"30s", "30s"}; !slices.Equal(got, want) {
		t.Errorf("timeouts = %v, want %v", got, want)
	}
}

// No declared shell or folder stays empty: the runner turns empty into the
// platform default, and the use case must not choose one.
func TestRunWithoutATargetKeepsThePlatformDefault(t *testing.T) {
	h := newHarness(catalog.Command{Name: "greet", Exec: "echo hi"})

	if err := h.run(executecommand.Request{Key: "greet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.runner.shells, []string{""}; !slices.Equal(got, want) {
		t.Errorf("shells = %v, want %v", got, want)
	}

	if got, want := h.runner.dirs, []string{""}; !slices.Equal(got, want) {
		t.Errorf("dirs = %v, want %v", got, want)
	}
}

func TestRunEmptyStepsAreSkipped(t *testing.T) {
	h := newHarness(catalog.Command{Name: "spaced", Exec: "unused", Steps: []string{"", "   ", "echo a"}})

	if err := h.run(executecommand.Request{Key: "spaced"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.runner.calls, []string{"echo a"}; !slices.Equal(got, want) {
		t.Errorf("ran %v, want %v", got, want)
	}
}

func TestRunReportsNothingForEmptyOutput(t *testing.T) {
	h := newHarness(catalog.Command{Name: "quiet", Exec: "true"})
	h.runner.out["true"] = "\n\n"

	if err := h.run(executecommand.Request{Key: "quiet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.output.results) != 0 {
		t.Errorf("results = %v, want nothing", h.output.results)
	}
}

func TestRunRunnerErrorFailsWithoutRecording(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness(catalog.Command{Name: "bad", Exec: "false"})
	h.runner.err = boom
	h.runner.errOnLine = "false"
	h.runner.out["false"] = "partial"

	err := h.run(executecommand.Request{Key: "bad"})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if len(h.output.results) != 0 {
		t.Errorf("results = %v, want the partial output discarded like the CLI always did", h.output.results)
	}

	if len(h.state.records) != 0 {
		t.Error("a failed run recorded an entry")
	}
}

func TestRunWritesDebugLines(t *testing.T) {
	h := newHarness(catalog.Command{Name: "multi", Exec: "unused", Steps: []string{"echo {1}"}})

	if err := h.run(executecommand.Request{Key: "multi", Params: []string{"a"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(strings.Join(h.output.debugs, "\n"), "[1/1]") {
		t.Errorf("debug = %v, want the step marker", h.output.debugs)
	}

	plain := newHarness(catalog.Command{Name: "greet", Exec: "echo hi"})
	if err := plain.run(executecommand.Request{Key: "greet"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !strings.Contains(strings.Join(plain.output.debugs, "\n"), "exec: echo hi") {
		t.Errorf("debug = %v, want the exec line", plain.output.debugs)
	}
}

func TestRunLast(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.state.last[operation.Create] = "build"

	if err := h.run(executecommand.Request{Last: operation.Create, Params: []string{"p"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.prompt.asked) != 0 || h.output.rows != 0 {
		t.Error("the --last form asked something")
	}

	if len(h.runner.calls) != 1 {
		t.Fatalf("ran %v, want the last command", h.runner.calls)
	}

	if got := h.state.records[0]; got.note != "note:last:create:build" {
		t.Errorf("recorded %+v, want the --last note", got)
	}
}

func TestRunLastWithoutHistoryFails(t *testing.T) {
	h := newHarness(twoCommands()...)

	err := h.run(executecommand.Request{Last: operation.Update})
	if err == nil || !strings.Contains(err.Error(), `command not found: "update"`) {
		t.Fatalf("error = %v, want command not found", err)
	}
}

func TestRunRecordErrorFailsTheRun(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness(catalog.Command{Name: "greet", Exec: "echo hi"})
	h.state.recordErr = boom

	if err := h.run(executecommand.Request{Key: "greet"}); !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}
}

// The UI entry point answers with the output and the command's name, and never
// prompts or records: those callers own their own interaction.
func TestExecuteForTheUIs(t *testing.T) {
	h := newHarness(twoCommands()...)
	h.runner.out["echo hi"] = "hi\n"

	out, name, err := h.useCase().Execute("hi", nil)
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	if out != "hi\n" {
		t.Errorf("output = %q, want the raw output", out)
	}

	if name != "greet" {
		t.Errorf("command = %q, want greet", name)
	}

	if len(h.state.records) != 0 || len(h.output.results) != 0 || len(h.prompt.asked) != 0 {
		t.Error("the UI entry point recorded, reported or prompted")
	}
}

func TestExecuteForTheUIsUnknownCommand(t *testing.T) {
	h := newHarness(twoCommands()...)

	_, name, err := h.useCase().Execute("nope", nil)
	if err == nil || !strings.Contains(err.Error(), "command not found") {
		t.Fatalf("error = %v, want command not found", err)
	}

	if name != "" {
		t.Errorf("command = %q, want empty", name)
	}
}
