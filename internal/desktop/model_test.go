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

// A page is newest first, bounded, and carries the command (empty when the
// operation touched none).
func TestHistoryPageIsNewestFirstAndBounded(t *testing.T) {
	var history timeline.Log
	for i := 0; i < 5; i++ {
		history.Add(operation.Create, "greet", "note")
	}
	history.Add(operation.Copy, "", "note")

	m := modelWithCommands(nil, history)
	page := m.HistoryPage(0, 2)

	if len(page) != 2 {
		t.Fatalf("page = %d rows, want 2", len(page))
	}

	if page[0].Version != 6 || page[0].Operation != "copy" || page[0].Command != "" {
		t.Errorf("page[0] = %+v, want the newest entry with no command", page[0])
	}

	if page[1].Version != 5 || page[1].Command != "greet" {
		t.Errorf("page[1] = %+v, want the create entry for greet", page[1])
	}

	if len(page[0].When) != len("2006-01-02 15:04") {
		t.Errorf("When = %q, want a formatted local timestamp", page[0].When)
	}

	if last := m.HistoryPage(5, 10); len(last) != 1 || last[0].Version != 1 {
		t.Errorf("last page = %+v, want the oldest entry only", last)
	}

	if empty := m.HistoryPage(20, 10); len(empty) != 0 {
		t.Errorf("out-of-range page = %+v, want empty", empty)
	}

	if empty := m.HistoryPage(0, 0); len(empty) != 0 {
		t.Errorf("zero-limit page = %+v, want empty", empty)
	}
}

// A timestamp that does not parse is shown as stored rather than dropped.
func TestHistoryPageKeepsAnUnparseableTimestamp(t *testing.T) {
	history := timeline.Log{AllData: []timeline.Content{{
		Version:   1,
		Timestemp: "not a timestamp",
		Operation: operation.Create,
	}}}

	m := modelWithCommands(nil, history)
	page := m.HistoryPage(0, 1)

	if len(page) != 1 || page[0].When != "not a timestamp" {
		t.Fatalf("page = %+v, want the raw timestamp", page)
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

// The detail is read on demand, by version, and knows its position from the
// newest entry so Restore can count the steps.
func TestHistoryDetailReturnsTheEntry(t *testing.T) {
	var history timeline.Log
	history.Add(operation.Copy, "", "main_operation=copy;sub_operation=file")
	history.Add(operation.Execute, "greet",
		`main_operation=execute;sub_operation=command;command_name=greet;command_exec=echo hi;command_params=["a b","-x"]`)

	m := modelWithCommands([]catalog.Command{{Name: "greet", Exec: "echo hi"}}, history)

	detail, found := m.HistoryDetail(2)
	if !found {
		t.Fatal("HistoryDetail(2) = missing")
	}

	if detail.Version != 2 || detail.Position != 0 {
		t.Errorf("detail version/position = %d/%d, want 2/0", detail.Version, detail.Position)
	}

	if len(detail.Fields) != 5 || detail.Fields[2] != (HistoryField{Key: "command_name", Value: "greet"}) {
		t.Errorf("Fields = %v, want the note in order with decoded values", detail.Fields)
	}

	if !slices.Equal(detail.Params, []string{"a b", "-x"}) || !detail.ParamsRecorded {
		t.Errorf("Params/recorded = %v/%v, want the recorded arguments", detail.Params, detail.ParamsRecorded)
	}

	if !strings.Contains(detail.Note, "command_exec=echo hi") {
		t.Errorf("Note = %q, want the decoded line", detail.Note)
	}

	older, found := m.HistoryDetail(1)
	if !found || older.Position != 1 {
		t.Fatalf("detail(1) = %+v, %v, want position 1", older, found)
	}

	if older.Params != nil || older.ParamsRecorded {
		t.Errorf("Params/recorded = %v/%v, want none recorded", older.Params, older.ParamsRecorded)
	}

	if _, found := m.HistoryDetail(99); found {
		t.Error("HistoryDetail(99) found an entry that is not there")
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
