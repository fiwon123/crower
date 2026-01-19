package restore

import (
	"github.com/fiwon123/crower/internal/core"
	restoreoperations "github.com/fiwon123/crower/internal/core/operations/restore"
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

		restoreoperations.Restore(args, app)
	},
}

func init() {

}
