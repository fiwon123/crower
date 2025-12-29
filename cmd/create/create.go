package create

import (
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var name string
var allAlias []string
var exec string
var process string

var folderFlag bool
var fileFlag bool

// Cmd represents the create command
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "create a command, file or folder",
	Long: `create a command, file or folder

you can use interactive input just either typing 'crower create' or using optional flags as name, alias and exec

Example:
	crower create
	crower create com1 "'echo com1'"
	crower create --file "C:\Users\Test\Desktop\Test\new_file.txt"
	crower create --folder "C:\Users\Test\Desktop\Test\new_folder"`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := core.InitApp(cfgFilePath)

		if process != "" {
			name = process
		}

		op := operation.Create
		if fileFlag {
			op = operation.CreateFile
		} else if folderFlag {
			op = operation.CreateFolder
		}

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

	Cmd.Flags().StringVarP(&process, "process", "p", "", "process name or pid")
	Cmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "create file using folder location and name")
	Cmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "create folder using folder location and name")
}
