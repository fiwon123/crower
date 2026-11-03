package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/listdom"

	"github.com/spf13/cobra"
)

var historyFlag bool

// Cmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List commands, history, processes, folders or variables.",
	Long: `List commands, history, processes, folders or variables.

With no flag, crower lists the managed commands in the order they were
saved.

	crower list
	crower list --history
	crower list --process
	crower list --folder ~/projects
	crower list --system
	crower list --syspath

--json prints one JSON object on stdout instead of a table.`,
	Example: `	crower list
	crower list --history
	crower list --folder ~/projects
	crower list --json`,
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
	listCmd.Flags().Bool(flagJSON, false, "print the result as one JSON object")
}
