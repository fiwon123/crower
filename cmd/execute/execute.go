package execute

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/data/state"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var last bool
var createFlag bool
var updateFlag bool

// executeCmd represents the execute command
var Cmd = &cobra.Command{
	Use:   "execute",
	Short: "execute command",
	Long:  `execute command`,
	Run: func(cmd *cobra.Command, args []string) {

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if last {
			operations.ExecuteLast(state.Execute, args, app)
		} else if createFlag {
			operations.ExecuteLast(state.Create, args, app)
		} else if updateFlag {
			operations.ExecuteLast(state.Update, args, app)
		} else {
			operations.Execute(args, app)
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// executeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// executeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	Cmd.Flags().BoolVarP(&last, "last", "l", false, "execute recent executed command")
	Cmd.Flags().BoolVarP(&createFlag, "create", "c", false, "execute recent created command")
	Cmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "execute recent updated command")
}
