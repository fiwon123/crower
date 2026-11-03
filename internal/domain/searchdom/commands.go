package searchdom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
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

func browserCommand(goos string, content string) (string, error) {
	switch goos {
	case "windows":
		// cmd.exe has no single-quote quoting: `start ' ' 'url'` makes the
		// first quote the window title and launches nothing. The empty ""
		// is the title; the URL is a real quoted argument, as opendom does.
		if utils.IsURL(content) {
			return fmt.Sprintf(`start "" "%s"`, content), nil
		}

		return fmt.Sprintf(`start "" "https://duckduckgo.com/?q=%s"`, content), nil

	case "linux":
		if utils.IsURL(content) {
			return fmt.Sprintf(`xdg-open "%s" `, content), nil
		}

		return fmt.Sprintf(`xdg-open "https://duckduckgo.com/?q=%s" `, content), nil

	case "darwin":
		if utils.IsURL(content) {
			return fmt.Sprintf(`open "%s"`, content), nil
		}

		return fmt.Sprintf(`open "https://duckduckgo.com/?q=%s"`, content), nil

	default:
		return "", unsupportedOSError(goos)
	}
}

func searchFileCommand(goos string, currentPath string, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf(`dir "%s" /s /b /a-d | findstr "%s"`, currentPath, content), nil
	case "linux", "darwin":
		return fmt.Sprintf(`find "%s" -type f -name "%s"`, currentPath, content), nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func searchFolderCommand(goos string, currentPath string, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf(`dir "%s" /s /b /ad | findstr "%s"`, currentPath, content), nil
	case "linux", "darwin":
		return fmt.Sprintf(`find "%s" -type d -name "%s"`, currentPath, content), nil
	default:
		return "", unsupportedOSError(goos)
	}
}

func searchFileAndFolderCommand(goos string, currentPath string, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf(`dir "%s" /s /b | findstr "%s"`, currentPath, content), nil
	case "linux", "darwin":
		return fmt.Sprintf(`find "%s" -name "%s"`, currentPath, content), nil
	default:
		return "", unsupportedOSError(goos)
	}
}
