package cmd

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/legacy/listdom"
	"github.com/fiwon123/crower/internal/legacy/updatedom"

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

		listHandler := listdom.NewHandler(cfg)
		core := updatedom.NewCore(cfg, listHandler)

		switch {
		case last:
			return core.UpdateLast(operation.Update, name, allAlias, exec)
		case create:
			return core.UpdateLast(operation.Create, name, allAlias, exec)
		case execute:
			return core.UpdateLast(operation.Execute, name, allAlias, exec)
		case len(args) > 0:
			return core.UpdateCommand(args, name, allAlias, exec)
		}

		return missingArgs(cmd, "a command name or alias")
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
