package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/movedom"

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
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "source paths and a destination")
		}

		return movedom.NewCore(cfg).Move(args)
	},
}
