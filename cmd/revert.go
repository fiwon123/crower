package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/revertdom"

	"github.com/spf13/cobra"
)

// Cmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert <steps>",
	Short: "Undo operations, newest first.",
	Long: `Undo operations, newest first.

The argument is how many history entries to roll back. "crower list
--history" shows what is there.

	crower revert 1
	crower revert 3`,
	Example: `	crower revert 1
	crower revert 3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "the number of steps to revert")
		}

		return revertdom.NewCore(cfg).Revert(args)
	},
}
