package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

func commandMap(t *testing.T) command.MapData {
	t.Helper()

	commands := command.NewMapData()
	commands.Add("greet", command.New("greet", []string{"hi"}, "echo hello"))
	commands.Add("build", &command.Data{
		Name:     "build",
		AllAlias: []string{},
		Exec:     "go build ./...",
		Steps:    []string{"go vet ./...", "go test ./..."},
	})

	return commands
}

// Load prepares the folder and the file, so a first run starts from an empty
// set instead of an error.
func TestCommandFileLoadCreatesTheInstallation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "crower.yaml")

	commands, order, err := storage.NewCommandFile(path).Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(commands) != 0 {
		t.Errorf("commands = %v, want empty", commands)
	}

	if len(order) != 0 {
		t.Errorf("order = %v, want empty", order)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("the config file was not created: %v", err)
	}
}

func TestCommandFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")
	file := storage.NewCommandFile(path)

	want := commandMap(t)

	if err := file.Save(want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, order, err := file.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(order) != len(want) {
		t.Errorf("order = %v, want %d keys", order, len(want))
	}

	for name, expected := range want {
		stored := got.Get(name)
		if stored == nil {
			t.Fatalf("Load lost %q", name)
		}

		if stored.Exec != expected.Exec {
			t.Errorf("%s: Exec = %q, want %q", name, stored.Exec, expected.Exec)
		}

		if len(stored.Steps) != len(expected.Steps) {
			t.Errorf("%s: Steps = %v, want %v", name, stored.Steps, expected.Steps)
		}
	}
}

func TestHistoryFilesEnsureAndLoadEmpty(t *testing.T) {
	folder := t.TempDir()
	files := storage.NewHistoryFiles(folder)

	if err := files.Ensure(); err != nil {
		t.Fatalf("Ensure: %v", err)
	}

	for _, path := range []string{
		filepath.Join(folder, "historyhelper.json"),
		filepath.Join(folder, "history"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("Ensure did not create %s: %v", path, err)
		}
	}

	data, err := files.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(data.AllData) != 0 {
		t.Errorf("AllData = %v, want empty", data.AllData)
	}
}

func TestHistoryIndexRoundTrip(t *testing.T) {
	folder := t.TempDir()
	files := storage.NewHistoryFiles(folder)

	if err := files.Ensure(); err != nil {
		t.Fatalf("Ensure: %v", err)
	}

	data := timeline.Log{}
	data.Add(operation.Create, "greet", "main_operation=create")

	if err := files.SaveIndex(data); err != nil {
		t.Fatalf("SaveIndex: %v", err)
	}

	back, err := files.Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(back.AllData) != 1 || back.AllData[0].CommandName != "greet" {
		t.Fatalf("Load = %+v, want the saved entry", back.AllData)
	}
}

func TestSnapshotRoundTrip(t *testing.T) {
	files := storage.NewHistoryFiles(t.TempDir())

	if err := files.Ensure(); err != nil {
		t.Fatalf("Ensure: %v", err)
	}

	if err := files.WriteSnapshot("00001.yaml", commandMap(t)); err != nil {
		t.Fatalf("WriteSnapshot: %v", err)
	}

	if !files.HasSnapshot("00001.yaml") {
		t.Fatal("HasSnapshot reported the file missing")
	}

	back, err := files.ReadSnapshot("00001.yaml")
	if err != nil {
		t.Fatalf("ReadSnapshot: %v", err)
	}

	if got := back.Get("greet"); got == nil {
		t.Fatal("ReadSnapshot lost the command")
	}

	if err := files.DeleteSnapshot("00001.yaml"); err != nil {
		t.Fatalf("DeleteSnapshot: %v", err)
	}

	if files.HasSnapshot("00001.yaml") {
		t.Fatal("the snapshot survived DeleteSnapshot")
	}
}

// A history file is JSON in a user's home directory and may be edited by hand,
// so the name is checked rather than trusted: anything with a path separator
// would reach outside the snapshot folder.
func TestSnapshotNameCannotEscapeTheFolder(t *testing.T) {
	files := storage.NewHistoryFiles(t.TempDir())

	if err := files.Ensure(); err != nil {
		t.Fatalf("Ensure: %v", err)
	}

	for _, name := range []string{"", "../escape.yaml", "sub/dir.yaml", `..\escape.yaml`} {
		if err := files.WriteSnapshot(name, commandMap(t)); err == nil {
			t.Errorf("WriteSnapshot(%q) was accepted", name)
		}

		if _, err := files.ReadSnapshot(name); err == nil {
			t.Errorf("ReadSnapshot(%q) was accepted", name)
		}

		if err := files.DeleteSnapshot(name); err == nil {
			t.Errorf("DeleteSnapshot(%q) was accepted", name)
		}
	}
}
