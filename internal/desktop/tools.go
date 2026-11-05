package desktop

import "github.com/fiwon123/crower/internal/domain/timeline"

// SearchMode is what the tools view looks at.
type SearchMode string

const (
	SearchBoth    SearchMode = "both"
	SearchFiles   SearchMode = "files"
	SearchFolders SearchMode = "folders"
	SearchWeb     SearchMode = "web"
)

// Tools runs the file operations the CLI offers, for a window.
//
// Every call returns the message to show and the history the operation
// recorded: a copy, a move, a search, an extract or an open is a history
// entry, so a front end has to see it without being restarted.
//
// It is declared here, at the point of use, so this package depends on neither
// the use cases nor their adapters; the command layer supplies one.
type Tools interface {
	Copy(sources []string, destination string) (string, timeline.Log, error)
	Move(sources []string, destination string) (string, timeline.Log, error)
	Search(mode SearchMode, term, path string) (string, timeline.Log, error)
	Extract(archives []string, outDir string) (string, timeline.Log, error)
	Open(paths []string) (string, timeline.Log, error)
}
