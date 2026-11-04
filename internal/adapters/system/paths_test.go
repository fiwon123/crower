package system_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
)

func TestPathsIsFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "a.txt")

	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	paths := system.NewPaths()

	if !paths.IsFile(file) {
		t.Errorf("IsFile(%s) = false, want true", file)
	}

	if paths.IsFile(dir) {
		t.Error("IsFile(folder) = true, want false")
	}

	if paths.IsFile(filepath.Join(dir, "missing")) {
		t.Error("IsFile(missing) = true, want false")
	}
}

// The transfer has to run, not only build: a path with a space is the case the
// quoting exists for.
func TestTransferCopiesAndMovesAFileWithASpace(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "a file.txt")
	destination := filepath.Join(dir, "out")

	if err := os.MkdirAll(destination, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	folder := filepath.Join(dir, "folder")
	if err := os.MkdirAll(folder, 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := os.WriteFile(source, []byte("x"), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	transfer := system.NewTransfer()

	if err := transfer.Copy(source, destination, false); err != nil {
		t.Fatalf("Copy: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destination, "a file.txt")); err != nil {
		t.Fatalf("copied file: %v", err)
	}

	if err := transfer.Move(filepath.Join(destination, "a file.txt"), folder); err != nil {
		t.Fatalf("Move: %v", err)
	}

	if _, err := os.Stat(filepath.Join(folder, "a file.txt")); err != nil {
		t.Errorf("moved file: %v", err)
	}

	if _, err := os.Stat(filepath.Join(destination, "a file.txt")); !os.IsNotExist(err) {
		t.Errorf("the source still exists: %v", err)
	}
}
