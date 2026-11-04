//go:build fyne

package cmd

import (
	"github.com/fiwon123/crower/internal/fyneui"

	"github.com/spf13/cobra"
)

// fyneCmd opens the Fyne window. It exists only in a binary built with
// -tags fyne and CGO; the default binary keeps the cgo-free desktop.
var fyneCmd = &cobra.Command{
	Use:   "fyne",
	Short: "Open the Fyne desktop window (CGO build).",
	Long: `Open the Fyne desktop window.

This command exists only in a binary built with -tags fyne and CGO enabled.
The default crower binary ships the cgo-free desktop instead.`,
	Example: "  crower fyne",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		fyneui.Run(cfg.Catalog(), cfg.History, newExecutor(cfg))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(fyneCmd)
}
