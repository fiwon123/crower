package revertcmd

import (
	"github.com/fiwon123/crower/internal/core"
	revertoperations "github.com/fiwon123/crower/internal/core/operations/revert"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the revert command
var Cmd = &cobra.Command{
	Use:   "revert",
	Short: "revert history by one",
	Long:  `revert history by one"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		revertoperations.Revert(args, app)
	},
}

func init() {
}
