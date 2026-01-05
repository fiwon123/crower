package open

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/crerrors"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var folderFlag bool
var systemFlag bool

// Cmd represents the open command
var Cmd = &cobra.Command{
	Use:   "open",
	Short: "open current configuration file path, folder and system",
	Long: `open current configuration file path, folder and system

Examples:
	crower open
	crower open --folder "FOLDER_PATH"
	crower open --system
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if folderFlag {
			operations.OpenFolder(args, app)
		} else if systemFlag {
			operations.OpenSystem(app)
		} else if len(args) > 0 {
			operations.Open(args, app)
		} else {
			crerrors.PrintCmdHelp("open")
		}

	},
}

func init() {

	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "open cfg folder")
	Cmd.Flags().BoolVarP(&systemFlag, "system", "s", false, "open system variable")
}
