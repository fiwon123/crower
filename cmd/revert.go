package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	revertdom "github.com/fiwon123/crower/internal/domain/revert"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "revert history by one",
	Long:  `revert history by one"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := revertdom.NewCore(app)

		core.Revert(args)
	},
}

func init() {
}
