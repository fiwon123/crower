package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/domain/resetdom"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all commands",
	Long:  `reset all commands`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		listHandler := listdom.NewHandler(app)
		core := resetdom.NewCore(app, listHandler)

		core.Reset()
	},
}

func init() {

}
