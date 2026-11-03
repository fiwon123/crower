package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/searchdom"

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

		core := searchdom.NewCore(cfg)

		switch {
		case browserFlag:
			return core.SearchBrowser(args)
		case fileFlag:
			return core.SearchFile(args)
		case folderFlag:
			return core.SearchFolder(args)
		case len(args) > 0:
			return core.SearchFileAndFolder(args)
		}

		return missingArgs(cmd, "a search term")
	},
}

func init() {
	searchCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "search for file name")
	searchCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "search for folder name")
	searchCmd.Flags().BoolVarP(&browserFlag, "browser", "b", false, "search on default browser")
	searchCmd.Flags().Bool(flagJSON, false, "print the result as one JSON object")
}
