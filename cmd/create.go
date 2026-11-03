package cmd

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/legacy/createdom"
	"github.com/fiwon123/crower/internal/legacy/opendom"

	"github.com/fiwon123/crower/pkg/utils"
	"github.com/spf13/cobra"
)

var allAlias []string

var scriptFlag bool
var processName string

// Cmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [name] [exec]",
	Short: "Create a command, a file or a folder.",
	Long: `Create a command, a file or a folder.

A command has a name, optional aliases, and the line it runs:
	crower create greet "echo hello"
	crower create greet "echo hello" --alias hi

It can come from a running process, or be a script:
	crower create -p server nginx
	crower create --script deploy

A file or a folder is created at the path you give:
	crower create --file notes.txt
	crower create --folder reports

With no arguments crower asks for each field. --execute makes a command from
the last command you ran.`,
	Example: `	crower create
	crower create greet "echo hello"
	crower create greet "echo hello" --alias hi
	crower create --file notes.txt
	crower create --folder reports
	crower create -p server nginx
	crower create --script deploy`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		openHandler := opendom.NewHandler(cfg)
		core := createdom.NewCore(cfg, *openHandler)

		switch {
		case processName != "":
			return core.CreateProcess(processName, allAlias, args)
		case fileFlag:
			return core.CreateFile(args)
		case folderFlag:
			return core.CreateFolder(args)
		case sysPathFlag:
			return core.CreateSystemPathVariable(args)
		case systemFlag:
			return core.CreateSystemVariable(args)
		case executeFlag:
			return core.CreateLastCommand(operation.Execute, allAlias, args)
		case scriptFlag:
			return core.CreateScriptCommand(allAlias, args)
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

	createCmd.Flags().StringVarP(&processName, "process", "p", "", "new command name for a process; pass the process name or pid as an argument")
	createCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	createCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	createCmd.Flags().BoolVar(&systemFlag, "system", false, "create system variable")
	createCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "create path variable")
	createCmd.Flags().BoolVar(&executeFlag, "execute", false, "create based on last executed command")
	createCmd.Flags().BoolVar(&scriptFlag, "script", false, "create a script file to help user create a complex command")
}
