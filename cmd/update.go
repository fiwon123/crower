package cmd

import (
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/domain/updatedom"

	"github.com/spf13/cobra"
)

var name string
var exec string

var create bool
var execute bool

// Cmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [command]",
	Short: "update commands",
	Long: `update commands

update command:
	- Use interactive input just typing 'crower update' without arguments
	- Use argument key (command name or command alias) 'crower update "COMMAND_KEY"'
	- using flags --last (update), --create or --execute to update last operation flag

Examples:
	crower update
	crower update com_name
	crower update com_alias
	crower update --last
	crower update --create
	crower update --execute
	crower update com_name --name "test" --exec "echo t"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		listHandler := listdom.NewHandler(cfg)
		core := updatedom.NewCore(cfg, listHandler)

		switch {
		case last:
			return core.UpdateLast(operations.Update, name, allAlias, exec)
		case create:
			return core.UpdateLast(operations.Create, name, allAlias, exec)
		case execute:
			return core.UpdateLast(operations.Execute, name, allAlias, exec)
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
