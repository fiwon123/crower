package cmd

import (
	"github.com/fiwon123/crower/internal/application/updatecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"

	"github.com/spf13/cobra"
)

var name string
var exec string

var create bool
var execute bool

// Cmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [command]",
	Short: "Change a command's name, aliases or exec.",
	Long: `Change a command's name, aliases or exec.

	crower update greet --exec "echo hi"
	crower update greet --name hello --alias hi

The command is found by name or by alias. The last touched command can be
updated instead:
	crower update --last
	crower update --create
	crower update --execute

With no field given, crower asks for each one.`,
	Example: `	crower update greet --exec "echo hi"
	crower update hi --name hello
	crower update --last
	crower update --create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		session := newLegacySession(cfg)

		request := updatecommand.Request{
			Changes: catalog.Changes{
				Name:    name,
				Aliases: allAlias,
				Exec:    exec,
			},
		}

		if len(args) != 0 {
			request.Key = args[0]
		}

		switch {
		case last:
			request.Last = operation.Update
		case create:
			request.Last = operation.Create
		case execute:
			request.Last = operation.Execute
		}

		return updatecommand.New(session, session, session, session).Run(request)
	},
}

func init() {
	AddNameFlag(updateCmd, &name)
	AddAllAliasFlag(updateCmd, &allAlias)
	AddExecFlag(updateCmd, &exec)

	updateCmd.Flags().BoolVarP(&last, "last", "l", false, "update recent updated command")
	updateCmd.Flags().BoolVarP(&create, "create", "c", false, "update recent created command")
	updateCmd.Flags().BoolVarP(&execute, "execute", "x", false, "update recent updated command")
}
