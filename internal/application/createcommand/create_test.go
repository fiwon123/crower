package createcommand_test

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/createcommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// recorder keeps the order of the effects a test cares about.
type recorder struct {
	events []string
}

func (r *recorder) add(event string) {
	if r == nil {
		return
	}

	r.events = append(r.events, event)
}

// fakeState records what the use case asked it to store.
type fakeState struct {
	commands  catalog.Catalog
	last      map[operation.Operation]createcommand.LastRun
	saved     []catalog.Catalog
	records   []record
	saveErr   error
	recordErr error
	events    *recorder
}

type record struct {
	op          operation.Operation
	commandName string
	note        string
}

func (s *fakeState) Commands() catalog.Catalog {
	return s.commands
}

func (s *fakeState) Save(commands catalog.Catalog) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	s.events.add("save")
	s.saved = append(s.saved, commands)
	return nil
}

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.events.add("record")
	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

func (s *fakeState) LastOperation(op operation.Operation) (createcommand.LastRun, bool) {
	run, ok := s.last[op]
	return run, ok
}

// fakeOutput keeps every diagnostic so a test can assert what the user saw.
type fakeOutput struct {
	lines  []string
	events *recorder
}

func (o *fakeOutput) Info(message string, data ...any) {
	o.events.add("info")
	o.lines = append(o.lines, fmt.Sprint(append([]any{message}, data...)...))
}

func (o *fakeOutput) text() string {
	return strings.Join(o.lines, "\n")
}

type fakePrompt struct {
	name      string
	exec      string
	aliases   []string
	answer    bool
	nameErr   error
	execErr   error
	aliasErr  error
	askErr    error
	asked     []string
	nameCalls int
	execCalls int
	aliasCall int
}

func (p *fakePrompt) CreateName() (string, error) {
	p.nameCalls++
	return p.name, p.nameErr
}

func (p *fakePrompt) CreateExec() (string, error) {
	p.execCalls++
	return p.exec, p.execErr
}

func (p *fakePrompt) NewAliases() ([]string, error) {
	p.aliasCall++
	return p.aliases, p.aliasErr
}

func (p *fakePrompt) Confirm(question string) (bool, error) {
	p.asked = append(p.asked, question)
	return p.answer, p.askErr
}

type fakeCreator struct {
	calls     []string
	errs      map[string]error
	events    *recorder
	lastError error
}

func (c *fakeCreator) File(path string) error {
	c.events.add("create file " + path)
	c.calls = append(c.calls, "file "+path)

	if err := c.errs[path]; err != nil {
		return err
	}

	return c.lastError
}

func (c *fakeCreator) Folder(path string) error {
	c.events.add("create folder " + path)
	c.calls = append(c.calls, "folder "+path)

	if err := c.errs[path]; err != nil {
		return err
	}

	return c.lastError
}

type fakeScripts struct {
	path   string
	folder string
	err    error
	names  []string
	events *recorder
}

func (s *fakeScripts) Create(name string) (string, string, error) {
	s.events.add("script")
	s.names = append(s.names, name)
	return s.path, s.folder, s.err
}

type fakeProcesses struct {
	execLine string
	err      error
	asked    []string
}

func (p *fakeProcesses) ExecLine(process string) (string, error) {
	p.asked = append(p.asked, process)
	return p.execLine, p.err
}

type fakeEnvironment struct {
	out     string
	err     error
	calls   []string
	events  *recorder
	setErr  error
	pathErr error
}

func (e *fakeEnvironment) SetVariable(name, value string) (string, error) {
	e.events.add("set variable")
	e.calls = append(e.calls, "variable "+name+"="+value)

	if e.setErr != nil {
		return "", e.setErr
	}

	return e.out, e.err
}

func (e *fakeEnvironment) AddPath(path string) (string, error) {
	e.events.add("add path")
	e.calls = append(e.calls, "path "+path)

	if e.pathErr != nil {
		return "", e.pathErr
	}

	return e.out, e.err
}

type fakeOpener struct {
	opened []string
	err    error
	events *recorder
}

func (o *fakeOpener) Open(path string) error {
	o.events.add("open")
	o.opened = append(o.opened, path)
	return o.err
}

// fakeNotes renders notes that name the mode, so a test can tell which one the
// use case asked for without depending on the note format.
type fakeNotes struct{}

func (fakeNotes) CreateCommandNote(command catalog.Command) string {
	return "note:command:" + command.Name
}

func (fakeNotes) CreateScriptNote(command catalog.Command) string {
	return "note:script:" + command.Name
}

func (fakeNotes) CreateLastNote(command catalog.Command) string {
	return "note:last:" + command.Name
}

func (fakeNotes) CreateProcessNote(command catalog.Command, args []string) string {
	return "note:process:" + command.Name + ":" + fmt.Sprint(args)
}

func (fakeNotes) CreateSystemVariableNote(args []string) string {
	return "note:system:" + fmt.Sprint(args)
}

func (fakeNotes) CreateSystemPathNote(args []string) string {
	return "note:syspath:" + fmt.Sprint(args)
}

func (fakeNotes) CreateFileNote(args []string) string {
	return "note:file:" + fmt.Sprint(args)
}

func (fakeNotes) CreateFolderNote(args []string) string {
	return "note:folder:" + fmt.Sprint(args)
}

type harness struct {
	state   *fakeState
	output  *fakeOutput
	prompt  *fakePrompt
	creator *fakeCreator
	scripts *fakeScripts
	process *fakeProcesses
	environ *fakeEnvironment
	opener  *fakeOpener
	events  *recorder
}

func newHarness(commands ...catalog.Command) *harness {
	events := &recorder{}

	h := &harness{
		events: events,
		state: &fakeState{
			commands: catalog.New(nil, commands),
			last:     map[operation.Operation]createcommand.LastRun{},
			events:   events,
		},
		output:  &fakeOutput{events: events},
		prompt:  &fakePrompt{name: "asked-name", exec: "asked-exec", aliases: []string{"asked-alias"}},
		creator: &fakeCreator{errs: map[string]error{}, events: events},
		scripts: &fakeScripts{path: "/cfg/scripts/deploy.sh", folder: "/cfg/scripts", events: events},
		process: &fakeProcesses{execLine: "/usr/bin/sleeper"},
		environ: &fakeEnvironment{out: "environment changed", events: events},
		opener:  &fakeOpener{events: events},
	}

	return h
}

func (h *harness) run(request createcommand.Request) error {
	return createcommand.New(
		h.state, h.output, h.prompt, h.creator, h.scripts, h.process, h.environ, h.opener, fakeNotes{},
	).Run(request)
}

func twoCommands() []catalog.Command {
	return []catalog.Command{
		{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"},
		{Name: "build", Exec: "go build ./..."},
	}
}

// The metadata flags reach the command the create mode adds.
func TestCreateCommandCarriesTheMetadata(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{
		Mode:        createcommand.ModeCommand,
		Args:        []string{"deploy", "make deploy"},
		Description: "ship it",
		Params:      []string{"env", "region"},
		Cwd:         "/srv/app",
		Shell:       "sh",
		Timeout:     "120s",
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	created, found := h.state.saved[0].Get("deploy")
	if !found {
		t.Fatal("the saved catalog does not hold the new command")
	}

	if created.Description != "ship it" {
		t.Errorf("Description = %q, want ship it", created.Description)
	}

	if !slices.Equal(created.Params, []string{"env", "region"}) {
		t.Errorf("Params = %v, want [env region]", created.Params)
	}

	if created.Cwd != "/srv/app" || created.Shell != "sh" || created.Timeout != "120s" {
		t.Errorf("Cwd/Shell/Timeout = %q/%q/%q, want /srv/app/sh/120s", created.Cwd, created.Shell, created.Timeout)
	}
}

func TestCreateCommandWithTwoArguments(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{Mode: createcommand.ModeCommand, Args: []string{"greet", "echo hi"}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.nameCalls != 0 || h.prompt.execCalls != 0 || len(h.prompt.asked) != 0 {
		t.Errorf("a complete command still prompted: name=%d exec=%d asked=%v", h.prompt.nameCalls, h.prompt.execCalls, h.prompt.asked)
	}

	if len(h.state.saved) != 1 {
		t.Fatalf("saved %d catalogs, want 1", len(h.state.saved))
	}

	if _, found := h.state.saved[0].Get("greet"); !found {
		t.Error("the saved catalog does not hold the new command")
	}

	if len(h.state.records) != 1 {
		t.Fatalf("recorded %d entries, want 1", len(h.state.records))
	}

	got := h.state.records[0]
	if got.op != operation.Create || got.commandName != "greet" || got.note != "note:command:greet" {
		t.Errorf("recorded %+v, want create/greet/note", got)
	}

	if !strings.Contains(h.output.text(), "added new command: greet") {
		t.Errorf("output = %q, want the success line", h.output.text())
	}
}

// Flags are used even when the two positional values complete the command.
func TestCreateCommandKeepsTheAliasFlag(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{
		Mode:    createcommand.ModeCommand,
		Aliases: []string{"hi"},
		Args:    []string{"greet", "echo hi"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if owner, found := h.state.saved[0].AliasOwner("hi"); !found || owner != "greet" {
		t.Errorf("alias hi = %q, %v; want greet, true", owner, found)
	}
}

// With no positional values create asks for each field, in the order it always
// has, and the alias question comes before the alias list.
func TestCreateCommandAsksForTheMissingFields(t *testing.T) {
	h := newHarness()
	h.prompt.answer = true

	if err := h.run(createcommand.Request{Mode: createcommand.ModeCommand}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.nameCalls != 1 || h.prompt.execCalls != 1 || h.prompt.aliasCall != 1 {
		t.Errorf("asked name=%d exec=%d aliases=%d, want one of each", h.prompt.nameCalls, h.prompt.execCalls, h.prompt.aliasCall)
	}

	if got, want := h.prompt.asked, []string{"Do you want to add alias"}; !slices.Equal(got, want) {
		t.Errorf("questions = %v, want %v", got, want)
	}

	created, found := h.state.saved[0].Get("asked-name")
	if !found {
		t.Fatal("the prompted command was not saved")
	}

	if created.Exec != "asked-exec" {
		t.Errorf("exec = %q, want the prompted one", created.Exec)
	}

	if !slices.Equal(created.Aliases, []string{"asked-alias"}) {
		t.Errorf("aliases = %v, want the prompted ones", created.Aliases)
	}
}

func TestCreateCommandDecliningAliasesSkipsThem(t *testing.T) {
	h := newHarness()
	h.prompt.answer = false

	if err := h.run(createcommand.Request{Mode: createcommand.ModeCommand}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.aliasCall != 0 {
		t.Errorf("the alias list was asked for after a decline")
	}
}

// Three positional values are not a command, so create asks for everything --
// the same fallback as before.
func TestCreateCommandWithThreeArgumentsAsks(t *testing.T) {
	h := newHarness()

	if err := h.run(createcommand.Request{Mode: createcommand.ModeCommand, Args: []string{"a", "b", "c"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.nameCalls != 1 || h.prompt.execCalls != 1 {
		t.Errorf("asked name=%d exec=%d, want one of each", h.prompt.nameCalls, h.prompt.execCalls)
	}
}

func TestCreateCommandPromptErrors(t *testing.T) {
	boom := errors.New("boom")

	tests := []struct {
		name  string
		setup func(*fakePrompt)
	}{
		{"name", func(p *fakePrompt) { p.nameErr = boom }},
		{"aliases", func(p *fakePrompt) { p.answer = true; p.aliasErr = boom }},
		{"exec", func(p *fakePrompt) { p.execErr = boom }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness()
			tt.setup(h.prompt)

			if err := h.run(createcommand.Request{Mode: createcommand.ModeCommand}); !errors.Is(err, boom) {
				t.Fatalf("error = %v, want boom", err)
			}

			if len(h.state.saved) != 0 || len(h.state.records) != 0 {
				t.Error("a prompt error still wrote something")
			}
		})
	}
}

// The collision rules live in the catalog; the use case wraps their error the
// way create always has, and writes nothing.
func TestCreateCommandCollision(t *testing.T) {
	h := newHarness(twoCommands()...)

	err := h.run(createcommand.Request{Mode: createcommand.ModeCommand, Args: []string{"greet", "echo again"}})
	if err == nil {
		t.Fatal("creating an existing name returned no error")
	}

	if want := `add command "greet": found name, command already added`; err.Error() != want {
		t.Errorf("error = %q, want %q", err, want)
	}

	if len(h.state.saved) != 0 || len(h.state.records) != 0 {
		t.Error("a refused create still wrote something")
	}
}

func TestCreateProcess(t *testing.T) {
	h := newHarness()
	h.prompt.answer = false

	err := h.run(createcommand.Request{
		Mode: createcommand.ModeProcess,
		Name: "sleeper",
		Args: []string{"sleep"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.process.asked, []string{"sleep"}; !slices.Equal(got, want) {
		t.Errorf("resolved processes = %v, want %v", got, want)
	}

	if !strings.Contains(h.output.text(), "added new command by process: sleeper") {
		t.Errorf("output = %q, want the process success line", h.output.text())
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:process:sleeper:[sleep]" {
		t.Errorf("recorded = %+v, want the process note with the original args", h.state.records)
	}

	created, found := h.state.saved[0].Get("sleeper")
	if !found || created.Exec != "/usr/bin/sleeper" {
		t.Errorf("saved command = %+v, %v; want the resolved exec", created, found)
	}
}

// With no --process value the first positional becomes the name, as before.
func TestCreateProcessTakesTheNameFromTheArgumentsWhenTheFlagIsEmpty(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{Mode: createcommand.ModeProcess, Args: []string{"sleeper", "sleep"}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.process.asked, []string{"sleep"}; !slices.Equal(got, want) {
		t.Errorf("resolved processes = %v, want %v", got, want)
	}

	if _, found := h.state.saved[0].Get("sleeper"); !found {
		t.Error("the command was not saved under the shifted name")
	}
}

func TestCreateProcessWithoutAProcessFails(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"no arguments", nil},
		{"empty", []string{""}},
		{"whitespace", []string{"   "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness()
			// The name comes from the flag, so the positional value stays the
			// process and is checked as such.
			err := h.run(createcommand.Request{Mode: createcommand.ModeProcess, Name: "x", Args: tt.args})

			if err == nil || !strings.Contains(err.Error(), "process name or pid") {
				t.Fatalf("error = %v, want it to name the missing process", err)
			}

			if len(h.process.asked) != 0 {
				t.Error("a process lookup ran with no process")
			}
		})
	}
}

func TestCreateProcessLookupErrorIsWrapped(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.process.err = boom

	err := h.run(createcommand.Request{Mode: createcommand.ModeProcess, Name: "x", Args: []string{"ghost"}})

	if !errors.Is(err, boom) || !strings.HasPrefix(err.Error(), "add command by process: ") {
		t.Fatalf("error = %v, want boom wrapped as a process failure", err)
	}

	if len(h.state.saved) != 0 {
		t.Error("a failed lookup still saved a catalog")
	}
}

// A script create writes the file, registers the command, opens the folder and
// only then records the entry -- the order it has always had, so the history
// holds a command the user was actually sent to edit.
func TestCreateScriptOrder(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{
		Mode:    createcommand.ModeScript,
		Aliases: []string{"d"},
		Args:    []string{"deploy"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	want := []string{"script", "save", "info", "open", "record"}
	if !slices.Equal(h.events.events, want) {
		t.Errorf("events = %v, want %v", h.events.events, want)
	}

	if got, want := h.opener.opened, []string{"/cfg/scripts"}; !slices.Equal(got, want) {
		t.Errorf("opened = %v, want %v", got, want)
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:script:deploy" {
		t.Errorf("recorded = %+v, want the script note", h.state.records)
	}

	created, found := h.state.saved[0].Get("deploy")
	if !found || created.Exec != "/cfg/scripts/deploy.sh" {
		t.Errorf("saved command = %+v, %v; want the script path as its exec", created, found)
	}
}

// If the folder cannot be opened the command is already written, but the entry
// is not recorded -- the run reports the failure.
func TestCreateScriptOpenFailure(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.opener.err = boom

	err := h.run(createcommand.Request{Mode: createcommand.ModeScript, Args: []string{"deploy"}})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if len(h.state.saved) != 1 {
		t.Errorf("saved %d catalogs, want the command file written", len(h.state.saved))
	}

	if len(h.state.records) != 0 {
		t.Errorf("recorded %d entries, want none", len(h.state.records))
	}
}

func TestCreateScriptWithoutANameFails(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{Mode: createcommand.ModeScript})
	if err == nil || !strings.Contains(err.Error(), "missing arguments: name") {
		t.Fatalf("error = %v, want the missing-name error", err)
	}

	if len(h.scripts.names) != 0 {
		t.Error("a script was created without a name")
	}
}

func TestCreateLast(t *testing.T) {
	h := newHarness()
	h.state.last[operation.Execute] = createcommand.LastRun{CommandName: "old", Exec: "echo old"}

	err := h.run(createcommand.Request{
		Mode:    createcommand.ModeLast,
		Last:    operation.Execute,
		Aliases: []string{"g"},
		Args:    []string{"greet"},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.prompt.nameCalls != 0 || h.prompt.execCalls != 0 || len(h.prompt.asked) != 0 {
		t.Error("create --execute prompted")
	}

	created, found := h.state.saved[0].Get("greet")
	if !found || created.Exec != "echo old" {
		t.Errorf("saved command = %+v, %v; want the last exec", created, found)
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:last:greet" {
		t.Errorf("recorded = %+v, want the last note", h.state.records)
	}
}

func TestCreateLastRefusesWhatItCannotRebuild(t *testing.T) {
	tests := []struct {
		name string
		run  createcommand.LastRun
		want string
	}{
		{
			"ran steps",
			createcommand.LastRun{CommandName: "multi", Exec: "unused", Steps: []string{"echo a", "echo b"}, HasSteps: true},
			`the last "execute" operation ("multi") ran 2 step(s); create --execute cannot rebuild them`,
		},
		{
			"unreadable steps",
			createcommand.LastRun{CommandName: "multi", Exec: "unused", HasSteps: true},
			`the last "execute" operation ("multi") has steps that cannot be read; create --execute cannot rebuild them`,
		},
		{
			"no exec",
			createcommand.LastRun{CommandName: "multi"},
			`the last "execute" operation carries no exec to reuse`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness()
			h.state.last[operation.Execute] = tt.run

			err := h.run(createcommand.Request{Mode: createcommand.ModeLast, Last: operation.Execute, Args: []string{"copy"}})
			if err == nil || err.Error() != tt.want {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}

			if len(h.state.saved) != 0 || len(h.state.records) != 0 {
				t.Error("a refused create --execute still wrote something")
			}
		})
	}
}

func TestCreateLastArgumentChecks(t *testing.T) {
	t.Run("no name", func(t *testing.T) {
		h := newHarness()

		err := h.run(createcommand.Request{Mode: createcommand.ModeLast, Last: operation.Execute})
		if err == nil || !strings.Contains(err.Error(), "missing arguments: name") {
			t.Fatalf("error = %v, want the missing-name error", err)
		}
	})

	t.Run("too many names", func(t *testing.T) {
		h := newHarness()

		err := h.run(createcommand.Request{Mode: createcommand.ModeLast, Last: operation.Execute, Args: []string{"a", "b"}})
		if err == nil || !strings.Contains(err.Error(), "too many arguments") {
			t.Fatalf("error = %v, want the too-many error", err)
		}
	})

	t.Run("nothing in history", func(t *testing.T) {
		h := newHarness()

		err := h.run(createcommand.Request{Mode: createcommand.ModeLast, Last: operation.Execute, Args: []string{"a"}})
		if err == nil || !strings.Contains(err.Error(), `command not found: "execute"`) {
			t.Fatalf("error = %v, want command not found", err)
		}
	})
}

func TestCreateFileAndFolder(t *testing.T) {
	t.Run("file", func(t *testing.T) {
		h := newHarness()

		if err := h.run(createcommand.Request{Mode: createcommand.ModeFile, Args: []string{"a.txt", "b.txt"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.creator.calls, []string{"file a.txt", "file b.txt"}; !slices.Equal(got, want) {
			t.Errorf("creator calls = %v, want %v", got, want)
		}

		if len(h.state.records) != 1 || h.state.records[0].note != "note:file:[a.txt b.txt]" {
			t.Errorf("recorded = %+v, want one file note", h.state.records)
		}

		if h.state.records[0].commandName != "" {
			t.Errorf("recorded command name = %q, want empty", h.state.records[0].commandName)
		}
	})

	t.Run("folder", func(t *testing.T) {
		h := newHarness()

		if err := h.run(createcommand.Request{Mode: createcommand.ModeFolder, Args: []string{"one", "two"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.creator.calls, []string{"folder one", "folder two"}; !slices.Equal(got, want) {
			t.Errorf("creator calls = %v, want %v", got, want)
		}

		if len(h.state.records) != 1 || h.state.records[0].note != "note:folder:[one two]" {
			t.Errorf("recorded = %+v, want one folder note", h.state.records)
		}
	})
}

// Every path is attempted; the failures are joined, and nothing is recorded
// because the operation did not fully happen.
func TestCreateFileReportsEveryFailure(t *testing.T) {
	first := errors.New("first")
	second := errors.New("second")

	h := newHarness()
	h.creator.errs["a.txt"] = first
	h.creator.errs["c.txt"] = second

	err := h.run(createcommand.Request{Mode: createcommand.ModeFile, Args: []string{"a.txt", "b.txt", "c.txt"}})

	if !errors.Is(err, first) || !errors.Is(err, second) {
		t.Fatalf("error = %v, want both failures", err)
	}

	if got, want := h.creator.calls, []string{"file a.txt", "file b.txt", "file c.txt"}; !slices.Equal(got, want) {
		t.Errorf("creator calls = %v, want every path attempted", got)
	}

	if len(h.state.records) != 0 {
		t.Error("a partially failed create recorded an entry")
	}
}

func TestCreatePathWithoutArgumentsFails(t *testing.T) {
	tests := []struct {
		mode createcommand.Mode
		want string
	}{
		{createcommand.ModeFile, "missing arguments: file path"},
		{createcommand.ModeFolder, "missing arguments: folder path"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			h := newHarness()

			err := h.run(createcommand.Request{Mode: tt.mode})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}

			if len(h.creator.calls) != 0 || len(h.state.records) != 0 {
				t.Error("a missing path still did work")
			}
		})
	}
}

func TestCreateSystemVariable(t *testing.T) {
	h := newHarness()

	if err := h.run(createcommand.Request{Mode: createcommand.ModeSystemVariable, Args: []string{"FOO", "bar"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.environ.calls, []string{"variable FOO=bar"}; !slices.Equal(got, want) {
		t.Errorf("environment calls = %v, want %v", got, want)
	}

	if !strings.Contains(h.output.text(), "environment changed") {
		t.Errorf("output = %q, want the adapter's message", h.output.text())
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:system:[FOO bar]" {
		t.Errorf("recorded = %+v, want the system note", h.state.records)
	}
}

func TestCreateSystemVariableNeedsNameAndValue(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{Mode: createcommand.ModeSystemVariable, Args: []string{"FOO"}})
	if err == nil || !strings.Contains(err.Error(), "missing arguments: var name and var value") {
		t.Fatalf("error = %v, want the missing-arguments error", err)
	}

	if len(h.environ.calls) != 0 {
		t.Error("the environment was edited without a value")
	}
}

func TestCreateSystemPath(t *testing.T) {
	h := newHarness()

	if err := h.run(createcommand.Request{Mode: createcommand.ModeSystemPath, Args: []string{"/opt/tools"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.environ.calls, []string{"path /opt/tools"}; !slices.Equal(got, want) {
		t.Errorf("environment calls = %v, want %v", got, want)
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:syspath:[/opt/tools]" {
		t.Errorf("recorded = %+v, want the syspath note", h.state.records)
	}
}

func TestCreateSystemPathWithoutAnArgumentFails(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{Mode: createcommand.ModeSystemPath})
	if err == nil || !strings.Contains(err.Error(), "missing arguments: path") {
		t.Fatalf("error = %v, want the missing-path error", err)
	}
}

func TestCreateUnknownModeFails(t *testing.T) {
	h := newHarness()

	err := h.run(createcommand.Request{Mode: createcommand.Mode(99)})
	if err == nil || !strings.Contains(err.Error(), "unknown create mode") {
		t.Fatalf("error = %v, want an unknown-mode error", err)
	}

	if len(h.state.saved) != 0 || len(h.state.records) != 0 {
		t.Error("an unknown mode did work")
	}
}

func TestCreateWriteErrorsFailTheRun(t *testing.T) {
	boom := errors.New("boom")

	t.Run("save", func(t *testing.T) {
		h := newHarness()
		h.state.saveErr = boom

		err := h.run(createcommand.Request{Mode: createcommand.ModeCommand, Args: []string{"greet", "echo hi"}})
		if !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed save still recorded an entry")
		}
	})

	t.Run("record", func(t *testing.T) {
		h := newHarness()
		h.state.recordErr = boom

		err := h.run(createcommand.Request{Mode: createcommand.ModeFile, Args: []string{"a.txt"}})
		if !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

func TestCreateEnvironmentErrorFailsTheRun(t *testing.T) {
	boom := errors.New("boom")

	h := newHarness()
	h.environ.setErr = boom

	err := h.run(createcommand.Request{Mode: createcommand.ModeSystemVariable, Args: []string{"FOO", "bar"}})
	if !errors.Is(err, boom) {
		t.Fatalf("error = %v, want boom", err)
	}

	if len(h.state.records) != 0 {
		t.Error("a failed environment write still recorded an entry")
	}
}
