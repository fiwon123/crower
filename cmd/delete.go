package cmd

import (
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/deletedom"
	"github.com/fiwon123/crower/internal/domain/listdom"

	"github.com/fiwon123/crower/pkg/utils"

	"github.com/spf13/cobra"
)

var historyFlagDelete bool

// Cmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [command]",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		listHandler := listdom.NewHandler(cfg)
		core := deletedom.NewCore(cfg, listHandler)

		switch {
		case createFlag:
			return core.DeleteLast(operations.Create)
		case updateFlag:
			return core.DeleteLast(operations.Update)
		case executeFlag:
			return core.DeleteLast(operations.Execute)
		case fileFlag:
			return core.DeleteFile(args)
		case folderFlag:
			return core.DeleteFolder(args)
		case systemFlag:
			return core.DeleteSystemVariable(args)
		case sysPathFlag:
			return core.DeleteSystemPathVariable(args)
		case historyFlagDelete:
			return core.DeleteHistoryContent(args)
		case len(args) == 0:
			return core.Delete(args)
		case utils.IsValidFilePath(args[0]):
			return core.DeleteFile(args)
		case utils.IsValidFolderPath(args[0]):
			return core.DeleteFolder(args)
		}

		return core.Delete(args)
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "delete recent created command")
	deleteCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "delete recent updated command")
	deleteCmd.Flags().BoolVarP(&executeFlag, "execute", "x", false, "delete recent executed command")
	deleteCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	deleteCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	deleteCmd.Flags().BoolVar(&systemFlag, "system", false, "delete system variable")
	deleteCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "delete path variable")
	deleteCmd.Flags().BoolVar(&historyFlagDelete, "history", false, "delete a history content and linked backup cfg")
}
