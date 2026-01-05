package copy

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the copy command
var Cmd = &cobra.Command{
	Use:   "copy",
	Short: "copy file or folder to other place",
	Long: `copy file or folder to other place

Examples:
	crower copy "C:\Users\Test\Desktop\Test\file.txt" "C:\Users\Test\Desktop\Test\Out"
	crower copy "C:\Users\Test\Desktop\Test\Folder\" "C:\Users\Test\Desktop\Test\Out"
`,
	Run: func(cmd *cobra.Command, args []string) {

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		operations.Copy(args, app)

	},
}

func init() {

}
