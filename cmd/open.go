package cmd

import (
	"github.com/fiwon123/crower/internal/domain/opendom"

	"github.com/spf13/cobra"
)

// Cmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [path...]",
	Short: "open current configuration file path, folder and system",
	Long: `open current configuration file path, folder and system

Examples:
	crower open
	crower open --file "FILE_PATH"
	crower open --folder "FOLDER_PATH"
	crower open --system
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		core := opendom.NewCore(cfg)

		switch {
		case folderFlag:
			return core.OpenFolder(args)
		case fileFlag:
			return core.OpenFile(args)
		case systemFlag:
			return core.OpenSystem()
		case len(args) > 0:
			return core.Open(args)
		}

		return missingArgs(cmd, "a path, or one of --file, --folder, --system")
	},
}

func init() {
	openCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "open cfg file or other file")
	openCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "open cfg folder or other folder")
	openCmd.Flags().BoolVarP(&systemFlag, "system", "s", false, "open system variable")
}
