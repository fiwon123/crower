package deletedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
)

// removeCommand builds the platform's delete invocation.
//
// goos is a parameter, not a read of runtime.GOOS, so a table test can cover
// every platform the release ships (#129).
func removeCommand(goos string, path string, folder bool) (string, error) {
	switch goos {
	case "windows":
		if folder {
			return "rmdir /s /q " + helper.QuoteFor(goos, path), nil
		}

		return "del " + helper.QuoteFor(goos, path), nil

	case "linux", "darwin":
		commandLine := "rm "
		if folder {
			commandLine = "rm -r "
		}

		return commandLine + helper.QuoteFor(goos, path), nil

	default:
		return "", fmt.Errorf("%w: delete on %q", errors.ErrUnsupportedOS, goos)
	}
}
