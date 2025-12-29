package delete

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var name string
var allAlias []string
var exec string

var createFlag bool
var updateFlag bool
var executeFlag bool

var folderFlag bool
var fileFlag bool

// Cmd represents the delete command
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command, file or folder",
	Long: `delete command, file or folder

you can use interactive input just either typing 'crower delete' or using optional flags as name, alias and exec

Example:
	crower delete
	crower delete com1"
	crower delete --file "C:\Users\Test\Desktop\Test\new_file.txt"
	crower delete --folder "C:\Users\Test\Desktop\Test\new_folder"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		op := operation.Delete
		if createFlag {
			op = operation.DeleteCreate
		} else if updateFlag {
			op = operation.DeleteUpdate
		} else if executeFlag {
			op = operation.DeleteExecute
		} else if fileFlag {
			op = operation.DeleteFile
		} else if folderFlag {
			op = operation.DeleteFolder
		}

		fmt.Println(op)
		core.HandlePayload(
			payload.New(op, args, name, allAlias, exec),
			app,
		)
	},
}

func init() {
	cmdsHelper.AddNameFlag(Cmd, &name)
	cmdsHelper.AddAllAliasFlag(Cmd, &allAlias)
	cmdsHelper.AddExecFlag(Cmd, &exec)

	Cmd.Flags().BoolVarP(&createFlag, "create", "c", false, "delete recent created command")
	Cmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "delete recent updated command")
	Cmd.Flags().BoolVarP(&executeFlag, "execute", "x", false, "delete recent executed command")
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "delete file using folder location and name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "delete folder using folder location and name")
}
