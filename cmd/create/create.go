package create

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
var process string

// Cmd represents the create command
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a command",
	Long: `create command can be used doing:

crower create com1 "'echo com1'"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if process != "" {
			name = process
		}

		core.HandlePayload(
			payload.New(operation.Create, args, name, allAlias, exec),
			app,
		)
	},
}

func init() {
	cmdsHelper.AddNameFlag(Cmd, &name)
	cmdsHelper.AddAllAliasFlag(Cmd, &allAlias)
	cmdsHelper.AddExecFlag(Cmd, &exec)

	Cmd.Flags().StringVarP(&process, "process", "p", "", "process name or pid")
}
