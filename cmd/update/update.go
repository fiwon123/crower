package update

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
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
	Long:  `update commands`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		op := operation.Update
		if last {
			op = operation.UpdateLast
		} else if create {
			op = operation.UpdateCreate
		} else if execute {
			op = operation.UpdateExecute
		}

		core.HandlePayload(
			payload.New(op, args, name, allAlias, exec),
			app,
		)
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
