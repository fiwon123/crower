package app

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"

	"github.com/fiwon123/crower/pkg/utils"
)

func newTestConfig(t *testing.T) *Config {
	t.Helper()

	cfg, err := InitApp(
		filepath.Join(t.TempDir(), "crower.yaml"),
		Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}},
	)
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	return cfg
}

func TestConfigPaths(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "nested", "crower.yaml")

	cfg, err := InitApp(cfgPath, Options{Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard}})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

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
	cfg := newTestConfig(t)

	cfg.AllCommandsByName.Add("greet", command.New("greet", []string{"hi"}, "echo hello"))

	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("SaveCommands: %v", err)
	}

	read := command.NewMapData()
	if err := utils.ReadToml(cfg.ConfigFile(), &read); err != nil {
		t.Fatalf("read back %s: %v", cfg.ConfigFile(), err)
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

	cfg.History.Add(operations.List, "note")
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
