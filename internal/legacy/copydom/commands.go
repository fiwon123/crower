package copydom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
)

// copyCommand builds the platform's copy invocation.
//
// It is a pure function of goos so a table test can cover every platform the
// release ships: runtime.GOOS cannot be overridden at run time, and an arm no
// test exercises is how macOS spent a release reporting success without
// copying anything (#129).
func copyCommand(goos string, source string, destFolder string, folder bool) (string, error) {
	switch goos {
	case "windows":
		if folder {
			return "xcopy " + helper.QuoteFor(goos, source) + " " + helper.QuoteFor(goos, destFolder) + " /E /I", nil
		}

		return "copy " + helper.QuoteFor(goos, source) + " " + helper.QuoteFor(goos, destFolder), nil

	case "linux", "darwin":
		commandLine := "cp "
		if folder {
			commandLine += "-r "
		}

		return commandLine + helper.QuoteFor(goos, source) + " " + helper.QuoteFor(goos, destFolder), nil

	default:
		return "", fmt.Errorf("%w: copy on %q", errors.ErrUnsupportedOS, goos)
	}
}
