package open

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
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

		op := operation.Open
		if folderFlag {
			op = operation.OpenFolder
		} else if systemFlag {
			op = operation.OpenSystem
		}

		core.HandlePayload(
			payload.New(op, args, "", nil, ""),
			app,
		)
	},
}

func init() {

	Cmd.Flags().BoolVarP(&folderFlag, "folder", "f", false, "open cfg folder")
	Cmd.Flags().BoolVarP(&systemFlag, "system", "s", false, "open system variable")
}
