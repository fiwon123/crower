package update

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/state"

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var name string
var allAlias []string
var exec string

var last bool
var create bool
var execute bool

// Cmd represents the update command
var Cmd = &cobra.Command{
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
			operations.UpdateLast(state.Update, name, allAlias, exec, app)
		} else if create {
			operations.UpdateLast(state.Create, name, allAlias, exec, app)
		} else if execute {
			operations.UpdateLast(state.Execute, name, allAlias, exec, app)
		} else if len(args) > 0 {
			key := ""
			if len(args) != 0 {
				key = args[0]
			}

			operations.Update(key, name, allAlias, exec, app)
		} else {
			crerrors.PrintCmdHelp("update")
		}

	},
}

func init() {
	cmdsHelper.AddNameFlag(Cmd, &name)
	cmdsHelper.AddAllAliasFlag(Cmd, &allAlias)
	cmdsHelper.AddExecFlag(Cmd, &exec)

	Cmd.Flags().BoolVarP(&last, "last", "l", false, "update recent updated command")
	Cmd.Flags().BoolVarP(&create, "create", "c", false, "update recent created command")
	Cmd.Flags().BoolVarP(&execute, "execute", "x", false, "update recent updated command")
}
