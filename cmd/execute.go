package cmd

import (
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/domain/listdom"

	"github.com/spf13/cobra"
)

var last bool

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute <command> [params...]",
	Short: "execute command",
	Long:  `execute command`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		listHandler := listdom.NewHandler(cfg)
		core := executedom.NewCore(cfg, listHandler)

		switch {
		case last:
			return core.ExecuteLast(operations.Execute, args)
		case createFlag:
			return core.ExecuteLast(operations.Create, args)
		case updateFlag:
			return core.ExecuteLast(operations.Update, args)
		}

		return core.ExecuteCommand(args)
	},
}

func init() {
	executeCmd.Flags().BoolVarP(&last, "last", "l", false, "execute recent executed command")
	executeCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "execute recent created command")
	executeCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "execute recent updated command")
}
