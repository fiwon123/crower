package cmd

import (
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/domain/resetdom"

	"github.com/spf13/cobra"
)

// Cmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all commands",
	Long:  `reset all commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		listHandler := listdom.NewHandler(cfg)

		return resetdom.NewCore(cfg, listHandler).Reset()
	},
}
