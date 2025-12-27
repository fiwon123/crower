package extract

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the extract command
var Cmd = &cobra.Command{
	Use:   "extract",
	Short: "extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...)",
	Long:  `extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...)`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		core.HandlePayload(
			payload.New(operation.Extract, args, "", nil, ""),
			app,
		)
	},
}

func init() {
}
