package cmd

import (
	"github.com/fiwon123/crower/internal/application/reverthistory"
	"github.com/fiwon123/crower/internal/legacy/listdom"

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

		// The usage text is part of the answer: a revert with no step count is
		// a usage error, and the fix belongs next to the message.
		if len(args) == 0 {
			return missingArgs(cmd, "the number of steps to revert")
		}

		session := newLegacySession(cfg, listdom.NewHandler(cfg).CommandsTable)

		return reverthistory.New(session, session, session, session).Run(args)
	},
}
