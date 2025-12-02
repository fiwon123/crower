package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crower",
	Short: "A dev tool that manages system commands to help developers in their daily workflow.",
	Long: `A dev tool that manages system commands by executing commands via custom aliases and
managing it with useful operations like add, edit, remove, list and more.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:s
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args", args)
		fmt.Println("cfg", cfgFile)

		// init app vars

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$HOME/.crower.yaml", "config file (default is $HOME/.crower.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}
