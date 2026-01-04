package search

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var folderFlag bool
var fileFlag bool

// Cmd represents the search command
var Cmd = &cobra.Command{
	Use:   "search",
	Short: "search on browser by default.",
	Long: `search on browser by default.

But can be used to search in folders too using either file flag or folder flag`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if folderFlag && fileFlag {
			operations.SearchFileAndFolder(args, app)
		} else if folderFlag {
			operations.SearchFolder(args, app)
		} else if fileFlag {
			operations.SearchFile(args, app)
		} else {
			operations.SearchBrowser(args, app)
		}
	},
}

func init() {
	Cmd.Flags().BoolVarP(&fileFlag, "file", "i", false, "search for file name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "f", false, "search for folder name")
}
