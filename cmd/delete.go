package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/deletedom"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/pkg/crowerutils"

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var historyFlag bool

// Cmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command, file or folder",
	Long: `delete command, file or folder

delete commands:
	- Use interactive input just typing 'crower delete' without arguments
	- Using either 'crower delete "COMMAND_NAME"' or 'crower delete "COMMAND_ALIAS"'
	- using flags --create, --update or --execute to delete last operation flag

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

		app := app.InitApp(cfgFilePath)

		listHandler := listdom.NewHandler(app)
		core := deletedom.NewCore(app, listHandler)

		if createFlag {
			core.DeleteLast(operations.Create)
		} else if updateFlag {
			core.DeleteLast(operations.Update)
		} else if executeFlag {
			core.DeleteLast(operations.Execute)
		} else if fileFlag {
			core.DeleteFile(args)
		} else if folderFlag {
			core.DeleteFolder(args)
		} else if systemFlag {
			core.DeleteSystemVariable(args)
		} else if sysPathFlag {
			core.DeleteSystemPathVariable(args)
		} else if historyFlag {
			core.DeleteHistoryContent(args)
		} else if len(args) > 0 {
			if crowerutils.IsValidFilePath(args[0]) {
				core.DeleteFile(args)
			} else if crowerutils.IsValidFolderPath(args[0]) {
				core.DeleteFolder(args)
			} else {
				core.Delete(args)
			}
		} else {
			core.Delete(args)
		}

	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "delete recent created command")
	deleteCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "delete recent updated command")
	deleteCmd.Flags().BoolVarP(&executeFlag, "execute", "x", false, "delete recent executed command")
	deleteCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	deleteCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	deleteCmd.Flags().BoolVar(&systemFlag, "system", false, "create system variable")
	deleteCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "create path variable")
	deleteCmd.Flags().BoolVar(&historyFlag, "history", false, "delete a history content and linked backup cfg")
}
