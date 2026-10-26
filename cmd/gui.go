package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/gui"
	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch interactive TUI GUI",
	Long:  "Launch an interactive terminal-based GUI for crower",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.InitApp(cfgFilePath)
		tui := gui.NewTUI(app)
		if err := tui.Run(); err != nil {
			app.Logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
