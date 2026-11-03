package cmd

import (
	"github.com/fiwon123/crower/internal/domain/revertdom"

	"github.com/spf13/cobra"
)

// Cmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert <steps>",
	Short: "revert history by one",
	Long:  `revert history by one`,
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
