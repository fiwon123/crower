package cmd

import (
	"github.com/fiwon123/crower/internal/application/resetcommands"

	"github.com/spf13/cobra"
)

// Cmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete every managed command.",
	Long: `Delete every managed command.

The command table is shown first, and nothing happens unless you confirm.
History snapshots are left alone, so "crower list --history" still works
and "crower revert" can bring the commands back.

	crower reset`,
	Example: `	crower reset`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		session := newSessionPorts(cfg)

		return resetcommands.New(session, session, session, session).Run()
	},
}
