package cmd

import (
	"github.com/fiwon123/crower/internal/domain/copydom"

	"github.com/spf13/cobra"
)

// Cmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy <source...> <destination>",
	Short: "copy file or folder to other location",
	Long: `copy file or folder to other location

Examples:
	crower copy "C:\Users\Test\Desktop\Test\file.txt" "C:\Users\Test\Desktop\Test\Out"
	crower copy "C:\Users\Test\Desktop\Test\file_1.txt" "C:\Users\Test\Desktop\Test\file_2.txt" "C:\Users\Test\Desktop\Test\Out"
	crower copy "C:\Users\Test\Desktop\Test\Folder" "C:\Users\Test\Desktop\Test\Out"
	crower copy "C:\Users\Test\Desktop\Test\Folder_1" "C:\Users\Test\Desktop\Test\Folder_2" "C:\Users\Test\Desktop\Test\Out"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "source paths and a destination")
		}

		return copydom.NewCore(cfg).Copy(args)
	},
}
