package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/spf13/cobra"
)

var desktopCmd = &cobra.Command{
	Use:   "desktop",
	Short: "Launch native desktop window UI",
	Long:  "Launch a native window UI for crower (cgo-free)",
	Run: func(cmd *cobra.Command, args []string) {
		a := app.InitApp(cfgFilePath)
		if err := desktop.Run(a); err != nil {
			a.Logger.Error(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(desktopCmd)
}
