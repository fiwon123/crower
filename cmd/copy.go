package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/copydom"

	"github.com/spf13/cobra"
)

// Cmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy <source...> <destination>",
	Short: "Copy files or folders to a destination.",
	Long: `Copy files or folders to a destination.

	crower copy notes.txt ~/backup
	crower copy a.txt b.txt ~/backup
	crower copy ./reports ~/backup`,
	Example: `	crower copy notes.txt ~/backup
	crower copy a.txt b.txt ~/backup
	crower copy ./reports ~/backup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "source paths and a destination")
		}

		return copydom.NewCore(cfg).Copy(args)
	},
}
