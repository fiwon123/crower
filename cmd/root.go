package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/cmd/create"
	"github.com/fiwon123/crower/cmd/delete"
	"github.com/fiwon123/crower/cmd/list"
	"github.com/fiwon123/crower/cmd/open"
	"github.com/fiwon123/crower/cmd/reset"
	"github.com/fiwon123/crower/cmd/revert"
	"github.com/fiwon123/crower/cmd/update"
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/data/payload"
	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var cfgFilePath string
var checkVersion bool

// Version is popualated when building with Makefile
var Version = "vx.x.x"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crower",
	Short: "A dev tool that manages system commands to help developers in their daily workflow.",
	Long: `A dev tool that manages system commands to help developers in their daily workflow.

It has useful operations like create, edit, remove, list and more.

By default after created your first command just use it by typing "crower 'command'" or "cr 'command'"`,
	Aliases: []string{"cr"},
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		if checkVersion {
			fmt.Println(Version)
			return
		}

		app := core.InitApp(cfgFilePath)

		core.HandlePayload(
			payload.New(operation.Execute, args, "", nil, ""),
			app,
		)

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
	rootCmd.AddCommand(create.Cmd)
	rootCmd.AddCommand(update.Cmd)
	rootCmd.AddCommand(delete.Cmd)
	rootCmd.AddCommand(list.Cmd)
	rootCmd.AddCommand(open.Cmd)
	rootCmd.AddCommand(reset.Cmd)
	rootCmd.AddCommand(revert.Cmd)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error could not get user home directory, ", err)
	}

	defaultCfgFilePath := filepath.Join(homePath, "crower", "crower.yaml")

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", defaultCfgFilePath, "configuration path")

	// Flags
	rootCmd.Flags().BoolVarP(&checkVersion, "version", "v", false, "check current version")
}
