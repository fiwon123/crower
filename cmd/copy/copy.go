package copy

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var fileFlag bool
var folderFlag bool

// Cmd represents the copy command
var Cmd = &cobra.Command{
	Use:   "copy",
	Short: "copy file or folder to other place",
	Long:  `copy file or folder to other place`,
	Run: func(cmd *cobra.Command, args []string) {

		if !fileFlag && !folderFlag {
			fmt.Println("file and folder flag not used")
			return
		}

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		op := operation.CopyFile
		if folderFlag {
			op = operation.CopyFolder
		}

		core.HandlePayload(
			payload.New(op, args, "", []string{}, ""),
			app,
		)
	},
}

func init() {
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "copy file using folder location and name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "copy folder using folder location and name")
}
