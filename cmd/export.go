package cmd

import (
	"github.com/fiwon123/crower/internal/domain/sharedom"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [name]",
	Short: "Export commands to a file",
	Long:  "Export a command or all commands to a TOML file for sharing",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		out, err := cmd.Flags().GetString("out")
		if err != nil {
			return err
		}

		return sharedom.NewCore(cfg).Export(name, out)
	},
}

var importCmd = &cobra.Command{
	Use:   "import <file>",
	Short: "Import commands from a file",
	Long:  "Import commands from a TOML file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return missingArgs(cmd, "a file to import")
		}

		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		overwrite, err := cmd.Flags().GetBool("overwrite")
		if err != nil {
			return err
		}

		return sharedom.NewCore(cfg).Import(args[0], overwrite)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return sharedom.NewCore(cfg).List()
	},
}

var sharePullCmd = &cobra.Command{
	Use:   "pull <name>",
	Short: "Pull a pack from registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return missingArgs(cmd, "a pack name")
		}

		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return sharedom.NewCore(cfg).Pull(args[0])
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
