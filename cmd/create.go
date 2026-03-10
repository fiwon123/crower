package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/opendom"

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/fiwon123/crower/pkg/crowerutils"
	"github.com/spf13/cobra"
)

var allAlias []string

var scriptFlag bool
var processName string

// Cmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a command, file or folder",
	Long: `create a command, file or folder

create commands:
	- Use interactive input just typing 'crower create' without arguments
	- Using 'crower create "COMAND_NAME" "COMMAND_EXEC"'
	- Using flag --alias in 'crower create' to already put an alias

create file:
	- Using 'crower create "FILE_PATH"'
	- Using flag --file for multiple paths 'crower create --file "FILE_PATH_1" "FILE_PATH_2" "FILE_PATH_3"'

create folder:
	- Using 'crower create "FOLDER_PATH"'
	- Using flag --folder for multiple paths 'crower create --folder "FOLDER_PATH_1" "FOLDER_PATH_2" "FOLDER_PATH_3"'

Example:
	crower create
	crower create com1 "'echo com1'"
	crower create "C:\Users\Test\Desktop\Test\new_file.txt"
	crower create --file "C:\Users\Test\Desktop\Test\new_file_1.txt" "C:\Users\Test\Desktop\Test\new_file_2.txt"
	crower create "C:\Users\Test\Desktop\Test\new_folder"
	crower create --folder "C:\Users\Test\Desktop\Test\new_folder_1" "C:\Users\Test\Desktop\Test\new_folder_2"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		openHandler := opendom.NewHandler(app)
		core := createdom.NewCore(app, *openHandler)

		if processName != "" {
			core.CreateProcess(processName, args)
		} else if fileFlag {
			core.CreateFile(args)
		} else if folderFlag {
			core.CreateFolder(args)
		} else if sysPathFlag {
			core.CreateSystemPathVariable(args)
		} else if systemFlag {
			core.CreateSystemVariable(args)
		} else if executeFlag {
			core.CreateLastCommand(operationsdata.Execute, args)
		} else if scriptFlag {
			core.CreateScriptCommand(args)
		} else if len(args) > 0 {
			if crowerutils.IsValidFilePath(args[0]) {
				core.CreateFile(args)
			} else if crowerutils.IsValidFolderPath(args[0]) {
				core.CreateFolder(args)
			} else {
				core.CreateCommand(allAlias, args)
			}
		} else {
			core.CreateCommand(allAlias, args)
		}

	},
}

func init() {
	cmdsHelper.AddAllAliasFlag(createCmd, &allAlias)

	createCmd.Flags().StringVarP(&processName, "process", "p", "", "process name or pid")
	createCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	createCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	createCmd.Flags().BoolVar(&systemFlag, "system", false, "create system variable")
	createCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "create path variable")
	createCmd.Flags().BoolVar(&executeFlag, "execute", false, "create based on last executed command")
	createCmd.Flags().BoolVar(&scriptFlag, "script", false, "create a script file to help user create a complex command")
}
