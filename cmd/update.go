package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/domain/updatedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

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
		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		listHandler := listdom.NewHandler(app)
		core := updatedom.NewCore(app, listHandler)

		if last {
			core.UpdateLast(operations.Update, name, allAlias, exec)
		} else if create {
			core.UpdateLast(operations.Create, name, allAlias, exec)
		} else if execute {
			core.UpdateLast(operations.Execute, name, allAlias, exec)
		} else if len(args) > 0 {
			core.UpdateCommand(args, name, allAlias, exec)
		} else {
			errors.PrintCmdHelp("update", app)
		}

	},
}

func init() {
	helper.AddNameFlag(updateCmd, &name)
	helper.AddAllAliasFlag(updateCmd, &allAlias)
	helper.AddExecFlag(updateCmd, &exec)

	updateCmd.Flags().BoolVarP(&last, "last", "l", false, "update recent updated command")
	updateCmd.Flags().BoolVarP(&create, "create", "c", false, "update recent created command")
	updateCmd.Flags().BoolVarP(&execute, "execute", "x", false, "update recent updated command")
}
