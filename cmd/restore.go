package cmd

import (
	"github.com/fiwon123/crower/internal/domain/restoredom"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <command>",
	Short: "restore specific command",
	Long:  `restore specific command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return restoredom.NewCore(cfg).Restore(args)
	},
}
