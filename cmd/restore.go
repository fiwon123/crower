package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/restoredom"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore specific command",
	Long:  `restore specific command`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := restoredom.NewCore(app)

		core.Restore(args)
	},
}

func init() {

}
