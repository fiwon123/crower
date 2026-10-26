package cmd

import (
	"github.com/fiwon123/crower/internal/domain/listdom"

	"github.com/spf13/cobra"
)

var historyFlag bool

// Cmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "list all commands, history, process, folder, system, path system",
	Long: `list all commands, history, process, folder, system, path system

list all commands by default

Example:
	crower list
	crower list --history
	crower list --process
	crower list --folder "FOLDER_PATH"
	crower list --system
	crower list --syspath
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		core := listdom.NewCore(cfg)

		switch {
		case processFlag:
			return core.ListProcess(args)
		case historyFlag:
			return core.ListHistory()
		case folderFlag:
			return core.ListFolder(args)
		case systemFlag:
			return core.ListSystem()
		case sysPathFlag:
			return core.ListSysPath()
		}

		return core.ListCommands()
	},
}

func init() {
	listCmd.Flags().BoolVarP(&processFlag, "process", "p", false, "list all process")
	listCmd.Flags().BoolVarP(&historyFlag, "history", "i", false, "list history")
	listCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "list folder files")
	listCmd.Flags().BoolVar(&systemFlag, "system", false, "list all system variables")
	listCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "list path system variable")
}
