package cmd

import (
	"github.com/fiwon123/crower/internal/application/transfercontent"
	"github.com/fiwon123/crower/internal/domain/operation"

	"github.com/spf13/cobra"
)

// Cmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move <source...> <destination>",
	Short: "Move files or folders to a destination.",
	Long: `Move files or folders to a destination.

	crower move notes.txt ~/archive
	crower move a.txt b.txt ~/archive
	crower move ./reports ~/archive`,
	Example: `	crower move notes.txt ~/archive
	crower move a.txt b.txt ~/archive
	crower move ./reports ~/archive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "source paths and a destination")
		}

		session := newTransferSession(cfg)

		return transfercontent.New(session, session, session, session).Run(
			transfercontent.Request{Op: operation.Move, Args: args},
		)
	},
}
