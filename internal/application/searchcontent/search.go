// Package searchcontent searches files, folders or the web.
//
// The operation is deliberately small: the platform commands are behind the
// Searcher port, and this package owns the argument parsing, the result shape
// and the recording. One detail is kept from the legacy flow on purpose -- a
// single positional value leaves both the query empty and the path `.` -- see
// the issue that tracks it; fixing it would change documented behaviour in a
// refactor.
package searchcontent

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/domain/operation"
)

// Mode is what a request searches.
type Mode int

const (
	// ModeBrowser opens a web search.
	ModeBrowser Mode = iota

	// ModeFile searches file names.
	ModeFile

	// ModeFolder searches folder names.
	ModeFolder

	// ModeBoth searches files and folders.
	ModeBoth
)

// Request is one search.
type Request struct {
	Mode Mode
	Args []string
}

// State is the history the operation is recorded in.
type State interface {
	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the result stream and its format.
type Output interface {
	// Emit writes text in text mode and payload as JSON in structured mode.
	Emit(text string, payload any) error
}

// Searcher runs the platform's search and browser commands.
type Searcher interface {
	// Browser opens a search in the default browser.
	Browser(query string) error

	// Files searches file names under path.
	Files(path, term string) (string, error)

	// Folders searches folder names under path.
	Folders(path, term string) (string, error)

	// Both searches file and folder names under path.
	Both(path, term string) (string, error)
}

// Notes renders the history note for each mode.
type Notes interface {
	SearchBrowserNote(args []string) string
	SearchFileNote(args []string) string
	SearchFolderNote(args []string) string
	SearchBothNote(args []string) string
}

// UseCase runs one search.
type UseCase struct {
	state    State
	output   Output
	searcher Searcher
	notes    Notes
}

// New builds the use case from its four dependencies.
func New(state State, output Output, searcher Searcher, notes Notes) *UseCase {
	return &UseCase{state: state, output: output, searcher: searcher, notes: notes}
}

// Run writes the search result and records it.
func (u *UseCase) Run(request Request) error {
	switch request.Mode {
	case ModeBrowser:
		return u.searchBrowser(request.Args)

	case ModeFile:
		return u.searchFiles(request.Args)

	case ModeFolder:
		return u.searchFolders(request.Args)

	case ModeBoth:
		return u.searchBoth(request.Args)

	default:
		return fmt.Errorf("unknown search mode %d", request.Mode)
	}
}

// searchBrowser opens a web search.
func (u *UseCase) searchBrowser(args []string) error {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	if err := u.searcher.Browser(content); err != nil {
		return err
	}

	return u.emit("", searchPayload(operation.Browser, content, "", ""), u.notes.SearchBrowserNote(args))
}

// searchFiles searches file names.
func (u *UseCase) searchFiles(args []string) error {
	content, currentPath := assertArgs(args)

	out, err := u.searcher.Files(currentPath, content)
	if err != nil {
		return err
	}

	return u.emit(out, searchPayload(operation.File, content, currentPath, out), u.notes.SearchFileNote(args))
}

// searchFolders searches folder names.
func (u *UseCase) searchFolders(args []string) error {
	content, currentPath := assertArgs(args)

	out, err := u.searcher.Folders(currentPath, content)
	if err != nil {
		return err
	}

	return u.emit(out, searchPayload(operation.Folder, content, currentPath, out), u.notes.SearchFolderNote(args))
}

// searchBoth searches file and folder names.
func (u *UseCase) searchBoth(args []string) error {
	content, currentPath := assertArgs(args)

	out, err := u.searcher.Both(currentPath, content)
	if err != nil {
		return err
	}

	return u.emit(out, searchPayload(operation.FileAndFolder, content, currentPath, out), u.notes.SearchBothNote(args))
}

// assertArgs splits the positional values the way the search has always done:
// the query and the path are read only when both are present.
func assertArgs(args []string) (string, string) {
	content := ""
	currentPath := "."

	if len(args) > 1 {
		content = args[0]
		currentPath = args[1]
	}

	return content, currentPath
}

// emit writes the result in the configured format and records the operation.
func (u *UseCase) emit(text string, payload any, note string) error {
	if err := u.output.Emit(text, payload); err != nil {
		return err
	}

	return u.state.Record(operation.Search, "", note)
}

// SearchResult is the JSON result of one search.
//
// Kind is the kind of thing searched for -- file, folder, file_and_folder or
// browser -- and Results is what was found, one entry per line of the shell
// output. Path is empty for a browser search, which returns nothing on stdout.
type SearchResult struct {
	Kind    string   `json:"kind"`
	Query   string   `json:"query"`
	Path    string   `json:"path"`
	Results []string `json:"results"`
}

// Result is the object `search --json` emits.
type Result struct {
	Search SearchResult `json:"search"`
}

func searchPayload(kind operation.SubOperation, query, path, out string) Result {
	return Result{Search: SearchResult{
		Kind:    string(kind),
		Query:   query,
		Path:    path,
		Results: splitResults(out),
	}}
}

// splitResults turns the shell's line-oriented output into entries.
//
// A path containing a newline cannot survive this, which is why the JSON
// contract documents results as lines rather than paths. Empty lines are
// dropped: a path is never empty, so an empty line is formatting, not a result.
func splitResults(out string) []string {
	results := make([]string, 0)

	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}

		results = append(results, line)
	}

	return results
}
