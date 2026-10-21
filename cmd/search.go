package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/searchdom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

var browserFlag bool

// Cmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
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
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := searchdom.NewCore(app)

		if browserFlag {
			core.SearchBrowser(args)
		} else if fileFlag {
			core.SearchFile(args)
		} else if folderFlag {
			core.SearchFolder(args)
		} else if len(args) > 0 {
			core.SearchFileAndFolder(args)
		} else {
			errors.PrintCmdHelp("search", app)
		}

	},
}

func init() {
	searchCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "search for file name")
	searchCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "search for folder name")
	searchCmd.Flags().BoolVarP(&browserFlag, "browser", "b", false, "search on default browser")
}
