package search

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
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

	But can be used to search in folders too`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		op := operation.Search
		if folderFlag && fileFlag {
			op = operation.SearchFileAndFolder
		} else if folderFlag {
			op = operation.SearchFolder
		} else if fileFlag {
			op = operation.SearchFile
		}

		core.HandlePayload(
			payload.New(op, args, "", nil, ""),
			app,
		)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&fileFlag, "file", "i", false, "search for file name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "f", false, "search for folder name")
}
