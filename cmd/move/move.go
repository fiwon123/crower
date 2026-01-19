package move

import (
	"github.com/fiwon123/crower/internal/core"
	moveoperations "github.com/fiwon123/crower/internal/core/operations/move"
	"github.com/fiwon123/crower/internal/crowererrors"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the move command
var Cmd = &cobra.Command{
	Use:   "move",
	Short: "move file or folder to other location",
	Long: `move file or folder to other location

Examples:
	crower move "C:\Users\Test\Desktop\Test\file.txt" "C:\Users\Test\Desktop\Test\Out"
	crower move "C:\Users\Test\Desktop\Test\file_1.txt" "C:\Users\Test\Desktop\Test\file_2.txt" "C:\Users\Test\Desktop\Test\Out"
	crower move "C:\Users\Test\Desktop\Test\Folder" "C:\Users\Test\Desktop\Test\Out"
	crower move "C:\Users\Test\Desktop\Test\Folder_1" "C:\Users\Test\Desktop\Test\Folder_2" "C:\Users\Test\Desktop\Test\Out"
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if len(args) > 0 {
			moveoperations.Move(args, app)
		} else {
			crowererrors.PrintCmdHelp("move", app)
		}
	},
}

func init() {
}
