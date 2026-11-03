package searchdom

import (
	"strings"

	"github.com/fiwon123/crower/internal/domain/operation"
)

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

func searchPayload(kind operation.SubOperation, query string, path string, out string) Result {
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
