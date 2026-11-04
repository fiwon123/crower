package system

import (
	"fmt"
	"net/url"
	"runtime"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/shell"
	"github.com/fiwon123/crower/pkg/utils"
)

// Search runs the platform's search and browser commands.
type Search struct {
	shell *shell.Runner

	// goos is a field so the dispatch can be exercised by a test.
	goos string
}

// NewSearch builds a Search that runs on the host.
func NewSearch() *Search {
	return &Search{shell: shell.New(), goos: runtime.GOOS}
}

// Browser opens a search in the default browser.
func (s *Search) Browser(query string) error {
	commandLine, err := BrowserCommand(s.goos, query)
	if err != nil {
		return err
	}

	return s.shell.Start(commandLine)
}

// Files searches files under path.
func (s *Search) Files(path, term string) (string, error) {
	commandLine, err := FilesSearchCommand(s.goos, path, term)
	if err != nil {
		return "", err
	}

	return s.shell.Exec(commandLine)
}

// Folders searches folders under path.
func (s *Search) Folders(path, term string) (string, error) {
	commandLine, err := FoldersSearchCommand(s.goos, path, term)
	if err != nil {
		return "", err
	}

	return s.shell.Exec(commandLine)
}

// Both searches files and folders under path.
func (s *Search) Both(path, term string) (string, error) {
	commandLine, err := AllSearchCommand(s.goos, path, term)
	if err != nil {
		return "", err
	}

	return s.shell.Exec(commandLine)
}

// The search commands are pure functions of goos so a table test can cover
// every platform the release ships: runtime.GOOS cannot be overridden at run
// time, and an arm no test exercises is how macOS spent a release reporting
// nothing found instead of searching (#129).
//
// Every user-supplied value goes through shell.QuoteFor: hand-written double
// quotes let a value containing `"` close the quoted region, after which `&`
// runs as a command separator on Windows.
func unsupportedSearchOS(goos string) error {
	return fmt.Errorf("%w: search %q", errors.ErrUnsupportedOS, goos)
}

func BrowserCommand(goos, content string) (string, error) {
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
		return "start \"\" " + shell.QuoteFor(goos, target), nil

	case "linux":
		return "xdg-open " + shell.QuoteFor(goos, target), nil

	case "darwin":
		return "open " + shell.QuoteFor(goos, target), nil

	default:
		return "", unsupportedSearchOS(goos)
	}
}

func FilesSearchCommand(goos, currentPath, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf("dir %s /s /b /a-d | findstr %s",
			shell.QuoteFor(goos, currentPath), shell.QuoteFor(goos, content)), nil
	case "linux", "darwin":
		return fmt.Sprintf("find %s -type f -name %s",
			shell.QuoteFor(goos, currentPath), shell.QuoteFor(goos, content)), nil
	default:
		return "", unsupportedSearchOS(goos)
	}
}

func FoldersSearchCommand(goos, currentPath, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf("dir %s /s /b /ad | findstr %s",
			shell.QuoteFor(goos, currentPath), shell.QuoteFor(goos, content)), nil
	case "linux", "darwin":
		return fmt.Sprintf("find %s -type d -name %s",
			shell.QuoteFor(goos, currentPath), shell.QuoteFor(goos, content)), nil
	default:
		return "", unsupportedSearchOS(goos)
	}
}

func AllSearchCommand(goos, currentPath, content string) (string, error) {
	switch goos {
	case "windows":
		return fmt.Sprintf("dir %s /s /b | findstr %s",
			shell.QuoteFor(goos, currentPath), shell.QuoteFor(goos, content)), nil
	case "linux", "darwin":
		return fmt.Sprintf("find %s -name %s",
			shell.QuoteFor(goos, currentPath), shell.QuoteFor(goos, content)), nil
	default:
		return "", unsupportedSearchOS(goos)
	}
}
