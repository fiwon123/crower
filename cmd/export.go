package cmd

import (
	"github.com/fiwon123/crower/internal/legacy/sharedom"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [name]",
	Short: "Export commands to a file.",
	Long:  `Export one command, or all of them, to a TOML file for sharing.`,
	Example: `	crower export
	crower export greet
	crower export --out ./commands.toml`,
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
	Short: "Import commands from a file.",
	Long: `Import commands from a TOML file. --overwrite replaces a command that
already exists with the same name.`,
	Example: `	crower import ./commands.toml
	crower import ./commands.toml --overwrite`,
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
	Short: "Share commands through a registry.",
	Long:  `List the packs a registry publishes, and pull one into the command set.`,
	Example: `	crower share list
	crower share pull tools`,
}

var shareListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List the packs available in the registry.",
	Long:    `List the packs available in the registry.`,
	Example: `	crower share list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return sharedom.NewCore(cfg).List()
	},
}

var sharePullCmd = &cobra.Command{
	Use:     "pull <name>",
	Short:   "Pull a pack from the registry.",
	Long:    `Pull a pack from the registry into the command set.`,
	Example: `	crower share pull tools`,
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
