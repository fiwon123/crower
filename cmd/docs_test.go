package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// The help text is the canonical documentation, but a user reading the wiki
// cannot see it. A command that ships without a section in the reference is a
// command they cannot find, so the reference is checked against the command
// tree rather than maintained by hope.
//
// This is the gate that keeps the wiki and the binary telling the same story.

// wikiPage reads one page from the wiki source.
func wikiPage(t *testing.T, name string) string {
	t.Helper()

	body, err := os.ReadFile(filepath.Join("..", "docs", "wiki", name))
	if err != nil {
		t.Fatalf("read docs/wiki/%s: %v", name, err)
	}

	return string(body)
}

// generatedCommands are the commands cobra adds at run time. They are not part
// of crower's own surface and are not documented. The check is on the whole
// path: cobra adds `completion` with four children of its own.
func generatedPath(path []string) bool {
	for _, name := range path {
		if name == "help" || name == "completion" {
			return true
		}
	}

	return false
}

func TestEveryCommandHasAReferenceSection(t *testing.T) {
	reference := wikiPage(t, "Command-Reference.md")

	var walk func(command *cobra.Command, path []string)

	walk = func(command *cobra.Command, path []string) {
		if len(path) > 0 && !generatedPath(path) {
			title := "crower " + strings.Join(path, " ")

			if !hasHeading(reference, title) {
				t.Errorf("docs/wiki/Command-Reference.md has no %q section", title)
			}
		}

		for _, child := range command.Commands() {
			walk(child, append(append([]string{}, path...), child.Name()))
		}
	}

	walk(rootCmd, nil)
}

// hasHeading reports whether the document has a Markdown heading that ends with
// title, so "### crower create" matches "crower create" wherever it sits.
func hasHeading(document, title string) bool {
	for _, line := range strings.Split(document, "\n") {
		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasSuffix(line, title) {
			return true
		}
	}

	return false
}

// Every top-level command has to be reachable from the sidebar: it is the only
// navigation the wiki has.
func TestEveryCommandIsInTheSidebar(t *testing.T) {
	sidebar := wikiPage(t, "_Sidebar.md")

	for _, child := range rootCmd.Commands() {
		if child.Hidden || generatedPath([]string{child.Name()}) {
			continue
		}

		if !strings.Contains(sidebar, child.Name()) {
			t.Errorf("docs/wiki/_Sidebar.md does not mention %q", child.Name())
		}
	}
}
