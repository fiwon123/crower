package create

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/fiwon123/crower/pkg/utils"
	"github.com/spf13/cobra"
)

var allAlias []string
var process string

var folderFlag bool
var fileFlag bool
var sysPathFlag bool
var systemFlag bool

// Cmd represents the create command
var Cmd = &cobra.Command{
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

		app := core.InitApp(cfgFilePath)

		if process != "" {
			operations.CreateProcess(process, args, app)
		} else if fileFlag {
			operations.CreateFile(args, app)
		} else if folderFlag {
			operations.CreateFolder(args, app)
		} else if sysPathFlag {
			operations.CreateSystemPathVariable(args, app)
		} else if systemFlag {
			operations.CreateSystemVariable(args, app)
		} else if len(args) > 0 {
			if utils.IsValidFilePath(args[0]) {
				operations.CreateFile(args, app)
			} else if utils.IsValidFolderPath(args[0]) {
				operations.CreateFolder(args, app)
			} else {
				operations.CreateCommand(allAlias, args, app)
			}
		} else {
			operations.CreateCommand(allAlias, args, app)
		}

	},
}

func init() {
	cmdsHelper.AddAllAliasFlag(Cmd, &allAlias)

	Cmd.Flags().StringVarP(&process, "process", "p", "", "process name or pid")
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	Cmd.Flags().BoolVar(&systemFlag, "system", false, "create system variable")
	Cmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "create path variable")
}
