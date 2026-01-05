package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/cmd/copy"
	"github.com/fiwon123/crower/cmd/create"
	"github.com/fiwon123/crower/cmd/delete"
	"github.com/fiwon123/crower/cmd/extract"
	"github.com/fiwon123/crower/cmd/list"
	"github.com/fiwon123/crower/cmd/move"
	"github.com/fiwon123/crower/cmd/open"
	"github.com/fiwon123/crower/cmd/reset"
	"github.com/fiwon123/crower/cmd/revert"
	"github.com/fiwon123/crower/cmd/search"
	"github.com/fiwon123/crower/cmd/update"
	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/core/operations"
	"github.com/fiwon123/crower/internal/crerrors"

	"github.com/fiwon123/crower/internal/data/state"

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var cfgFilePath string
var checkVersion bool

var last bool
var createFlag bool
var updateFlag bool

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
		if last {
			operations.ExecuteLast(state.Execute, args, app)
		} else if createFlag {
			operations.ExecuteLast(state.Create, args, app)
		} else if updateFlag {
			operations.ExecuteLast(state.Update, args, app)
		} else {
			operations.Execute("", args, app)
		}

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
	rootCmd.AddCommand(search.Cmd)
	rootCmd.AddCommand(extract.Cmd)
	rootCmd.AddCommand(copy.Cmd)
	rootCmd.AddCommand(move.Cmd)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(crerrors.GetNotUserHomeFoundString(), err)
	}

	defaultCfgFilePath := filepath.Join(homePath, "crower", "crower.yaml")

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", defaultCfgFilePath, "configuration path")

	// Flags
	rootCmd.Flags().BoolVarP(&checkVersion, "version", "v", false, "check current version")
	rootCmd.Flags().BoolVarP(&last, "last", "l", false, "execute recent executed command")
	rootCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "execute recent created command")
	rootCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "execute recent updated command")
}
