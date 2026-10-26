package cmd

import (
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/opendom"

	"github.com/fiwon123/crower/pkg/utils"
	"github.com/spf13/cobra"
)

var allAlias []string

var scriptFlag bool
var processName string

// Cmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [name] [exec]",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		openHandler := opendom.NewHandler(cfg)
		core := createdom.NewCore(cfg, *openHandler)

		switch {
		case processName != "":
			return core.CreateProcess(processName, args)
		case fileFlag:
			return core.CreateFile(args)
		case folderFlag:
			return core.CreateFolder(args)
		case sysPathFlag:
			return core.CreateSystemPathVariable(args)
		case systemFlag:
			return core.CreateSystemVariable(args)
		case executeFlag:
			return core.CreateLastCommand(operations.Execute, args)
		case scriptFlag:
			return core.CreateScriptCommand(args)
		case len(args) == 0:
			return core.CreateCommand(allAlias, args)
		case utils.IsValidFilePath(args[0]):
			return core.CreateFile(args)
		case utils.IsValidFolderPath(args[0]):
			return core.CreateFolder(args)
		}

		return core.CreateCommand(allAlias, args)
	},
}

func init() {
	AddAllAliasFlag(createCmd, &allAlias)

	createCmd.Flags().StringVarP(&processName, "process", "p", "", "process name or pid")
	createCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	createCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	createCmd.Flags().BoolVar(&systemFlag, "system", false, "create system variable")
	createCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "create path variable")
	createCmd.Flags().BoolVar(&executeFlag, "execute", false, "create based on last executed command")
	createCmd.Flags().BoolVar(&scriptFlag, "script", false, "create a script file to help user create a complex command")
}
