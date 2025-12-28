package extract

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var outDirFlag string

// Cmd represents the extract command
var Cmd = &cobra.Command{
	Use:   "extract",
	Short: "extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...)",
	Long:  `extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...)`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		operations.Extract(args, outDirFlag, app)
	},
}

func init() {
	Cmd.Flags().StringVarP(&outDirFlag, "out", "o", "", "out folder where to be extracted")
}
