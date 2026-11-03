package cmd

import (
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/gui"

	"github.com/spf13/cobra"
)

var desktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Launch native desktop window UI",
	Long:  "Launch a native window UI for crower (cgo-free)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return desktop.Run(cfg)
	},
}

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch interactive TUI GUI",
	Long:  "Launch an interactive terminal-based GUI for crower",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return gui.NewTUI(cfg).Run()
	},
}

func init() {
	rootCmd.AddCommand(desktopCmd)
	rootCmd.AddCommand(guiCmd)
}
