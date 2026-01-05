package search

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/crerrors"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var folderFlag bool
var fileFlag bool
var browserFlag bool

// Cmd represents the search command
var Cmd = &cobra.Command{
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
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if browserFlag {
			operations.SearchBrowser(args, app)
		} else if fileFlag {
			operations.SearchFile(args, app)
		} else if folderFlag {
			operations.SearchFolder(args, app)
		} else if len(args) > 0 {
			operations.SearchFileAndFolder(args, app)
		} else {
			crerrors.PrintCmdHelp("search")
		}

	},
}

func init() {
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "search for file name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "search for folder name")
	Cmd.Flags().BoolVarP(&browserFlag, "browser", "b", false, "search on default browser")
}
