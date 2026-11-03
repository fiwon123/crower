package searchdom

import (
	"fmt"
	"net/url"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/fiwon123/crower/pkg/utils"
)

// unsupportedOSError reports a search that has no implementation for the
// running platform. It is an error rather than an empty result: a search that
// matched nothing and a search that cannot run must not look the same.
func unsupportedOSError(goos string) error {
	return fmt.Errorf("%w: search %q", errors.ErrUnsupportedOS, goos)
}

// The search commands are pure functions of goos so a table test can cover
// every platform the release ships: runtime.GOOS cannot be overridden at run
// time, and an arm no test exercises is how macOS spent a release reporting
// nothing found instead of searching (#129).
//
// Every user-supplied value goes through helper.QuoteFor: hand-written double
// quotes let a value containing `"` close the quoted region, after which `&`
// runs as a command separator on Windows.

func browserCommand(goos string, content string) (string, error) {
	target := content
	if !utils.IsURL(content) {
		// A search term is not a URL fragment: encode it, or `&` and `#`
		// change the query instead of being searched for.
		target = "https://duckduckgo.com/?q=" + url.QueryEscape(content)
	}

	switch goos {
	case "windows":
		// cmd.exe has no single-quote quoting, so the empty "" is the window
		// title and the target is a real (quoted when needed) argument.
		return "start \"\" " + helper.QuoteFor(goos, target), nil

	case "linux":
		return "xdg-open " + helper.QuoteFor(goos, target), nil

	case "darwin":
		return "open " + helper.QuoteFor(goos, target), nil

	default:
		return "", unsupportedOSError(goos)
	}
}

func searchFileCommand(goos string, currentPath string, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf("dir %s /s /b /a-d | findstr %s",
			helper.QuoteFor(goos, currentPath), helper.QuoteFor(goos, content)), nil
	case "linux", "darwin":
		return fmt.Sprintf("find %s -type f -name %s",
			helper.QuoteFor(goos, currentPath), helper.QuoteFor(goos, content)), nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func searchFolderCommand(goos string, currentPath string, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf("dir %s /s /b /ad | findstr %s",
			helper.QuoteFor(goos, currentPath), helper.QuoteFor(goos, content)), nil
	case "linux", "darwin":
		return fmt.Sprintf("find %s -type d -name %s",
			helper.QuoteFor(goos, currentPath), helper.QuoteFor(goos, content)), nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func searchFileAndFolderCommand(goos string, currentPath string, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf("dir %s /s /b | findstr %s",
			helper.QuoteFor(goos, currentPath), helper.QuoteFor(goos, content)), nil
	case "linux", "darwin":
		return fmt.Sprintf("find %s -name %s",
			helper.QuoteFor(goos, currentPath), helper.QuoteFor(goos, content)), nil
	default:
		return "", unsupportedOSError(goos)
	}
}
