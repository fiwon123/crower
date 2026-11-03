package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/restoredom"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <command>",
	Short: "Restore a command from the history.",
	Long: `Restore a command from the history.

The command is looked up by name or alias in the snapshot of a past
operation, so a deleted or overwritten command can come back.

	crower restore greet
	crower restore hi`,
	Example: `	crower restore greet
	crower restore hi`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return restoredom.NewCore(cfg).Restore(args)
	},
}
