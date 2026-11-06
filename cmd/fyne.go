//go:build fyne

package cmd

import (
	"fmt"
	"os"

	"github.com/fiwon123/crower/internal/fyneui"

	"github.com/spf13/cobra"
)

// fyneCmd opens the Fyne window. The crower-gui executable invokes it through
// ExecuteFyne; it is also reachable as `crower fyne` in a tagged build.
var fyneCmd = &cobra.Command{
	Use:   "fyne",
	Short: "Open the Fyne window (CGO build).",
	Long: `Open the Fyne window.

This command is the entrypoint of the crower-gui executable, built with
-tags fyne and CGO. The default CLI binary does not include it.`,
	Example: "  crower-gui",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		fyneui.Run(cfg.Catalog(), cfg.History, fyneui.Deps{
			Exec:     newExecutor(cfg),
			Editor:   newFyneEditor(cfg),
			Reverter: newFyneReverter(cfg),
			Updater:  newFyneUpdater(cfg, Version),
			Tools:    newFyneTools(cfg),
			Groups:   newFyneGroups(cfg),
			Setups:   newFyneSetups(cfg),
			Window:   newFyneWindow(cfg.ConfigFolder()),
		})
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fyneCmd)
}

// ExecuteFyne is the crower-gui entrypoint. The root command parses the
// persistent flags (--config, --setup) and dispatches to the fyne command.
func ExecuteFyne(args []string) {
	rootCmd.SetArgs(append([]string{"fyne"}, args...))

	if err := rootCmd.Execute(); err != nil {
		message := "error: " + err.Error()
		fmt.Fprintln(streams.Err, message)

		// On Windows the executable is linked without a console, so the line
		// above goes nowhere a double-clicking user can see. The message box
		// is what turns "nothing happens" into a readable failure.
		showFyneStartupError(message)

		os.Exit(1)
	}
}
