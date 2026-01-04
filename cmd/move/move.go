package move

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/cterrors"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var fileFlag bool
var folderFlag bool

// Cmd represents the move command
var Cmd = &cobra.Command{
	Use:   "move",
	Short: "move file or folder to other location",
	Long: `move file or folder to other location

Example:
	crower move --file	"FILE_PATH" "OUTPUT_FOLDER_PATH"
	crower move --folder "FOLDER_PATH" "OUTPUT_FOLDER_PATH"
`,
	Run: func(cmd *cobra.Command, args []string) {
		if !fileFlag && !folderFlag {
			cterrors.PrintFileAndFolderFlagsNotUsed()
			return
		}

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if folderFlag {
			operations.MoveFolder(args, app)
		} else {
			operations.MoveFile(args, app)
		}
	},
}

func init() {
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "copy file using folder location and name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "copy folder using folder location and name")
}
