package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/revertdom"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

// Cmd represents the revert command
var revertCmd = &cobra.Command{
	Use:   "revert",
	Short: "revert history by one",
	Long:  `revert history by one"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := revertdom.NewCore(app)

		core.Revert(args)
	},
}

func init() {
}
