package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"

	"github.com/fiwon123/crower/pkg/utils"
)

func newTestConfigAt(t *testing.T, cfgPath string) *Config {
	t.Helper()

	cfg, err := InitApp(
		cfgPath,
		Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}},
	)
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	return cfg
}

func newTestConfig(t *testing.T) *Config {
	t.Helper()

	return newTestConfigAt(t, filepath.Join(t.TempDir(), "crower.yaml"))
}

func TestConfigPaths(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "nested", "crower.yaml")
	cfg := newTestConfigAt(t, cfgPath)

	if got := cfg.ConfigFile(); got != cfgPath {
		t.Errorf("ConfigFile() = %q, want %q", got, cfgPath)
	}

	if got, want := cfg.ConfigFolder(), filepath.Dir(cfgPath); got != want {
		t.Errorf("ConfigFolder() = %q, want %q", got, want)
	}
}

// SaveCommands is the one writer of the configuration file. A domain calls it
// instead of knowing that the file is TOML or where it is.
func TestSaveCommandsWritesTheCommandMap(t *testing.T) {
	// Read back the path handed to InitApp, not ConfigFile(): a wrong path on
	// both sides would otherwise hide itself.
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")
	cfg := newTestConfigAt(t, cfgPath)

	cfg.AllCommandsByName.Add("greet", command.New("greet", []string{"hi"}, "echo hello"))

	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("SaveCommands: %v", err)
	}

	read := command.NewMapData()
	if err := utils.ReadToml(cfgPath, &read); err != nil {
		t.Fatalf("read back %s: %v", cfgPath, err)
	}

	written := read.Get("greet")
	if written == nil {
		t.Fatal("the command map was not written to the configuration file")
	}

	if written.Exec != "echo hello" {
		t.Errorf("exec = %q, want %q", written.Exec, "echo hello")
	}
}

// Deleting a history entry has to remove the snapshot it points at, and the
// caller should not have to know the history folder to do it.
func TestDeleteHistoryBackup(t *testing.T) {
	cfg := newTestConfig(t)

	cfg.History.Add(operation.List, "", "note")
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snapshot := filepath.Join(cfg.historyFiles.Folder(), cfg.History.GetLast().File)
	if _, err := os.Stat(snapshot); err != nil {
		t.Fatalf("Save did not write the snapshot: %v", err)
	}

	if err := cfg.DeleteHistoryBackup(cfg.History.GetLast().File); err != nil {
		t.Fatalf("DeleteHistoryBackup: %v", err)
	}

	if _, err := os.Stat(snapshot); !os.IsNotExist(err) {
		t.Errorf("snapshot still exists after DeleteHistoryBackup: %v", err)
	}
}

// The file name comes from a history entry, which can be hand-edited: it must
// not be able to point outside the history folder.
func TestDeleteHistoryBackupRejectsAPath(t *testing.T) {
	cfg := newTestConfig(t)

	for _, name := range []string{"", "../escape", "/abs/path"} {
		if err := cfg.DeleteHistoryBackup(name); err == nil {
			t.Errorf("DeleteHistoryBackup(%q) returned no error", name)
		}
	}
}

// A history file written by an older build has no CommandName field: the name
// lives only inside the note. Loading it must resolve the name, or `--last`
// stays broken for every entry made before the fix.
func TestNewAppBackfillsTheCommandNameFromTheNote(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "crower.yaml")

	legacy := timeline.Log{}
	legacy.Add(operation.Create, "", "main_operation=create;sub_operation=command;command_name=greet;command_exec=echo hi")

	if err := utils.WriteJson(legacy, filepath.Join(dir, "historyhelper.json")); err != nil {
		t.Fatalf("write legacy history: %v", err)
	}

	cfg, err := InitApp(cfgPath, Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	last := cfg.GetLast(operation.Create)
	if last == nil {
		t.Fatal("the legacy entry did not survive the load")
	}

	if last.CommandName != "greet" {
		t.Errorf("CommandName = %q, want it backfilled from the note", last.CommandName)
	}
}

// A file note can carry the same text inside a path; resolving that would point
// `--last` at a command the operation never touched.
func TestNewAppDoesNotBackfillAFileNote(t *testing.T) {
	dir := t.TempDir()

	legacy := timeline.Log{}
	legacy.Add(operation.Create, "", `main_operation=create;sub_operation=file;crower_exec=crower create --file "x;command_name=greet"`)
	legacy.Add(operation.Execute, "real", "main_operation=execute;sub_operation=command;command_name=other")

	if err := utils.WriteJson(legacy, filepath.Join(dir, "historyhelper.json")); err != nil {
		t.Fatalf("write legacy history: %v", err)
	}

	cfg, err := InitApp(filepath.Join(dir, "crower.yaml"), Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	if got := cfg.History.AllData[0].CommandName; got != "" {
		t.Errorf("file note backfilled %q, want no name", got)
	}

	if got := cfg.History.AllData[1].CommandName; got != "real" {
		t.Errorf("stored name = %q, want it left alone", got)
	}
}

// Commit is the one unit of work a command-changing operation uses: the
// snapshot, the index and the command file all describe the same new state.
func TestCommitWritesSnapshotIndexAndCommandFile(t *testing.T) {
	cfg := newTestConfig(t)

	commands := catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}})

	if err := cfg.Commit(commands, operation.Create, "greet", "main_operation=create"); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if cfg.AllCommandsByName.Get("greet") == nil {
		t.Error("the running configuration does not hold the committed command")
	}

	configBody, err := os.ReadFile(cfg.ConfigFile())
	if err != nil {
		t.Fatalf("read command file: %v", err)
	}

	if !strings.Contains(string(configBody), "greet") {
		t.Errorf("command file = %q, want the committed command", configBody)
	}

	indexBody, err := os.ReadFile(filepath.Join(filepath.Dir(cfg.ConfigFile()), "historyhelper.json"))
	if err != nil {
		t.Fatalf("read history index: %v", err)
	}

	if !strings.Contains(string(indexBody), `"Version": 1`) {
		t.Errorf("history index = %q, want the new entry", indexBody)
	}

	entry := cfg.History.GetLast()
	if entry == nil {
		t.Fatal("Commit did not add a history entry")
	}

	snapshot, err := cfg.historyFiles.ReadSnapshot(entry.File)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}

	if got := snapshot["greet"].Name; got != "greet" {
		t.Errorf("snapshot command = %q, want greet", got)
	}
}

// The second commit starts from what the first wrote: the history version and
// the snapshot name advance together, and nothing is left half-written.
func TestCommitTwice(t *testing.T) {
	cfg := newTestConfig(t)

	first := catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo one"}})
	if err := cfg.Commit(first, operation.Create, "greet", "one"); err != nil {
		t.Fatalf("first Commit: %v", err)
	}

	second, err := first.Add(catalog.Command{Name: "build", Exec: "go build"})
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if err := cfg.Commit(second, operation.Create, "build", "two"); err != nil {
		t.Fatalf("second Commit: %v", err)
	}

	if got := cfg.History.GetLast(); got == nil || got.Version != 2 {
		t.Fatalf("last entry = %+v, want version 2", got)
	}

	if _, err := os.Stat(filepath.Join(cfg.ConfigFolder(), "history", "00002.yaml")); err != nil {
		t.Errorf("snapshot 00002.yaml: %v", err)
	}
}
