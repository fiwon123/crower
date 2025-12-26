package search

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the search command
var Cmd = &cobra.Command{
	Use:   "search",
	Short: "search on browser by default.",
	Long: `search on browser by default.

	But can be used to search in folders too`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		core.HandlePayload(
			payload.New(operation.Search, args, "", nil, ""),
			app,
		)
	},
}

func init() {

}
