package reset

import (
	"github.com/fiwon123/crower/internal/core"
	resetoperations "github.com/fiwon123/crower/internal/core/operations/reset"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the reset command
var Cmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all commands",
	Long:  `reset all commands`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		resetoperations.Reset(app)
	},
}

func init() {

}
