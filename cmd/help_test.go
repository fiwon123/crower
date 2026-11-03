package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// The help text is the only documentation a user has when they are already in
// the terminal, and the refactor makes it the source of truth for the wiki. A
// command that cannot render its own help, or has no one-line description, is
// undocumented no matter what the wiki says, so the surface is pinned here
// before anything moves.

// helpOf renders the help of one command path and returns what reached the
// streams. Cobra writes help to cmd.OutOrStdout, which is os.Stdout unless a
// test redirects it -- redirecting per call keeps it out of the test log.
func helpOf(t *testing.T, path ...string) string {
	t.Helper()

	resetFlags()

	buf := &bytes.Buffer{}

	previousOut, previousErr := rootCmd.OutOrStdout(), rootCmd.ErrOrStderr()
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	t.Cleanup(func() {
		rootCmd.SetOut(previousOut)
		rootCmd.SetErr(previousErr)
	})

	rootCmd.SetArgs(append(append([]string{}, path...), "--help"))
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("%v --help returned an error: %v", append([]string{"crower"}, path...), err)
	}

	return buf.String()
}

// walkCommands visits every command in the tree with the path that reaches it.
func walkCommands(c *cobra.Command, path []string, visit func(path []string, c *cobra.Command)) {
	visit(path, c)

	for _, child := range c.Commands() {
		walkCommands(child, append(append([]string{}, path...), child.Name()), visit)
	}
}

func commandPath(c *cobra.Command, path []string) string {
	return strings.TrimSpace("crower " + strings.Join(path, " "))
}

func TestEveryCommandRendersHelp(t *testing.T) {
	walkCommands(rootCmd, nil, func(path []string, c *cobra.Command) {
		t.Run(commandPath(c, path), func(t *testing.T) {
			out := helpOf(t, path...)

			if !strings.Contains(out, "Usage:") {
				t.Errorf("help has no usage section:\n%s", out)
			}

			if !strings.Contains(out, c.Name()) {
				t.Errorf("help does not name the command %q:\n%s", c.Name(), out)
			}
		})
	})
}

// Short is what a user sees in `crower --help` and in every parent listing. An
// empty one makes a command invisible in the one place it is discoverable.
func TestEveryCommandHasAShortDescription(t *testing.T) {
	walkCommands(rootCmd, nil, func(path []string, c *cobra.Command) {
		if strings.TrimSpace(c.Short) == "" {
			t.Errorf("%s has no Short description", commandPath(c, path))
		}
	})
}

// A command's Use line starts with the word that dispatches to it. When the two
// disagree the documented invocation does not run the command it names.
func TestUseLineMatchesTheCommandName(t *testing.T) {
	walkCommands(rootCmd, nil, func(path []string, c *cobra.Command) {
		if path == nil {
			return
		}

		first, _, _ := strings.Cut(strings.TrimSpace(c.Use), " ")
		if first != c.Name() {
			t.Errorf("%s: Use %q does not start with %q", commandPath(c, path), c.Use, c.Name())
		}
	})
}

// Every flag is printed in the help. A flag with no usage string is a flag the
// user cannot discover a purpose for.
func TestEveryFlagHasAUsage(t *testing.T) {
	walkCommands(rootCmd, nil, func(path []string, c *cobra.Command) {
		c.Flags().VisitAll(func(flag *pflag.Flag) {
			if strings.TrimSpace(flag.Usage) == "" {
				t.Errorf("%s: flag --%s has no usage text", commandPath(c, path), flag.Name)
			}
		})
	})
}

// The root help is where a user looks for the command list, so every top-level
// command has to appear in it.
func TestRootHelpListsEveryCommand(t *testing.T) {
	out := helpOf(t)

	for _, child := range rootCmd.Commands() {
		if child.Hidden {
			continue
		}

		if !strings.Contains(out, child.Name()) {
			t.Errorf("root help does not list %q:\n%s", child.Name(), out)
		}
	}
}

// The short alias is part of the interface: a user types `cr` all day. It is
// declared on the root command, so it must survive any restructure of the tree.
func TestRootKeepsTheShortAlias(t *testing.T) {
	if !contains(rootCmd.Aliases, "cr") {
		t.Errorf("root aliases = %v, want cr", rootCmd.Aliases)
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}

	return false
}
