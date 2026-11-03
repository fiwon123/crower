package movedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
)

// moveCommand builds the platform's move invocation.
//
// goos is a parameter, not a read of runtime.GOOS, so a table test can cover
// every platform the release ships (#129).
func moveCommand(goos string, source string, destFolder string) (string, error) {
	switch goos {
	case "windows":
		return "move " + helper.QuoteFor(goos, source) + " " + helper.QuoteFor(goos, destFolder), nil
	case "linux", "darwin":
		return "mv " + helper.QuoteFor(goos, source) + " " + helper.QuoteFor(goos, destFolder), nil
	default:
		return "", fmt.Errorf("%w: move on %q", errors.ErrUnsupportedOS, goos)
	}
}
