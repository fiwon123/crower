package cmd

import (
		"github.com/fiwon123/crower/internal/app"
		"github.com/fiwon123/crower/internal/domain/sharedom"
		"github.com/spf13/cobra"
	)

var exportCmd = &cobra.Command{
	Use:   "export [name]",
	Short: "Export commands to a file",
	Long:  "Export a command or all commands to a TOML file for sharing",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.InitApp(cfgFilePath)
		core := sharedom.NewCore(app)
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		out := ""
		if f := cmd.Flag("out"); f != nil {
			out = f.Value.String()
		}
		if err := core.Export(name, out); err != nil {
			app.Logger.Error(err.Error())
			return
		}
		app.Logger.Info("Exported successfully")
	},
}

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import commands from a file",
	Long:  "Import commands from a TOML file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			app := app.InitApp(cfgFilePath)
			app.Logger.Error("file required")
			return
		}
		app := app.InitApp(cfgFilePath)
		core := sharedom.NewCore(app)
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		if err := core.Import(args[0], overwrite); err != nil {
			app.Logger.Error(err.Error())
			return
		}
		app.Logger.Info("Imported successfully")
	},
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share commands via registry",
	Long:  "List and pull commands from a registry",
}

var shareListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available packs in registry",
	Run: func(cmd *cobra.Command, args []string) {
		app := app.InitApp(cfgFilePath)
		core := sharedom.NewCore(app)
		if err := core.List(); err != nil {
			app.Logger.Error(err.Error())
			return
		}
	},
}

var sharePullCmd = &cobra.Command{
	Use:   "pull <name>",
	Short: "Pull a pack from registry",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			app := app.InitApp(cfgFilePath)
			app.Logger.Error("pack name required")
			return
		}
		app := app.InitApp(cfgFilePath)
		core := sharedom.NewCore(app)
		if err := core.Pull(args[0]); err != nil {
			app.Logger.Error(err.Error())
			return
		}
	},
}

func init() {
	exportCmd.Flags().StringP("out", "o", "", "Output file path")
	importCmd.Flags().Bool("overwrite", false, "Overwrite existing commands")
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
	rootCmd.AddCommand(shareCmd)
	shareCmd.AddCommand(shareListCmd)
	shareCmd.AddCommand(sharePullCmd)
}
