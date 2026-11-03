package cmd

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/legacy/executedom"
	"github.com/fiwon123/crower/internal/legacy/listdom"

	"github.com/spf13/cobra"
)

var last bool

// executeCmd represents the execute command
var executeCmd = &cobra.Command{
	Use:   "execute <command> [params...]",
	Short: "Run a managed command.",
	Long: `Run a managed command, by name or by alias.

	crower execute greet
	crower greet

The last touched command can be run instead:
	crower execute --last
	crower execute --create
	crower execute --update

Arguments after the command are applied to its steps.`,
	Example: `	crower execute greet
	crower execute greet --help
	crower execute --last`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		listHandler := listdom.NewHandler(cfg)
		core := executedom.NewCore(cfg, listHandler)

		switch {
		case last:
			return core.ExecuteLast(operation.Execute, args)
		case createFlag:
			return core.ExecuteLast(operation.Create, args)
		case updateFlag:
			return core.ExecuteLast(operation.Update, args)
		}

		return core.ExecuteCommand(args)
	},
}

func init() {
	executeCmd.Flags().BoolVarP(&last, "last", "l", false, "execute recent executed command")
	executeCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "execute recent created command")
	executeCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "execute recent updated command")
}
