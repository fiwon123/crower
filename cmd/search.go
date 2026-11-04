package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/searchcontent"

	"github.com/spf13/cobra"
)

var browserFlag bool

// Cmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <term> [path]",
	Short: "Search files, folders, or the web.",
	Long: `Search files, folders, or the web.

	crower search report ./notes
	crower search --file report ./notes
	crower search --folder projects ~
	crower search --browser "golang context"

--json prints one JSON object on stdout instead of a table.`,
	Example: `	crower search report ./notes
	crower search --folder projects ~
	crower search --browser "golang context"
	crower search report ./notes --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		request := searchcontent.Request{Args: args}

		// The flags decide what is searched; the use case owns how.
		switch {
		case browserFlag:
			request.Mode = searchcontent.ModeBrowser
		case fileFlag:
			request.Mode = searchcontent.ModeFile
		case folderFlag:
			request.Mode = searchcontent.ModeFolder
		case len(args) > 0:
			request.Mode = searchcontent.ModeBoth
		default:
			return missingArgs(cmd, "a search term")
		}

		session := newSearchSession(cfg)

		return searchcontent.New(session, session, session, session).Run(request)
	},
}

// searchSession adapts the configuration to the ports the search use case
// declares.
type searchSession struct {
	legacySession

	searcher *system.Search
}

// newSearchSession wires the configuration and the platform search.
func newSearchSession(cfg *app.Config) searchSession {
	return searchSession{
		legacySession: newLegacySession(cfg),
		searcher:      system.NewSearch(),
	}
}

// Emit writes the text, or the payload when --json was asked for.
func (s searchSession) Emit(text string, payload any) error {
	return s.cfg.Result(text, payload)
}

// Browser opens a search in the default browser.
func (s searchSession) Browser(query string) error {
	return s.searcher.Browser(query)
}

// Files searches file names under path.
func (s searchSession) Files(path, term string) (string, error) {
	return s.searcher.Files(path, term)
}

// Folders searches folder names under path.
func (s searchSession) Folders(path, term string) (string, error) {
	return s.searcher.Folders(path, term)
}

// Both searches file and folder names under path.
func (s searchSession) Both(path, term string) (string, error) {
	return s.searcher.Both(path, term)
}

// SearchBrowserNote is the history note for `search --browser`.
func (s searchSession) SearchBrowserNote(args []string) string {
	return notes.SearchBrowserNote(args)
}

// SearchFileNote is the history note for `search --file`.
func (s searchSession) SearchFileNote(args []string) string {
	return notes.SearchFileNote(args)
}

// SearchFolderNote is the history note for `search --folder`.
func (s searchSession) SearchFolderNote(args []string) string {
	return notes.SearchFolderNote(args)
}

// SearchBothNote is the history note for the positional search.
func (s searchSession) SearchBothNote(args []string) string {
	return notes.SearchBothNote(args)
}

func init() {
	searchCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "search for file name")
	searchCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "search for folder name")
	searchCmd.Flags().BoolVarP(&browserFlag, "browser", "b", false, "search on default browser")
	searchCmd.Flags().Bool(flagJSON, false, "print the result as one JSON object")
}
