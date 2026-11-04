package cmd

import (
	"github.com/fiwon123/crower/internal/application/executecommand"
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/session"
	"github.com/fiwon123/crower/internal/tui"

	"github.com/spf13/cobra"
)

// newExecutor builds the use case the terminal and desktop interfaces run
// through. They own their own interaction, so they use its Execute entry point,
// which neither prompts nor records history.
func newExecutor(cfg *session.Session) *executecommand.UseCase {
	session := newExecuteSession(cfg)

	return executecommand.New(session, session, session, session, session)
}

var desktopCmd = &cobra.Command{
	Use:     "desktop",
	Short:   "Open the native desktop window.",
	Long:    `Open a native window UI (cgo-free).`,
	Example: `	crower desktop`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		return desktop.Run(cfg.Catalog(), cfg.History, newExecutor(cfg))
	},
}

var guiCmd = &cobra.Command{
	Use:     "gui",
	Short:   "Open the terminal UI.",
	Long:    `Open an interactive terminal UI.`,
	Example: `	crower gui`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		return tui.NewTUI(cfg.Catalog(), cfg.History, newExecutor(cfg)).Run()
	},
}

func init() {
	rootCmd.AddCommand(desktopCmd)
	rootCmd.AddCommand(guiCmd)
}
