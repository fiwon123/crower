package main

import (
	"fmt"
	"os"

	"github.com/fiwon123/crower/internal/tools/winpe"
)

// verifyPath checks a built executable against what its name promises: the
// window binary is a GUI program with the icon, manifest and version; the CLI
// pair are console programs.
func verifyPath(path string, gui bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := winpe.Inspect(file)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	if gui {
		if err := info.VerifyGUI(); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}

		return nil
	}

	if err := info.VerifyConsole(); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	return nil
}
