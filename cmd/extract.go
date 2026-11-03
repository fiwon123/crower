package cmd

import (
	"github.com/fiwon123/crower/internal/domain/extractdom"

	"github.com/spf13/cobra"
)

var outDirFlag string

// Cmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract <archive...>",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "compressed file paths")
		}

		return extractdom.NewCore(cfg).Extract(args, outDirFlag)
	},
}

func init() {
	extractCmd.Flags().StringVarP(&outDirFlag, "out", "o", "", "out folder where to be extracted")
}
