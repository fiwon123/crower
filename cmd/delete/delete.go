package delete

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/cterrors"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/pkg/utils"

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

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
delete commands:
	- Use interactive input just typing 'crower delete' without parameter
	- Using either 'crower delete "COMMAND_NAME"' or 'crower delete "COMMAND_ALIAS"'

delete file:
	- Using 'crower delete "FILE_PATH"'
	- Using flag --file for multiple paths 'crower delete --file "FILE_PATH_1" "FILE_PATH_2" "FILE_PATH_3"'

delete folder:
	- Using 'crower delete "FOLDER_PATH"'
	- Using flag --folder for multiple paths 'crower delete --folder "FOLDER_PATH_1" "FOLDER_PATH_2" "FOLDER_PATH_3"'

Example:
	crower delete
	crower delete com_name
	crower delete com_alias
	crower delete --file "C:\Users\Test\Desktop\Test\new_file.txt"
	crower delete --file "C:\Users\Test\Desktop\Test\new_file_1.txt" "C:\Users\Test\Desktop\Test\new_file_2.txt"
	crower delete --folder "C:\Users\Test\Desktop\Test\new_folder"
	crower delete --folder "C:\Users\Test\Desktop\Test\new_folder_1" "C:\Users\Test\Desktop\Test\new_folder_2"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if createFlag {
			operations.DeleteLast(state.Create, app)
		} else if updateFlag {
			operations.DeleteLast(state.Update, app)
		} else if executeFlag {
			operations.DeleteLast(state.Execute, app)
		} else if fileFlag {
			operations.DeleteFile(args, app)
		} else if folderFlag {
			operations.DeleteFolder(args, app)
		} else if len(args) > 0 {
			if utils.IsValidFilePath(args[0]) {
				operations.DeleteFile(args, app)
			} else if utils.IsValidFolderPath(args[0]) {
				operations.DeleteFolder(args, app)
			} else {
				operations.Delete(args, app)
			}
		} else {
			cterrors.PrintCmdHelp("delete")
		}

	},
}

func init() {
	Cmd.Flags().BoolVarP(&createFlag, "create", "c", false, "delete recent created command")
	Cmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "delete recent updated command")
	Cmd.Flags().BoolVarP(&executeFlag, "execute", "x", false, "delete recent executed command")
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
}
