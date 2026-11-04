package utils

import (
	"os"
	"path/filepath"
)

// writeFileAtomic writes data through a temporary file in the same directory,
// then renames it over path.
//
// A rename within one directory is atomic, so a reader -- or the next run
// after a crash -- sees either the old file or the new one, never a truncated
// one. The temporary file is created beside the target so the rename cannot
// cross a filesystem, and it is removed if anything fails before the rename.
func writeFileAtomic(path string, data []byte) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}

	name := tmp.Name()

	// A successful rename removes this name; the deferred call is for every
	// failure path in between.
	defer os.Remove(name)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Chmod(0o644); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(name, path)
}
