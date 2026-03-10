package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/domain/listdom"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var last bool

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute",
	Short: "execute command",
	Long:  `execute command`,
	Run: func(cmd *cobra.Command, args []string) {

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		listHandler := listdom.NewHandler(app)
		core := executedom.NewCore(app, listHandler)

		if last {
			core.ExecuteLast(operationsdata.Execute, args)
		} else if createFlag {
			core.ExecuteLast(operationsdata.Create, args)
		} else if updateFlag {
			core.ExecuteLast(operationsdata.Update, args)
		} else {
			core.ExecuteCommand(args)
		}
	},
}

func init() {
	executeCmd.Flags().BoolVarP(&last, "last", "l", false, "execute recent executed command")
	executeCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "execute recent created command")
	executeCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "execute recent updated command")
}
