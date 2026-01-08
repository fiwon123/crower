package restore

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var Cmd = &cobra.Command{
	Use:   "restore",
	Short: "restore specific command",
	Long:  `restore specific command`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		operations.Restore(args, app)
	},
}

func init() {

}
