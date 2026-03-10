package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	opendom "github.com/fiwon123/crower/internal/domain/open"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

// Cmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open current configuration file path, folder and system",
	Long: `open current configuration file path, folder and system

Examples:
	crower open
	crower open --folder "FOLDER_PATH"
	crower open --system
`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := opendom.NewCore(app)

		if folderFlag {
			core.OpenFolder(args)
		} else if fileFlag {
			core.OpenFile(args)
		} else if systemFlag {
			core.OpenSystem()
		} else if len(args) > 0 {
			core.Open(args)
		} else {
			crowererrors.PrintCmdHelp("open", app)
		}

	},
}

func init() {

	openCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "open cfg file or other file")
	openCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "open cfg folder or other folder")
	openCmd.Flags().BoolVarP(&systemFlag, "system", "s", false, "open system variable")
}
