package revert

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
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

		core.HandlePayload(
			payload.New(operation.Revert, args, "", nil, ""),
			app,
		)
	},
}

func init() {

}
