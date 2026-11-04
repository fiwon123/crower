package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var initApply bool

// initCmd explains, or applies, the one environment change crower needs: the
// install directory on PATH.
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Show (or apply) the PATH setup for your shell.",
	Long: `Show the line your shell needs so "crower" and "cr" are on PATH, and how to
enable completion. With --apply, append the line to the shell's startup file if
it is not there already.

	crower init
	crower init --apply`,
	Example: "  crower init\n  crower init --apply",
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		shell := shellName(os.Getenv("SHELL"))
		target, line := initTarget(shell, home)

		if target == "" {
			fmt.Fprintln(streams.Out, "Could not recognise your shell.")
			fmt.Fprintln(streams.Out, `Add the crower install directory to PATH yourself, then run "crower completion <shell>" for completion.`)
			return nil
		}

		if !initApply {
			fmt.Fprintf(streams.Out, "Shell:   %s\n", shell)
			fmt.Fprintf(streams.Out, "Startup: %s\n", target)
			fmt.Fprintf(streams.Out, "Add:     %s\n", line)
			fmt.Fprintln(streams.Out)
			fmt.Fprintln(streams.Out, `Or run "crower init --apply" to add it now.`)
			fmt.Fprintf(streams.Out, "Completion: crower completion %s\n", shell)
			return nil
		}

		changed, err := applyInit(target, line)
		if err != nil {
			return err
		}

		if changed {
			fmt.Fprintf(streams.Out, "Updated %s. Restart the shell (or source it) to pick it up.\n", target)
		} else {
			fmt.Fprintf(streams.Out, "%s already sets PATH. Nothing changed.\n", target)
		}

		fmt.Fprintf(streams.Out, "Completion: crower completion %s\n", shell)

		return nil
	},
}

// shellName turns $SHELL into the name the commands and the startup files use.
func shellName(shell string) string {
	base := strings.ToLower(filepath.Base(shell))

	switch {
	case strings.Contains(base, "zsh"):
		return "zsh"
	case strings.Contains(base, "bash"):
		return "bash"
	case strings.Contains(base, "fish"):
		return "fish"
	case strings.Contains(base, "pwsh"), strings.Contains(base, "powershell"):
		return "powershell"
	default:
		return ""
	}
}

// initTarget is the startup file and the line that puts the install directory
// on PATH for one shell.
func initTarget(shell, home string) (string, string) {
	switch shell {
	case "bash":
		return filepath.Join(home, ".bashrc"), `export PATH="$HOME/.local/bin:$PATH"`
	case "zsh":
		return filepath.Join(home, ".zshrc"), `export PATH="$HOME/.local/bin:$PATH"`
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish"), `fish_add_path "$HOME/.local/bin"`
	default:
		return "", ""
	}
}

// applyInit appends the line unless the file already contains it.
//
// It reports whether the file changed: a second run is a no-op, because a
// startup file that grows on every init is noise.
func applyInit(path, line string) (bool, error) {
	body, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	if strings.Contains(string(body), line) {
		return false, nil
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return false, err
	}
	defer file.Close()

	if _, err := fmt.Fprintf(file, "\n%s\n", line); err != nil {
		return false, err
	}

	return true, nil
}

func init() {
	initCmd.Flags().BoolVar(&initApply, "apply", false, "append the line to the shell startup file")
	rootCmd.AddCommand(initCmd)
}
