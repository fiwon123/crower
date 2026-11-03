package cmd

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/legacy/deletedom"
	"github.com/fiwon123/crower/internal/legacy/listdom"

	"github.com/fiwon123/crower/pkg/utils"

	"github.com/spf13/cobra"
)

var historyFlagDelete bool

// Cmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [command]",
	Short: "Delete a command, a file or a folder.",
	Long: `Delete a command, a file or a folder.

A command is found by name or by alias:
	crower delete greet
	crower delete hi

The last operation can be undone instead:
	crower delete --create
	crower delete --update
	crower delete --execute
	crower delete --history

A file or a folder is deleted at the path you give:
	crower delete --file notes.txt
	crower delete --folder reports

With no arguments crower asks which command to delete.`,
	Example: `	crower delete
	crower delete greet
	crower delete --file notes.txt
	crower delete --folder reports
	crower delete --create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		listHandler := listdom.NewHandler(cfg)
		core := deletedom.NewCore(cfg, listHandler)

		switch {
		case createFlag:
			return core.DeleteLast(operation.Create)
		case updateFlag:
			return core.DeleteLast(operation.Update)
		case executeFlag:
			return core.DeleteLast(operation.Execute)
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
