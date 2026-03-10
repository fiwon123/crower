package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/listdom"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var process bool
var history bool

// Cmd represents the list command
var ListCmd = &cobra.Command{
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

		app := app.InitApp(cfgFilePath)

		core := listdom.NewCore(app)

		if process {
			core.ListProcess(args)
		} else if history {
			core.ListHistory()
		} else if folderFlag {
			core.ListFolder(args)
		} else if systemFlag {
			core.ListSystem()
		} else if sysPathFlag {
			core.ListSysPath()
		} else {
			core.ListCommands()
		}

	},
}

func init() {

	ListCmd.Flags().BoolVarP(&process, "process", "p", false, "list all process")
	ListCmd.Flags().BoolVarP(&history, "history", "i", false, "list history")
	ListCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "list folder files")
	ListCmd.Flags().BoolVar(&systemFlag, "system", false, "list all system variables")
	ListCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "list path system variable")
}
