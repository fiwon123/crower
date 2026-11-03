package opendom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
)

// openCommand builds the platform's "open with the default application"
// invocation.
//
// goos is a parameter, not a read of runtime.GOOS, so a table test can cover
// every platform the release ships (#129).
func openCommand(goos string, path string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf(`start "" %s`, helper.QuoteFor(goos, path)), nil
	case "linux":
		return "xdg-open " + helper.QuoteFor(goos, path), nil
	case "darwin":
		return "open " + helper.QuoteFor(goos, path), nil
	default:
		return "", fmt.Errorf("%w: open on %q", errors.ErrUnsupportedOS, goos)
	}
}
