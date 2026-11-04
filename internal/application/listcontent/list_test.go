package listcontent_test

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/application/listcontent"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type fakeState struct {
	commands  catalog.Catalog
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

func (s *fakeState) Record(op operation.Operation, commandName, note string) error {
	if s.recordErr != nil {
		return s.recordErr
	}

	s.records = append(s.records, record{op: op, commandName: commandName, note: note})
	return nil
}

type emitted struct {
	text    string
	payload any
}

type fakeOutput struct {
	table      string
	structured bool
	emitted    []emitted
	emitErr    error
}

func (o *fakeOutput) Table(catalog.Catalog) string {
	return o.table
}

func (o *fakeOutput) Emit(text string, payload any) error {
	if o.emitErr != nil {
		return o.emitErr
	}

	o.emitted = append(o.emitted, emitted{text: text, payload: payload})
	return nil
}

func (o *fakeOutput) Structured() bool {
	return o.structured
}

type fakeHistory struct {
	table string
	rows  []listcontent.HistoryRow
}

func (h *fakeHistory) HistoryTable() string {
	return h.table
}

func (h *fakeHistory) HistoryRows() []listcontent.HistoryRow {
	return h.rows
}

type fakeProcesses struct {
	rows  []listcontent.ProcessRow
	err   error
	asked []string
}

func (p *fakeProcesses) List(part string) ([]listcontent.ProcessRow, error) {
	p.asked = append(p.asked, part)
	return p.rows, p.err
}

type fakeFiles struct {
	text       string
	textErr    error
	entries    []listcontent.FolderEntry
	entriesErr error
	textCalls  []string
	entryCalls []string
}

func (f *fakeFiles) FolderText(path string) (string, error) {
	f.textCalls = append(f.textCalls, path)
	return f.text, f.textErr
}

func (f *fakeFiles) FolderEntries(path string) ([]listcontent.FolderEntry, error) {
	f.entryCalls = append(f.entryCalls, path)
	return f.entries, f.entriesErr
}

type fakeEnvironment struct {
	variables    string
	variablesErr error
	path         string
	pathErr      error
}

func (e *fakeEnvironment) Variables() (string, error) {
	return e.variables, e.variablesErr
}

func (e *fakeEnvironment) Path() (string, error) {
	return e.path, e.pathErr
}

type fakeNotes struct{}

func (fakeNotes) ListCommandsNote() string        { return "note:commands" }
func (fakeNotes) ListProcessNote([]string) string { return "note:process" }
func (fakeNotes) ListHistoryNote() string         { return "note:history" }
func (fakeNotes) ListFolderNote([]string) string  { return "note:folder" }
func (fakeNotes) ListSystemNote() string          { return "note:system" }
func (fakeNotes) ListPathNote() string            { return "note:path" }

type harness struct {
	state  *fakeState
	output *fakeOutput
	hist   *fakeHistory
	procs  *fakeProcesses
	files  *fakeFiles
	env    *fakeEnvironment
}

func newHarness(commands ...catalog.Command) *harness {
	return &harness{
		state:  &fakeState{commands: catalog.New(nil, commands)},
		output: &fakeOutput{table: "TABLE"},
		hist:   &fakeHistory{table: "HISTORY"},
		procs:  &fakeProcesses{},
		files:  &fakeFiles{},
		env:    &fakeEnvironment{},
	}
}

func (h *harness) run(request listcontent.Request) error {
	return listcontent.New(h.state, h.output, h.hist, h.procs, h.files, h.env, fakeNotes{}).Run(request)
}

func TestListCommands(t *testing.T) {
	h := newHarness()
	h.state.commands = catalog.New(
		[]string{"greet", "build"},
		[]catalog.Command{
			{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi", Description: "Say hi", Params: []string{"name"}, Shell: "pwsh", Cwd: "~/api", Env: map[string]string{"STAGE": "prod"}, Timeout: "30s"},
			{Name: "build", Exec: "go build ./..."},
		},
	)

	if err := h.run(listcontent.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if len(h.output.emitted) != 1 || h.output.emitted[0].text != "TABLE" {
		t.Fatalf("emitted = %+v, want the table text", h.output.emitted)
	}

	payload, ok := h.output.emitted[0].payload.(listcontent.CommandsResult)
	if !ok {
		t.Fatalf("payload = %T, want CommandsResult", h.output.emitted[0].payload)
	}

	if len(payload.Commands) != 2 {
		t.Fatalf("payload holds %d rows, want 2", len(payload.Commands))
	}

	if payload.Commands[0].Row != 0 || payload.Commands[0].Name != "greet" {
		t.Errorf("first row = %+v, want greet at row 0", payload.Commands[0])
	}

	if payload.Commands[0].Description != "Say hi" || !slices.Equal(payload.Commands[0].Params, []string{"name"}) {
		t.Errorf("metadata = %+v, want the description and params", payload.Commands[0])
	}

	if payload.Commands[0].Shell != "pwsh" || payload.Commands[0].Cwd != "~/api" {
		t.Errorf("target = %+v, want the shell and cwd", payload.Commands[0])
	}

	if payload.Commands[0].Env["STAGE"] != "prod" || payload.Commands[0].Timeout != "30s" {
		t.Errorf("target = %+v, want the env and timeout", payload.Commands[0])
	}

	// An empty alias list encodes as [] rather than null.
	if payload.Commands[1].Aliases == nil {
		t.Error("aliases = nil, want an empty slice")
	}

	if len(h.state.records) != 1 || h.state.records[0].note != "note:commands" {
		t.Errorf("records = %+v, want the commands note", h.state.records)
	}

	if h.state.records[0].op != operation.List || h.state.records[0].commandName != "" {
		t.Errorf("recorded %+v, want list with no command name", h.state.records[0])
	}
}

// The write path refuses a literal under a secret-looking key, but an old or
// hand-edited config still loads. The JSON result must not be the one place
// the literal is printed; a ${VAR} reference stays readable.
func TestListCommandsRedactsASecretLiteral(t *testing.T) {
	h := newHarness()
	h.state.commands = catalog.New(nil, []catalog.Command{{
		Name: "deploy",
		Exec: "echo hi",
		Env: map[string]string{
			"STAGE":     "prod",
			"TOKEN":     "${TOKEN}",
			"API_TOKEN": "literal",
		},
	}})

	if err := h.run(listcontent.Request{}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	payload, ok := h.output.emitted[0].payload.(listcontent.CommandsResult)
	if !ok {
		t.Fatalf("payload = %T, want CommandsResult", h.output.emitted[0].payload)
	}

	env := payload.Commands[0].Env

	if env["STAGE"] != "prod" || env["TOKEN"] != "${TOKEN}" {
		t.Errorf("env = %v, want the ordinary literal and the reference", env)
	}

	if env["API_TOKEN"] != "${REDACTED}" {
		t.Errorf("API_TOKEN = %q, want the literal redacted", env["API_TOKEN"])
	}
}

func TestListProcess(t *testing.T) {
	h := newHarness()
	h.procs.rows = []listcontent.ProcessRow{{PID: 7, Name: "go", Path: "/usr/bin/go"}}

	if err := h.run(listcontent.Request{Mode: listcontent.ModeProcess, Args: []string{"go"}}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if got, want := h.procs.asked, []string{"go"}; !slices.Equal(got, want) {
		t.Errorf("filter = %v, want %v", got, want)
	}

	text := h.output.emitted[0].text
	if !strings.Contains(text, "PID") || !strings.Contains(text, "go") {
		t.Errorf("text = %q, want the process table", text)
	}

	payload, ok := h.output.emitted[0].payload.(listcontent.ProcessesResult)
	if !ok || len(payload.Processes) != 1 || payload.Processes[0].PID != 7 {
		t.Errorf("payload = %+v, want the process row", h.output.emitted[0].payload)
	}
}

func TestListProcessWithNoMatches(t *testing.T) {
	h := newHarness()

	if err := h.run(listcontent.Request{Mode: listcontent.ModeProcess}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	payload := h.output.emitted[0].payload.(listcontent.ProcessesResult)
	if payload.Processes == nil || len(payload.Processes) != 0 {
		t.Errorf("processes = %#v, want an empty slice", payload.Processes)
	}
}

func TestListHistory(t *testing.T) {
	h := newHarness()
	h.hist.rows = []listcontent.HistoryRow{{Version: 2, Note: "greet"}}

	if err := h.run(listcontent.Request{Mode: listcontent.ModeHistory}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if h.output.emitted[0].text != "HISTORY" {
		t.Errorf("text = %q, want the history table", h.output.emitted[0].text)
	}

	payload := h.output.emitted[0].payload.(listcontent.HistoryResult)
	if len(payload.History) != 1 || payload.History[0].Version != 2 {
		t.Errorf("payload = %+v, want the history row", payload)
	}

	if h.state.records[0].note != "note:history" {
		t.Errorf("note = %q, want the history note", h.state.records[0].note)
	}
}

// The text form runs the platform listing; the structured form reads the
// directory and never touches the shell.
func TestListFolderChoosesTheSourceByFormat(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		h := newHarness()
		h.files.text = "ls output"

		if err := h.run(listcontent.Request{Mode: listcontent.ModeFolder, Args: []string{"~/p"}}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if got, want := h.files.textCalls, []string{"~/p"}; !slices.Equal(got, want) {
			t.Errorf("folder text calls = %v, want %v", got, want)
		}

		if len(h.files.entryCalls) != 0 {
			t.Errorf("the text form read the directory: %v", h.files.entryCalls)
		}

		if h.output.emitted[0].text != "ls output" || h.output.emitted[0].payload != nil {
			t.Errorf("emitted = %+v, want the listing text and no payload", h.output.emitted[0])
		}
	})

	t.Run("structured", func(t *testing.T) {
		h := newHarness()
		h.output.structured = true
		h.files.entries = []listcontent.FolderEntry{{Name: "a.txt"}}

		if err := h.run(listcontent.Request{Mode: listcontent.ModeFolder}); err != nil {
			t.Fatalf("Run: %v", err)
		}

		if len(h.files.textCalls) != 0 {
			t.Errorf("the structured form ran the shell listing: %v", h.files.textCalls)
		}

		if got, want := h.files.entryCalls, []string{"./"}; !slices.Equal(got, want) {
			t.Errorf("folder entry calls = %v, want the default path", got)
		}

		payload := h.output.emitted[0].payload.(listcontent.FolderResult)
		if payload.Folder != "./" || len(payload.Entries) != 1 {
			t.Errorf("payload = %+v, want the folder and its entry", payload)
		}

		if h.output.emitted[0].text != "" {
			t.Errorf("text = %q, want empty in structured mode", h.output.emitted[0].text)
		}
	})
}

func TestListSystem(t *testing.T) {
	h := newHarness()
	h.env.variables = "FOO=bar\nEMPTY=\nBROKEN\n"

	if err := h.run(listcontent.Request{Mode: listcontent.ModeSystem}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	text := h.output.emitted[0].text
	for _, want := range []string{"FOO\n0- bar\n", "EMPTY\n0- \n"} {
		if !strings.Contains(text, want) {
			t.Errorf("text = %q, want it to contain %q", text, want)
		}
	}

	if strings.Contains(text, "BROKEN") {
		t.Errorf("text = %q, want a line without = to be skipped", text)
	}

	payload := h.output.emitted[0].payload.(listcontent.VariablesResult)
	if len(payload.Variables) != 2 {
		t.Fatalf("variables = %+v, want two", payload.Variables)
	}

	if payload.Variables[1].Values == nil || len(payload.Variables[1].Values) != 0 {
		t.Errorf("empty variable values = %#v, want an empty slice", payload.Variables[1].Values)
	}
}

func TestListPath(t *testing.T) {
	h := newHarness()
	h.env.path = "/usr/bin" + string(os.PathListSeparator) + "/opt" + "\n"

	if err := h.run(listcontent.Request{Mode: listcontent.ModePath}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	text := h.output.emitted[0].text
	if !strings.Contains(text, "PATH\n0- /usr/bin\n1- /opt\n") {
		t.Errorf("text = %q, want the PATH entries", text)
	}

	payload := h.output.emitted[0].payload.(listcontent.VariableResult)
	if payload.Variable.Name != "PATH" || len(payload.Variable.Values) != 2 {
		t.Errorf("payload = %+v, want two PATH values", payload)
	}
}

func TestListErrors(t *testing.T) {
	boom := errors.New("boom")

	t.Run("process", func(t *testing.T) {
		h := newHarness()
		h.procs.err = boom

		if err := h.run(listcontent.Request{Mode: listcontent.ModeProcess}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed listing recorded an entry")
		}
	})

	t.Run("folder text", func(t *testing.T) {
		h := newHarness()
		h.files.textErr = boom

		if err := h.run(listcontent.Request{Mode: listcontent.ModeFolder}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("folder entries", func(t *testing.T) {
		h := newHarness()
		h.output.structured = true
		h.files.entriesErr = boom

		if err := h.run(listcontent.Request{Mode: listcontent.ModeFolder}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("variables", func(t *testing.T) {
		h := newHarness()
		h.env.variablesErr = boom

		if err := h.run(listcontent.Request{Mode: listcontent.ModeSystem}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("path", func(t *testing.T) {
		h := newHarness()
		h.env.pathErr = boom

		if err := h.run(listcontent.Request{Mode: listcontent.ModePath}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})

	t.Run("emit", func(t *testing.T) {
		h := newHarness()
		h.output.emitErr = boom

		if err := h.run(listcontent.Request{}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}

		if len(h.state.records) != 0 {
			t.Error("a failed result still recorded an entry")
		}
	})

	t.Run("record", func(t *testing.T) {
		h := newHarness()
		h.state.recordErr = boom

		if err := h.run(listcontent.Request{}); !errors.Is(err, boom) {
			t.Fatalf("error = %v, want boom", err)
		}
	})
}

func TestListUnknownMode(t *testing.T) {
	h := newHarness()

	err := h.run(listcontent.Request{Mode: listcontent.Mode(99)})
	if err == nil || !strings.Contains(err.Error(), "unknown list mode") {
		t.Fatalf("error = %v, want an unknown-mode error", err)
	}

	if len(h.output.emitted) != 0 || len(h.state.records) != 0 {
		t.Error("an unknown mode did work")
	}
}
