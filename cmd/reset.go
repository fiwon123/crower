package cmd

import (
	"github.com/fiwon123/crower/internal/application/resetcommands"
	"github.com/fiwon123/crower/internal/legacy/listdom"

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

		table := listdom.NewHandler(cfg).CommandsTable
		session := newLegacySession(cfg, table)

		return resetcommands.New(session, session, session).Run()
	},
}
