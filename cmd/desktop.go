package cmd

import (
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/tui"

	"github.com/spf13/cobra"
)

var desktopCmd = &cobra.Command{
	Use:     "desktop",
	Short:   "Open the native desktop window.",
	Long:    `Open a native window UI (cgo-free).`,
	Example: `	crower desktop`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return desktop.Run(cfg)
	},
}

var guiCmd = &cobra.Command{
	Use:     "gui",
	Short:   "Open the terminal UI.",
	Long:    `Open an interactive terminal UI.`,
	Example: `	crower gui`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return tui.NewTUI(cfg).Run()
	},
}

func init() {
	rootCmd.AddCommand(desktopCmd)
	rootCmd.AddCommand(guiCmd)
}
