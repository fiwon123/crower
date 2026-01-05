package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/cmd/copy"
	"github.com/fiwon123/crower/cmd/create"
	"github.com/fiwon123/crower/cmd/delete"
	"github.com/fiwon123/crower/cmd/execute"
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

	cmdsHelper "github.com/fiwon123/crower/internal/helper/cmds"
	"github.com/spf13/cobra"
)

var cfgFilePath string
var checkVersion bool
var checkNewVersion bool
var upgradeFlag bool

// Version is popualated when building with Makefile
var Version = "vx.x.x"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crower",
	Short: "A dev tool that manages system commands to help developers in their daily workflow.",
	Long: `A dev tool that manages system commands to help developers in their daily workflow.

It has useful operations like create, edit, remove, list and more.

Execute Command:
	- Use 'crower "command"' or 'cr "command"'
	- Use 'crower execute "command"' or 'cr execute "command"'`,
	Aliases: []string{"cr"},
	Args:    cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {

		cfgFilePath, _ := cmdsHelper.GetPersistentConfigFlag(cmd)

		if checkVersion {
			fmt.Println(Version)
			return
		}

		app := core.InitApp(cfgFilePath)

		if checkNewVersion {
			operations.CheckNewVersion(Version, app)
			return
		}

		if upgradeFlag {
			operations.UpgradeApp(Version, app)
			return
		}

		if len(args) > 0 {
			operations.Execute(args, app)
		} else {
			crerrors.PrintCmdHelp("")
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
	rootCmd.AddCommand(execute.Cmd)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(crerrors.GetNotUserHomeFoundString(), err)
	}

	defaultCfgFilePath := filepath.Join(homePath, "crower", "crower.yaml")

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", defaultCfgFilePath, "configuration path")

	// Flags
	rootCmd.Flags().BoolVarP(&checkVersion, "version", "v", false, "check current version")
	rootCmd.Flags().BoolVar(&upgradeFlag, "upgrade", false, "upgrade to new version")
	rootCmd.Flags().BoolVar(&checkNewVersion, "check", false, "check new version")
}
