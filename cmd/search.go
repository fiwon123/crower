package cmd

import (
	"github.com/fiwon123/crower/internal/domain/searchdom"

	"github.com/spf13/cobra"
)

var browserFlag bool

// Cmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <term> [path]",
	Short: "search files, folder and on browser",
	Long: `search files, folder and on browser

search files:
	- Using 'crower search "PART_NAME" ./'
	- Using 'crower search --file "PART_NAME" ./'

search folder:
	- Using 'crower search "PART_NAME" ./'
	- Using 'crower search --folder "PART_NAME" ./'

search on browser:
	- Using 'crower search --browser "CONTENT"'

	`,
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
}
