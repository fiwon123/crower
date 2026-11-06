package desktop

import (
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// stubExecutor stands in for the command layer's executor: these tests are
// about the model's state, not about running anything.
type stubExecutor struct{}

func (stubExecutor) Execute(string, []string) (string, string, error) {
	return "", "", nil
}

// testModel builds the model over one command and an empty history.
func testModel(t *testing.T) *Model {
	t.Helper()

	commands := catalog.New(nil, []catalog.Command{{Name: "hello", Exec: "echo hello"}})

	return NewModel(commands, timeline.Log{}, stubExecutor{})
}

func modelWithCommands(commands []catalog.Command, history timeline.Log) *Model {
	return NewModel(catalog.New(nil, commands), history, stubExecutor{})
}

func TestNewModel(t *testing.T) {
	m := testModel(t)
	if m.State != ViewList {
		t.Fatalf("want list state")
	}
	if len(m.Keys) != 1 {
		t.Fatalf("want 1 key")
	}
}

func TestMove(t *testing.T) {
	commands := catalog.New(
		[]string{"hello", "b"},
		[]catalog.Command{
			{Name: "hello", Exec: "echo hello"},
			{Name: "b", Exec: "echo b"},
		},
	)

	m := NewModel(commands, timeline.Log{}, stubExecutor{})
	m.Move(1)
	if m.Selected() != "b" {
		t.Fatalf("want b")
	}
	m.Move(1)
	if m.Selected() != "hello" {
		t.Fatalf("want wrap")
	}
}

// recordingExecutor remembers what the model asked it to run.
type recordingExecutor struct {
	keys   []string
	params [][]string
}

func (e *recordingExecutor) Execute(key string, params []string) (string, string, error) {
	e.keys = append(e.keys, key)
	e.params = append(e.params, params)

	return "out", key, nil
}

// A row's Run button runs that row, whether or not it was selected first, and
// the params variant only records the selection until the dialog submits.
func TestStartRunKeyRunsTheChosenRow(t *testing.T) {
	exec := &recordingExecutor{}
	m := NewModel(
		catalog.New(nil, []catalog.Command{
			{Name: "greet", Exec: "echo hi"},
			{Name: "ship", Exec: "echo ship"},
		}),
		timeline.Log{},
		exec,
	)

	m.StartRunKey("ship", false)

	if len(exec.keys) != 1 || exec.keys[0] != "ship" {
		t.Fatalf("executed %v, want [ship]", exec.keys)
	}

	if m.State != ViewOutput || m.LastCmd != "ship" {
		t.Errorf("state/last = %v/%q, want the output view for ship", m.State, m.LastCmd)
	}

	m.StartRunKey("greet", true)

	if m.State != ViewParams || m.SelectedKey != "greet" {
		t.Errorf("state/selected = %v/%q, want the params view for greet", m.State, m.SelectedKey)
	}

	if len(exec.keys) != 1 {
		t.Errorf("parameters path executed %v, want nothing until submit", exec.keys)
	}
}

// The rows are newest first and carry the command, or nothing when the
// operation touched none (the table renders the dash).
func TestHistoryRowsAreNewestFirst(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Create, "greet", "note")
	history.Add(operation.Copy, "", "note")

	m := modelWithCommands([]catalog.Command{{Name: "hello", Exec: "echo hello"}}, history)
	rows := m.HistoryRows()

	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	if rows[0].Version != 2 || rows[0].Operation != "copy" || rows[0].Command != "" {
		t.Errorf("row 0 = %+v, want the newest entry with no command", rows[0])
	}

	if rows[1].Version != 1 || rows[1].Command != "greet" {
		t.Errorf("row 1 = %+v, want the create entry for greet", rows[1])
	}

	if len(rows[0].When) != len("2006-01-02 15:04") {
		t.Errorf("When = %q, want a formatted local timestamp", rows[0].When)
	}
}

// A timestamp that does not parse is shown as stored rather than dropped.
func TestHistoryRowsKeepAnUnparseableTimestamp(t *testing.T) {
	history := timeline.Log{AllData: []timeline.Content{{
		Version:   1,
		Timestemp: "not a timestamp",
		Operation: operation.Create,
	}}}

	m := modelWithCommands(nil, history)
	rows := m.HistoryRows()

	if len(rows) != 1 || rows[0].When != "not a timestamp" {
		t.Fatalf("rows = %+v, want the raw timestamp", rows)
	}
}

// A tools operation changes the log but not the command set, so it replaces
// the history and leaves the list and its selection alone.
func TestReplaceHistoryKeepsTheCommandSet(t *testing.T) {
	m := testModel(t)
	m.Index = 0

	var history timeline.Log
	history.Add(operation.Copy, "", "note")

	m.ReplaceHistory(history)

	if len(m.History.AllData) != 1 {
		t.Fatalf("history = %+v, want the new log", m.History.AllData)
	}

	if got := m.Selected(); got != "hello" {
		t.Errorf("Selected() = %q, want the command set untouched", got)
	}
}

func TestFilter(t *testing.T) {
	m := modelWithCommands([]catalog.Command{
		{Name: "hello", Exec: "echo hello"},
		{Name: "build", Exec: "make build"},
	}, timeline.Log{})
	m.Filter = "bui"
	keys := m.FilteredKeys()
	if len(keys) != 1 || keys[0] != "build" {
		t.Fatalf("bad filter")
	}
}

func TestReloadReplacesTheSetAndResetsTheSelection(t *testing.T) {
	m := testModel(t)
	m.Filter = "hello"
	m.Index = 3

	commands := catalog.New(nil, []catalog.Command{{Name: "build", Exec: "make"}})
	m.Reload(commands, timeline.Log{})

	if got := m.Selected(); got != "build" {
		t.Errorf("Selected() = %q, want build", got)
	}

	if m.Filter != "" || m.Index != 0 {
		t.Errorf("filter/index = %q/%d, want reset", m.Filter, m.Index)
	}
}

// The history screen names the command an entry touched, and shows a dash for
// an operation that touched no command.
func TestHistoryTextShowsTheCommand(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Create, "greet", "note")
	history.Add(operation.Copy, "", "note")

	m := modelWithCommands([]catalog.Command{{Name: "hello", Exec: "echo hello"}}, history)
	text := m.historyText()

	if !strings.Contains(text, "[create] greet") {
		t.Errorf("history = %q, want the command on the create row", text)
	}

	if !strings.Contains(text, "[copy] -") {
		t.Errorf("history = %q, want a dash on the command-less row", text)
	}
}

// The group filter narrows the list to one group; an empty group shows nothing
// rather than falling back to the whole set the way the text filter does.
func TestGroupFilterAndGroups(t *testing.T) {
	m := modelWithCommands([]catalog.Command{
		{Name: "ship", Exec: "x", Group: "deploy"},
		{Name: "rollback", Exec: "x", Group: "deploy"},
		{Name: "backup", Exec: "x", Group: "backup"},
		{Name: "loose", Exec: "x"},
	}, timeline.Log{})

	if got := m.Groups(); !slices.Equal(got, []string{"backup", "deploy"}) {
		t.Errorf("Groups() = %v, want [backup deploy]", got)
	}

	m.Group = "deploy"
	if got := m.FilteredKeys(); !slices.Equal(got, []string{"rollback", "ship"}) {
		t.Errorf("group deploy = %v, want [rollback ship]", got)
	}

	m.Group = "missing"
	if got := m.FilteredKeys(); len(got) != 0 {
		t.Errorf("empty group = %v, want nothing", got)
	}

	m.Group = ""
	m.Filter = "zzz"
	if got := m.FilteredKeys(); len(got) != 4 {
		t.Errorf("text filter fallback = %v, want the whole set", got)
	}
}

// A row carries the note split into fields and the parameters an execute entry
// recorded, so the window can label the entry and run it again.
func TestHistoryRowsExposeFieldsAndParams(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Copy, "", "main_operation=copy;sub_operation=file")
	history.Add(operation.Execute, "greet",
		`main_operation=execute;sub_operation=command;command_name=greet;command_exec=echo hi;command_params=["a b","-x"]`)

	m := modelWithCommands([]catalog.Command{{Name: "greet", Exec: "echo hi"}}, history)
	rows := m.HistoryRows()

	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}

	execute := rows[0]
	if len(execute.Fields) != 5 || execute.Fields[2] != (HistoryField{Key: "command_name", Value: "greet"}) {
		t.Errorf("Fields = %v, want the note in order with decoded values", execute.Fields)
	}

	if !slices.Equal(execute.Params, []string{"a b", "-x"}) {
		t.Errorf("Params = %v, want the recorded arguments", execute.Params)
	}

	if !execute.ParamsRecorded {
		t.Error("ParamsRecorded = false, want true for a note that carries the field")
	}

	if rows[1].Params != nil || rows[1].ParamsRecorded {
		t.Errorf("Params/recorded on a copy row = %v/%v, want none recorded", rows[1].Params, rows[1].ParamsRecorded)
	}
}

// The history Run action replays the entry's arguments and does not keep a
// reference the caller could mutate afterwards.
func TestRunWithParamsUsesTheRecordedArguments(t *testing.T) {
	exec := &recordingExecutor{}
	m := NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}),
		timeline.Log{},
		exec,
	)

	params := []string{"a b", "-x"}
	m.RunWithParams("greet", params)
	params[0] = "changed"

	if !slices.Equal(exec.keys, []string{"greet"}) {
		t.Fatalf("executed %v, want [greet]", exec.keys)
	}

	if !slices.Equal(exec.params[0], []string{"a b", "-x"}) {
		t.Errorf("params = %v, want the recorded arguments", exec.params[0])
	}

	if m.State != ViewOutput || m.SelectedKey != "greet" {
		t.Errorf("state/selected = %v/%q, want the output view for greet", m.State, m.SelectedKey)
	}
}
