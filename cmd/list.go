package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

var history bool

// Cmd represents the list command
var listCmd = &cobra.Command{
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
		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := listdom.NewCore(app)

		if processFlag {
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

	listCmd.Flags().BoolVarP(&processFlag, "process", "p", false, "list all process")
	listCmd.Flags().BoolVarP(&history, "history", "i", false, "list history")
	listCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "list folder files")
	listCmd.Flags().BoolVar(&systemFlag, "system", false, "list all system variables")
	listCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "list path system variable")
}
