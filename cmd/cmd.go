package cmd

import (
	"fmt"
	"os"

	"github.com/fiwon123/crower/internal/application/checkrelease"
	"github.com/fiwon123/crower/internal/legacy/executedom"
	"github.com/fiwon123/crower/internal/legacy/listdom"
	"github.com/fiwon123/crower/internal/legacy/upgradedom"

	"github.com/spf13/cobra"
)

var versionFlag bool
var checkNewVersion bool
var upgradeFlag bool

// Shared flags, parsed by more than one command.
var executeFlag bool
var fileFlag bool
var folderFlag bool
var sysPathFlag bool
var systemFlag bool
var createFlag bool
var updateFlag bool
var processFlag bool

// Version is set at build time with -X github.com/fiwon123/crower/cmd.Version.
var Version = "vx.x.x"

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "crower",
	Short: "A terminal file manager and command-alias manager with undoable history.",
	Long: `crower keeps the commands you run often and the files you touch, and it can
undo what it did.

Run a managed command:
	crower greet
	crower execute greet

Manage the command set:
	crower create greet "echo hello" --alias hi
	crower list
	crower update greet --exec "echo hi"
	crower delete greet

Work with files:
	crower create --file notes.txt
	crower copy notes.txt ~/backup
	crower open --folder .

Undo what happened:
	crower list --history
	crower revert 1

Results go to stdout and diagnostics go to stderr, so "crower list | wc -l"
counts rows. A failure exits 1, and a prompt with no input available fails
instead of waiting.

Every command has its own help: crower <command> --help.`,
	Aliases: []string{"cr"},
	Args:    cobra.ArbitraryArgs,
	// Cobra must not print the error itself: Execute prints it once, to the
	// diagnostics stream, and turns it into the exit status.
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// --version answers without reading the configuration, so the answer
		// works on a machine where crower has never run.
		if versionFlag {
			fmt.Fprintln(streams.Out, Version)
			return nil
		}

		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		switch {
		case checkNewVersion:
			session := newLegacySession(cfg, listdom.NewHandler(cfg).CommandsTable)
			return checkrelease.New(latestTag{}, session, session, session).Run(Version)

		case upgradeFlag:
			return upgradedom.NewCore(cfg).UpgradeApp(Version)

		case len(args) == 0:
			return cmd.Help()
		}

		listHandler := listdom.NewHandler(cfg)

		return executedom.NewCore(cfg, listHandler).ExecuteCommand(args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
//
// A returned error becomes exit status 1 and one line on the diagnostics
// stream. That is the contract a script depends on: the exit status says
// whether the work happened, and the message says why it did not.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(streams.Err, "error:", err)
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

	// Persistent flags
	rootCmd.PersistentFlags().String(flagConfig, "", "configuration file path (default: $HOME/crower/crower.yaml)")
	rootCmd.PersistentFlags().Bool(flagNoInput, false, "never prompt; fail instead of waiting for input")

	// Flags
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "print the current version")
	rootCmd.Flags().BoolVar(&upgradeFlag, "upgrade", false, "upgrade to the latest version")
	rootCmd.Flags().BoolVar(&checkNewVersion, "check", false, "check whether a newer version exists")
}
