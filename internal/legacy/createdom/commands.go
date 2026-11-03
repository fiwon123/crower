package createdom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
)

// The create commands are pure functions of goos so a table test can cover
// every platform the release ships: runtime.GOOS cannot be overridden at run
// time, and an arm no test exercises is how macOS spent a release reporting
// success without creating anything (#129).

func createFileCommand(goos string, filePath string) (string, error) {
	switch goos {
	case "windows":
		return "type nul > " + helper.QuoteFor(goos, filePath), nil
	case "linux", "darwin":
		return "touch " + helper.QuoteFor(goos, filePath), nil
	default:
		return "", fmt.Errorf("%w: create file on %q", errors.ErrUnsupportedOS, goos)
	}
}

func createFolderCommand(goos string, folderPath string) (string, error) {
	switch goos {
	case "windows":
		return "mkdir " + helper.QuoteFor(goos, folderPath), nil
	case "linux", "darwin":
		return "mkdir -p " + helper.QuoteFor(goos, folderPath), nil
	default:
		return "", fmt.Errorf("%w: create folder on %q", errors.ErrUnsupportedOS, goos)
	}
}

// scriptExtension is the suffix of a script command's file: a batch file on
// Windows, a shell script everywhere else.
func scriptExtension(goos string) (string, error) {
	switch goos {
	case "windows":
		return ".bat", nil
	case "linux", "darwin":
		return ".sh", nil
	default:
		return "", fmt.Errorf("%w: script command on %q", errors.ErrUnsupportedOS, goos)
	}
}
