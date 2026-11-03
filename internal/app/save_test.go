package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/data/operations"

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

	cfg.History.Add(operations.List, "", "note")
	if err := cfg.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	snapshot := filepath.Join(cfg.historyFolderPath, cfg.History.GetLast().File)
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

	legacy := history.New(filepath.Join(dir, "history"))
	legacy.Add(operations.Create, "", "main_operation=create;command_name=greet;command_exec=echo hi")

	if err := utils.WriteJson(legacy, filepath.Join(dir, "historyhelper.json")); err != nil {
		t.Fatalf("write legacy history: %v", err)
	}

	cfg, err := InitApp(cfgPath, Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	last := cfg.GetLast(operations.Create)
	if last == nil {
		t.Fatal("the legacy entry did not survive the load")
	}

	if last.CommandName != "greet" {
		t.Errorf("CommandName = %q, want it backfilled from the note", last.CommandName)
	}
}
