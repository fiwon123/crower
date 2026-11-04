// Package extractcontent expands compressed archives.
//
// The operation owns the argument expansion, the per-archive destination rule
// and the recording; which tool opens which extension on which platform is the
// Archiver port's business.
package extractcontent

import (
	stderrors "errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/errors"
)

// Request is one extraction.
type Request struct {
	// Args are the archive paths, possibly glob patterns.
	Args []string

	// OutDir is where everything goes. Empty means beside each archive.
	OutDir string
}

// State is the history the operation is recorded in.
type State interface {
	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the result stream: the extractor's own output.
type Output interface {
	Result(text string)
}

// Archives expands one archive.
type Archives interface {
	// Extract expands archive with the tool for ext into outDir and returns
	// the tool's output.
	Extract(ext, archive, outDir string) (string, error)
}

// Notes renders the history note for the extraction.
type Notes interface {
	ExtractNote(out string, args []string) string
}

// UseCase runs one extraction.
type UseCase struct {
	state    State
	output   Output
	archives Archives
	notes    Notes
}

// New builds the use case from its four dependencies.
func New(state State, output Output, archives Archives, notes Notes) *UseCase {
	return &UseCase{state: state, output: output, archives: archives, notes: notes}
}

// Run expands every archive and records the operation.
func (u *UseCase) Run(request Request) error {
	paths := expand(request.Args)
	if len(paths) == 0 {
		return errors.EmptyPaths()
	}

	var failures []error

	for _, path := range paths {
		// The destination belongs to this archive: an omitted --out means
		// "next to this archive", not "next to the first archive". Assigning
		// it once outside the loop sent every later archive into the first
		// one's directory.
		outDir := request.OutDir
		if outDir == "" {
			outDir = filepath.Dir(path)
		}

		ext := extension(path)
		if ext == "" {
			failures = append(failures, fmt.Errorf("%s: no file extension to choose an extractor", path))
			continue
		}

		output, err := u.archives.Extract(ext, path, outDir)
		if err != nil {
			failures = append(failures, fmt.Errorf("%s: %w", path, err))
			continue
		}

		u.output.Result(strings.TrimRight(output, "\n"))
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	return u.state.Record(operation.Extract, "", u.notes.ExtractNote(request.OutDir, request.Args))
}

// expand turns the arguments into the paths to extract.
//
// A pattern that matches nothing is kept literally, which is how an archive the
// shell did not expand still works; a pattern that does not compile is skipped.
func expand(args []string) []string {
	paths := make([]string, 0, len(args))

	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err != nil {
			continue
		}

		if len(matches) > 0 {
			paths = append(paths, matches...)
		} else {
			paths = append(paths, arg)
		}
	}

	return paths
}

// extension is the last dot-separated part of the file name.
//
// A name with no dot returns itself, which the archiver then refuses as an
// unknown extension; a name ending in a dot returns "", which is the explicit
// "no extension" error.
func extension(path string) string {
	parts := strings.Split(filepath.Base(path), ".")

	return parts[len(parts)-1]
}
