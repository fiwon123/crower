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
var sysPathFlag bool
var systemFlag bool

// Cmd represents the list command
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "list all commands, history, process, folder, system, path system",
	Long: `list all commands, history, process, folder, system, path system

list all commands by default

Example:
	crower --list
	crower --list --history
	crower --list --process
	crower --list --folder "FOLDER_PATH"
	crower --list --system
	crower --list --syspath
`,
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
		} else if systemFlag {
			op = operation.ListSystem
		} else if sysPathFlag {
			op = operation.ListSysPath
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
	Cmd.Flags().BoolVarP(&systemFlag, "system", "a", false, "list all system variables")
	Cmd.Flags().BoolVarP(&sysPathFlag, "syspath", "s", false, "list path system variable")
}
