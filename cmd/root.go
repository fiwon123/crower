package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/pkg/utils"
	"github.com/spf13/cobra"
)

var cfgFilePath string
var createOp bool
var listOp bool
var resetOp bool
var index int
var name string
var exec string
var alias []string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crower",
	Short: "A dev tool that manages system commands to help developers in their daily workflow.",
	Long: `A dev tool that manages system commands by executing commands via custom aliases and
managing it with useful operations like add, edit, remove, list and more.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:s
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) > 0 {
			fmt.Println("args", args)
			name = args[0]
			alias = append(alias, args[0])
		}
		fmt.Println("cfg", cfgFilePath)

		utils.CreateTomlIfNotExists(cfgFilePath)
		fmt.Println("cfgfilepath: ", cfgFilePath)

		commandsMap := utils.ReadToml(cfgFilePath)
		aliasMap := core.GetAliasMap(commandsMap)

		var op data.CommandOperation
		if createOp {
			op = data.Create
		} else if listOp {
			op = data.List
		} else if resetOp {
			op = data.Reset
		} else {
			op = data.Execute
		}

		core.HandlePayload(
			data.Payload{
				Op:      op,
				Command: *data.NewCommand(name, alias, exec),
			},
			&data.App{
				CfgFilePath: cfgFilePath,
				AliasMap:    aliasMap,
				CommandsMap: commandsMap,
			})

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

	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Error could not get user home directory, ", err)
	}

	rootCmd.Flags().StringVar(&cfgFilePath, "config", filepath.Join(homePath, "crower.yaml"), "config file (default is $HOME/.crower.yaml)")
	rootCmd.Flags().IntVarP(&index, "index", "i", 0, "command index")
	rootCmd.Flags().BoolVar(&createOp, "create", false, "create a command")
	rootCmd.Flags().BoolVar(&listOp, "list", false, "list all commands")
	rootCmd.Flags().BoolVar(&resetOp, "reset", false, "reset all commands")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "command name")
	rootCmd.Flags().StringVarP(&exec, "exec", "e", "", `define the command (example "echo 'Hello World!'")`)
	rootCmd.Flags().StringSliceVarP(&alias, "alias", "a", []string{}, `define alias (example "--alias 'a1,a2,a3'")`)
}
