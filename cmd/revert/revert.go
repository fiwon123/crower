package revert

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
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

		operations.Revert(app)
	},
}

func init() {

}
