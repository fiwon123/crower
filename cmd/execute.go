package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/executecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"

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

		request := executecommand.Request{}

		// The flags decide which command to run; the use case owns how.
		switch {
		case last:
			request.Last, request.Params = operation.Execute, args
		case createFlag:
			request.Last, request.Params = operation.Create, args
		case updateFlag:
			request.Last, request.Params = operation.Update, args
		case len(args) > 0:
			request.Key, request.Params = args[0], args[1:]
		}

		session := newExecuteSession(cfg)

		return executecommand.New(session, session, session, session, session).Run(request)
	},
}

// executeSession adapts the configuration to the ports the execute use case
// declares.
//
// It embeds legacySession for the catalog, the diagnostics stream, the result
// stream, the debug log, the confirmation and the history entry, and adds the
// shell runner and the two execution notes.
type executeSession struct {
	legacySession

	runner *system.Runner
}

// newExecuteSession wires the configuration and the shell runner.
func newExecuteSession(cfg *app.Config) executeSession {
	return executeSession{
		legacySession: newLegacySession(cfg),
		runner:        system.NewRunner(),
	}
}

// Run runs one shell line in the command's shell and working directory.
func (s executeSession) Run(line, shellName, dir string) (string, error) {
	return s.runner.Run(line, shellName, dir)
}

// RunNested runs a line that invokes crower again, with the depth guard.
func (s executeSession) RunNested(line, shellName, dir string) (string, error) {
	return s.runner.RunNested(line, shellName, dir)
}

// ExecuteCommandNote is the history note for `execute <command>`.
func (s executeSession) ExecuteCommandNote(command catalog.Command, params []string) string {
	return notes.ExecuteCommandNote(command, params)
}

// ExecuteLastNote is the history note for the --last forms.
func (s executeSession) ExecuteLastNote(op operation.Operation, command catalog.Command, args []string) string {
	return notes.ExecuteLastNote(op, command, args)
}

func init() {
	executeCmd.Flags().BoolVarP(&last, "last", "l", false, "execute recent executed command")
	executeCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "execute recent created command")
	executeCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "execute recent updated command")
}
