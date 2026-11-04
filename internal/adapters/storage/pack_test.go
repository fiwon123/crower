package storage_test

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/domain/catalog"
)

func TestPacksRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pack.toml")

	written := []catalog.Command{
		{Name: "greet", Aliases: []string{"hi", "hello"}, Exec: "echo hi", Description: "Say hi", Params: []string{"name"}},
		{Name: "multi", Exec: "unused", Steps: []string{"echo a", "echo b"}},
	}

	if err := storage.NewPacks().Write(path, written); err != nil {
		t.Fatalf("Write: %v", err)
	}

	read, err := storage.NewPacks().Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	byName := map[string]catalog.Command{}
	for _, command := range read {
		byName[command.Name] = command
	}

	greet := byName["greet"]
	if !slices.Equal(greet.Aliases, []string{"hi", "hello"}) || greet.Exec != "echo hi" {
		t.Errorf("greet = %+v, want the aliases and exec", greet)
	}

	if greet.Description != "Say hi" || !slices.Equal(greet.Params, []string{"name"}) {
		t.Errorf("greet = %+v, want the description and params", greet)
	}

	if multi := byName["multi"]; !slices.Equal(multi.Steps, []string{"echo a", "echo b"}) {
		t.Errorf("multi = %+v, want the steps", multi)
	}
}

// The section key is the command's name everywhere else, so it wins over a
// stored Name that disagrees with it.
func TestPacksReadUsesTheSectionKeyAsTheName(t *testing.T) {
	path := filepath.Join(t.TempDir(), "pack.toml")

	body := "[greet]\nName = \"other\"\nExec = \"echo hi\"\n"
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	read, err := storage.NewPacks().Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}

	if len(read) != 1 || read[0].Name != "greet" {
		t.Errorf("read = %+v, want the key as the name", read)
	}
}

func TestPacksReadAMissingFile(t *testing.T) {
	_, err := storage.NewPacks().Read(filepath.Join(t.TempDir(), "missing.toml"))
	if err == nil || !strings.Contains(err.Error(), "read pack") {
		t.Fatalf("error = %v, want read pack", err)
	}
}
