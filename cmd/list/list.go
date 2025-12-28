package list

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var process bool
var history bool
var folderFlag bool

// Cmd represents the list command
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "list all commands by default",
	Long:  `list all commands by default`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		op := operation.List
		if process {
			op = operation.ListProcess
		} else if history {
			op = operation.ListHistory
		} else if folderFlag {
			op = operation.ListFolder
		}

		core.HandlePayload(
			payload.New(op, args, "", nil, ""),
			app,
		)
	},
}

func init() {

	Cmd.Flags().BoolVarP(&process, "process", "p", false, "list all process")
	Cmd.Flags().BoolVarP(&history, "history", "i", false, "list history")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "list folder files")
}
