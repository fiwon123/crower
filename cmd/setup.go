package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fiwon123/crower/internal/adapters/storage"

	"github.com/spf13/cobra"
)

var setupForce bool

// setupCmd manages the named command sets.
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Manage named command setups with their own history.",
	Long: `Manage named command setups. Each setup keeps its own commands and its own
history under ~/crower/setups/<name>/, so a revert in one never touches another.

Select one for a run with --setup, CROWER_SETUP, or "crower setup use". The
default setup keeps the original ~/crower/crower.yaml and needs no migration.`,
	Example: "  crower setup list\n  crower setup new streamer\n  crower setup use streamer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var setupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List the named setups.",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		active := activeSetup()

		entries, err := os.ReadDir(filepath.Join(home, "crower", "setups"))
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		fmt.Fprintln(streams.Out, "default")
		if active == "" || active == "default" {
			fmt.Fprintln(streams.Out, "  (active)")
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			fmt.Fprintln(streams.Out, entry.Name())

			if entry.Name() == active {
				fmt.Fprintln(streams.Out, "  (active)")
			}
		}

		return nil
	},
}

var setupNewCmd = &cobra.Command{
	Use:   "new <name>",
	Short: "Create an empty named setup.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		path, err := createSetup(home, args[0])
		if err != nil {
			return err
		}

		fmt.Fprintf(streams.Out, "created setup %s (%s)\n", args[0], path)

		return nil
	},
}

var setupUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Make a setup the default for future runs.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		if args[0] != "default" {
			path, err := setupConfigPath(home, args[0])
			if err != nil {
				return err
			}

			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("setup %q does not exist; create it with crower setup new %s", args[0], args[0])
			}
		}

		marker := filepath.Join(home, "crower", "setup")

		if err := os.MkdirAll(filepath.Dir(marker), 0o755); err != nil {
			return err
		}

		if err := os.WriteFile(marker, []byte(args[0]+"\n"), 0o644); err != nil {
			return err
		}

		fmt.Fprintf(streams.Out, "default setup is now %s\n", args[0])

		return nil
	},
}

var setupShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the active setup and where it lives.",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := activeSetup()
		if name == "" {
			name = "default"
		}

		path, err := resolveConfigPath("", name)
		if err != nil {
			return err
		}

		commands := 0

		if file, _, err := storage.NewCommandFile(path).Load(); err == nil {
			commands = len(file)
		}

		fmt.Fprintf(streams.Out, "Setup:    %s\n", name)
		fmt.Fprintf(streams.Out, "Config:   %s\n", path)
		fmt.Fprintf(streams.Out, "Commands: %d\n", commands)

		return nil
	},
}

var setupDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a named setup and its history.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		active := activeSetup()
		if setupForce {
			active = ""
		}

		if err := deleteSetup(home, args[0], active); err != nil {
			return err
		}

		fmt.Fprintf(streams.Out, "deleted setup %s\n", args[0])

		return nil
	},
}

// createSetup makes an empty named setup and returns the config path it wrote.
// The CLI and the window both call it, so `setup new` and the window's New…
// cannot drift apart.
func createSetup(home, name string) (string, error) {
	path, err := setupConfigPath(home, name)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("setup %q already exists", name)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}

	if err := os.WriteFile(path, nil, 0o644); err != nil {
		return "", err
	}

	return path, nil
}

// deleteSetup removes a named setup and its folder. It refuses the default setup,
// which is the original installation, and the active one -- a running session
// would otherwise point at a file that no longer exists.
func deleteSetup(home, name, active string) error {
	if name == "default" {
		return fmt.Errorf("the default setup is not a folder; delete its commands with crower reset")
	}

	if name == active {
		return fmt.Errorf("setup %q is active; switch first or pass --force", name)
	}

	path, err := setupConfigPath(home, name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("setup %q does not exist", name)
	}

	return os.RemoveAll(filepath.Dir(path))
}

// activeSetup resolves the selected setup the way a run does: --setup, then
// CROWER_SETUP, then the persisted default, then "default".
func activeSetup() string {
	if value, err := rootCmd.PersistentFlags().GetString(flagSetup); err == nil && value != "" {
		return value
	}

	if value := strings.TrimSpace(os.Getenv("CROWER_SETUP")); value != "" {
		return value
	}

	return persistedSetup()
}

func init() {
	setupDeleteCmd.Flags().BoolVar(&setupForce, "force", false, "delete the active setup too")

	setupCmd.AddCommand(setupListCmd, setupNewCmd, setupUseCmd, setupShowCmd, setupDeleteCmd)
	rootCmd.AddCommand(setupCmd)
}
