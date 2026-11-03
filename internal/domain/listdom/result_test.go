package listdom

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/testutil"
	"github.com/fiwon123/crower/pkg/utils"
)

func TestCommandsPayload(t *testing.T) {
	all := command.NewMapData()
	all.Add("greet", command.New("greet", []string{"hi", "hello"}, "echo hi"))
	all.Add("bare", command.New("bare", nil, "true"))

	// "missing" has an order key but no entry: the table skips it, so the JSON
	// result must skip it too, without shifting the row numbers of the rest.
	payload := commandsPayload([]string{"greet", "missing", "bare"}, all)

	if got, want := len(payload.Commands), 2; got != want {
		t.Fatalf("commands = %d, want %d", got, want)
	}

	first := payload.Commands[0]
	if first.Row != 0 || first.Name != "greet" || first.Exec != "echo hi" {
		t.Errorf("first row = %+v", first)
	}
	if strings.Join(first.Aliases, ",") != "hi,hello" {
		t.Errorf("aliases = %v", first.Aliases)
	}

	second := payload.Commands[1]
	if second.Row != 2 {
		t.Errorf("row after a skipped entry = %d, want 2 (the table's index)", second.Row)
	}
	if second.Aliases == nil {
		t.Error("a command with no aliases encoded a null list")
	}
}

func TestHistoryPayloadIsNewestFirst(t *testing.T) {
	rows := []history.Content{
		{Version: 1, File: "00001.yaml", Timestemp: "2026-01-01T00:00:00Z", Operation: operation.List, Note: "first"},
		{Version: 2, File: "00002.yaml", Timestemp: "2026-01-02T00:00:00Z", Operation: operation.Create, CommandName: "greet", Note: "second"},
	}

	payload := historyPayload(rows)

	if got, want := len(payload.History), 2; got != want {
		t.Fatalf("history = %d rows, want %d", got, want)
	}

	if payload.History[0].Version != 2 {
		t.Errorf("first row version = %d, want the newest (2)", payload.History[0].Version)
	}

	if payload.History[0].Timestamp != "2026-01-02T00:00:00Z" {
		t.Errorf("timestamp = %q", payload.History[0].Timestamp)
	}

	if payload.History[0].Operation != string(operation.Create) {
		t.Errorf("operation = %q", payload.History[0].Operation)
	}

	if payload.History[0].Command != "greet" {
		t.Errorf("command = %q, want greet", payload.History[0].Command)
	}
}

func TestHistoryPayloadEmptyIsACollection(t *testing.T) {
	if payload := historyPayload(nil); payload.History == nil {
		t.Error("empty history encoded a null list")
	}
}

func TestProcessesPayloadEmptyIsACollection(t *testing.T) {
	if payload := processesPayload(nil); payload.Processes == nil {
		t.Error("empty process list encoded a null list")
	}
}

func TestProcessTableLayout(t *testing.T) {
	table := processTable([]utils.ProcessInfo{{PID: 42, Name: "crower", Path: "/usr/bin/crower"}})

	if !strings.HasPrefix(table, "PID") {
		t.Errorf("table does not start with the header: %q", table)
	}

	if !strings.Contains(table, "42") || !strings.Contains(table, "/usr/bin/crower") {
		t.Errorf("row missing from table: %q", table)
	}

	if !strings.HasSuffix(table, "\n") {
		t.Errorf("table does not end with a newline: %q", table)
	}
}

func TestVariablesPayload(t *testing.T) {
	raw := "HOME=/home/dev\nNO_EQUALS_LINE\nPATH=/bin" + string(os.PathListSeparator) + "/usr/bin\n"

	payload := variablesPayload(raw)

	if got, want := len(payload.Variables), 2; got != want {
		t.Fatalf("variables = %d, want %d (a line without = is not a variable)", got, want)
	}

	if payload.Variables[0].Name != "HOME" || payload.Variables[0].Values[0] != "/home/dev" {
		t.Errorf("HOME = %+v", payload.Variables[0])
	}

	path := payload.Variables[1]
	if len(path.Values) != 2 || path.Values[1] != "/usr/bin" {
		t.Errorf("PATH values = %v", path.Values)
	}
}

func TestVariablePayloadTrimsCarriageReturn(t *testing.T) {
	payload := variablePayload("PATH", "/bin:/usr/bin\r")

	values := payload.Variable.Values
	if len(values) != 2 || strings.HasSuffix(values[1], "\r") {
		t.Errorf("values = %q, want the trailing CR removed", values)
	}
}

func TestFolderResult(t *testing.T) {
	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("InitCrowerTests: %v", err)
	}

	handler := NewHandler(cfg)
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("abc"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	payload, err := handler.FolderResult(dir)
	if err != nil {
		t.Fatalf("FolderResult: %v", err)
	}

	if payload.Folder != dir {
		t.Errorf("folder = %q, want %q", payload.Folder, dir)
	}

	if got, want := len(payload.Entries), 2; got != want {
		t.Fatalf("entries = %d, want %d", got, want)
	}

	// os.ReadDir sorts by name, so the file comes before the directory.
	file := payload.Entries[0]
	if file.Name != "a.txt" || file.Dir || file.Size != 3 || file.Mode == "" || file.Modified == "" {
		t.Errorf("file entry = %+v", file)
	}

	if sub := payload.Entries[1]; !sub.Dir || sub.Name != "sub" {
		t.Errorf("directory entry = %+v", sub)
	}
}

func TestFolderResultMissingFolderIsAnError(t *testing.T) {
	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("InitCrowerTests: %v", err)
	}

	if _, err := NewHandler(cfg).FolderResult(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("FolderResult accepted a folder that does not exist")
	}
}
