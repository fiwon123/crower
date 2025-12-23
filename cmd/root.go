package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data"
	"github.com/spf13/cobra"
)

var cfgFilePath string
var addOp bool
var deleteOp bool
var updateOp bool
var listOp bool
var resetOp bool
var index int
var name string
var exec string
var alias []string
var openOp bool
var processOp bool
var historyOp bool

var checkVersion bool

// Version is popualted when building with Makefile
var Version = "vx.x.x"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crower",
	Short: "A dev tool that manages system commands to help developers in their daily workflow.",
	Long: `A dev tool that manages system commands by executing commands via custom aliases and
managing it with useful operations like add, edit, remove, list and more.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:s
	Run: func(cmd *cobra.Command, args []string) {

		if checkVersion {
			fmt.Println(Version)
			return
		}

		fmt.Println("cfg", cfgFilePath)
		app := core.InitApp(cfgFilePath)

		op := getOperation()

		core.HandlePayload(
			data.Payload{
				Op:    op,
				Args:  args,
				Name:  name,
				Alias: alias,
				Exec:  exec,
			},
			app,
		)

	},
}

func getOperation() data.Operation {
	if addOp {

		if processOp {
			return data.AddProcess
		}

		return data.AddOp

	} else if listOp {
		return data.ListOp
	} else if resetOp {
		return data.ResetOp
	} else if deleteOp {
		return data.DeleteOp
	} else if updateOp {
		return data.UpdateOp
	} else if openOp {
		return data.OpenOp
	} else if processOp {
		return data.ProcessOp
	} else if historyOp {
		return data.HistoryOp
	}

	return data.ExecuteOp

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

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error could not get user home directory, ", err)
	}

	defaultCfgFilePath := filepath.Join(homePath, "crower", "crower.yaml")

	rootCmd.Flags().StringVar(&cfgFilePath, "config", defaultCfgFilePath, "config file (default is $HOME/.crower.yaml)")
	rootCmd.Flags().IntVarP(&index, "index", "i", 0, "command index")
	rootCmd.Flags().BoolVar(&addOp, "add", false, "add a command (--add ip ifconfig)")
	rootCmd.Flags().BoolVarP(&checkVersion, "version", "v", false, "check current version")
	rootCmd.Flags().BoolVar(&listOp, "list", false, "list all commands")
	rootCmd.Flags().BoolVar(&resetOp, "reset", false, "reset all commands")
	rootCmd.Flags().BoolVar(&updateOp, "update", false, "update command")
	rootCmd.Flags().BoolVar(&deleteOp, "delete", false, "delete commands")
	rootCmd.Flags().BoolVar(&openOp, "open", false, "open cfg file path")
	rootCmd.Flags().BoolVar(&processOp, "process", false, "list all process")
	rootCmd.Flags().BoolVar(&historyOp, "history", false, "list history")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "command name")
	rootCmd.Flags().StringVarP(&exec, "exec", "e", "", `define the command (--exec "echo 'Hello World!'")`)
	rootCmd.Flags().StringSliceVarP(&alias, "alias", "a", []string{}, `define alias (--alias 'a1,a2,a3')`)
}
