package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/extractdom"

	"github.com/spf13/cobra"
)

var outDirFlag string

// Cmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract <archive...>",
	Short: "Extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...).",
	Long: `Extract compressed files (.zip, .tar, .tar.gz, .7z, .rar...).

By default the content goes to the archive's own folder; --out changes it.

	crower extract files.zip
	crower extract files.zip --out ./build
	crower extract "*.zip"`,
	Example: `	crower extract files.zip
	crower extract files.zip --out ./build
	crower extract "*.zip" --out ./build`,
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
