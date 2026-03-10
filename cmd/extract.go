package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/crowererrors"
	"github.com/fiwon123/crower/internal/domain/extractdom"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var outDirFlag string

// Cmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...)",
	Long: `extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...)

Type 'crower/cr extract "FILE_PATH"' to extract all content to the current folder
You can specify out to change output folder: 'crower/cr extract "COMPRESSED_FILE_PATH" --out/-o "OUTPUT_FOLDER_PATH"'
If you have more than 1 compressed file to extract you can put * as the name file

Examples:
	crower extract "C:\Users\Test\Desktop\Test\files.zip"
	crower extract "C:\Users\Test\Desktop\Test\files.zip" -o "C:\Users\Test\Desktop\Out\"
	crower extract "C:\Users\Test\Desktop\Test\*.zip"
	crower extract "C:\Users\Test\Desktop\Test\*.zip" -o "C:\Users\Test\Desktop\Out\"

	`,
	Run: func(cmd *cobra.Command, args []string) {
		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		app := app.InitApp(cfgFilePath)

		core := extractdom.NewCore(app)

		if len(args) > 0 {
			core.Extract(args, outDirFlag)
		} else {
			crowererrors.PrintCmdHelp("extract", app)
		}
	},
}

func init() {
	extractCmd.Flags().StringVarP(&outDirFlag, "out", "o", "", "out folder where to be extracted")
}
