package desktop

import (
	"context"

	"github.com/fiwon123/crower/internal/domain/timeline"
)

// SearchMode is what the tools view looks at.
type SearchMode string

const (
	SearchBoth    SearchMode = "both"
	SearchFiles   SearchMode = "files"
	SearchFolders SearchMode = "folders"
	SearchWeb     SearchMode = "web"
)

// Progress is how far a long run has got: the item that just started and the
// count around it.
type Progress struct {
	Done    int
	Total   int
	Current string
}

// Tools runs the file operations the CLI offers, for a window.
//
// Every call returns the message to show and the history the operation
// recorded: a copy, a move, a search, an extract or an open is a history
// entry, so a front end has to see it without being restarted.
//
// It is declared here, at the point of use, so this package depends on neither
// the use cases nor their adapters; the command layer supplies one.
//
// Copy, Move and Extract take a context and a progress callback: they can run
// for a long time, and a window has to be able to say how far they are and stop
// them between items.
type Tools interface {
	Copy(ctx context.Context, report func(Progress), sources []string, destination string) (string, timeline.Log, error)
	Move(ctx context.Context, report func(Progress), sources []string, destination string) (string, timeline.Log, error)
	Search(mode SearchMode, term, path string) (string, timeline.Log, error)
	Extract(ctx context.Context, report func(Progress), archives []string, outDir string) (string, timeline.Log, error)
	Open(paths []string) (string, timeline.Log, error)
}
