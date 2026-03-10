package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/checkdom"
	"github.com/fiwon123/crower/internal/domain/executedom"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/domain/upgradedom"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

var cfgFilePath string
var checkVersion bool
var checkNewVersion bool
var upgradeFlag bool

// Shared
var executeFlag bool
var fileFlag bool
var folderFlag bool
var sysPathFlag bool
var systemFlag bool
var createFlag bool
var updateFlag bool
var processFlag bool

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

		cfgFilePath, _ := helper.GetPersistentConfigFlag(cmd)

		if checkVersion {
			fmt.Println(Version)
			return
		}

		app := app.InitApp(cfgFilePath)
		checkCore := checkdom.NewCore(app)
		upgradeHandler := upgradedom.NewHandler(app)
		listHandler := listdom.NewHandler(app)

		executeCore := executedom.NewCore(app, listHandler)

		if checkNewVersion {
			checkCore.CheckNewVersion(Version)
			return
		}

		if upgradeFlag {
			upgradeHandler.UpgradeApp(Version)
			return
		}

		if len(args) > 0 {
			executeCore.ExecuteCommand(args)
		} else {
			errors.PrintCmdHelp("", app)
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
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(openCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(revertCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(extractCmd)
	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(moveCmd)
	rootCmd.AddCommand(executeCmd)
	rootCmd.AddCommand(restoreCmd)

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(errors.GetNotUserHomeFoundString(), err)
	}

	defaultCfgFilePath := filepath.Join(homePath, "crower", "crower.yaml")

	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", defaultCfgFilePath, "configuration path")

	// Flags
	rootCmd.Flags().BoolVarP(&checkVersion, "version", "v", false, "check current version")
	rootCmd.Flags().BoolVar(&upgradeFlag, "upgrade", false, "upgrade to new version")
	rootCmd.Flags().BoolVar(&checkNewVersion, "check", false, "check new version")
}
