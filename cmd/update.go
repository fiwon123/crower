package cmd

import (
	"github.com/fiwon123/crower/internal/core"
	updateoperations "github.com/fiwon123/crower/internal/core/operations/update"
	"github.com/fiwon123/crower/internal/crowererrors"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var name string
var exec string

var create bool
var execute bool

// Cmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
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
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if last {
			updateoperations.UpdateLast(operationsdata.Update, name, allAlias, exec, app)
		} else if create {
			updateoperations.UpdateLast(operationsdata.Create, name, allAlias, exec, app)
		} else if execute {
			updateoperations.UpdateLast(operationsdata.Execute, name, allAlias, exec, app)
		} else if len(args) > 0 {
			updateoperations.UpdateCommand(args, name, allAlias, exec, app)
		} else {
			crowererrors.PrintCmdHelp("update", app)
		}

	},
}

func init() {
	cmdsHelper.AddNameFlag(updateCmd, &name)
	cmdsHelper.AddAllAliasFlag(updateCmd, &allAlias)
	cmdsHelper.AddExecFlag(updateCmd, &exec)

	updateCmd.Flags().BoolVarP(&last, "last", "l", false, "update recent updated command")
	updateCmd.Flags().BoolVarP(&create, "create", "c", false, "update recent created command")
	updateCmd.Flags().BoolVarP(&execute, "execute", "x", false, "update recent updated command")
}
