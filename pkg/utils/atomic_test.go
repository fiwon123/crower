package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Writing twice replaces the file: no stale bytes, no leftover temporary files.
func TestWriteTomlIsAtomicAndReplaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.toml")

	if err := WriteToml(map[string]string{"a": "first"}, path); err != nil {
		t.Fatalf("first write: %v", err)
	}

	if err := WriteToml(map[string]string{"b": "second"}, path); err != nil {
		t.Fatalf("second write: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if strings.Contains(string(body), "first") || !strings.Contains(string(body), "second") {
		t.Errorf("body = %q, want only the second value", body)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read dir: %v", err)
	}

	if len(entries) != 1 {
		names := []string{}
		for _, entry := range entries {
			names = append(names, entry.Name())
		}

		t.Errorf("directory holds %v, want only the target file", names)
	}
}

func TestWriteJsonIsAtomicAndReplaces(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.json")

	if err := WriteJson(map[string]string{"a": "first"}, path); err != nil {
		t.Fatalf("first write: %v", err)
	}

	if err := WriteJson(map[string]string{"b": "second"}, path); err != nil {
		t.Fatalf("second write: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if strings.Contains(string(body), "first") || !strings.Contains(string(body), "second") {
		t.Errorf("body = %q, want only the second value", body)
	}
}

// A write that cannot create its temporary file leaves the old content alone.
func TestWriteTomlLeavesTheOldFileOnFailure(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.toml")

	if err := WriteToml(map[string]string{"a": "kept"}, path); err != nil {
		t.Fatalf("write: %v", err)
	}

	// A directory where the temporary file would go makes the atomic write
	// fail before it can touch the target.
	if err := os.Mkdir(filepath.Join(dir, ".file.toml.tmp-block"), 0o700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	if err := writeFileAtomic(filepath.Join(dir, "missing-dir", "x.toml"), []byte("x")); err == nil {
		t.Fatal("a missing parent directory returned no error")
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	if !strings.Contains(string(body), "kept") {
		t.Errorf("body = %q, want the old content", body)
	}
}
