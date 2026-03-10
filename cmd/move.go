package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/movedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

// Cmd represents the move command
var moveCmd = &cobra.Command{
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
		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := movedom.NewCore(app)

		if len(args) > 0 {
			core.Move(args)
		} else {
			errors.PrintCmdHelp("move", app)
		}
	},
}

func init() {
}
