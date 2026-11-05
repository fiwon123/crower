// Command bodycheck verifies a release body against the house style in
// docs/release-notes-template.md.
//
// The release workflow runs it on the generated notes before creating the
// draft, so a generator change that drops a section or reintroduces a badge
// fails the dispatch instead of shipping a body that reads differently from
// the rest of the release page.
package main

import (
	"fmt"
	"os"
	"strings"
)

// requiredSections must appear, in this order, in every release body.
var requiredSections = []string{
	"## 📦 Downloads",
	"## ⚙️ Install",
	"## ⚡ Upgrade",
	"## 🔐 Verify",
	"## 📝 Changelog",
}

// optionalSections are written from the release window's pull requests. When
// present they must precede Downloads; a fix-only release omits Highlights.
var optionalSections = []string{
	"## ✨ Highlights",
	"## 🛠 Fixes",
}

// downloadRows pin the platform icon table.
var downloadRows = []string{
	"| Platform | CLI | GUI |",
	"| 🐧 Linux amd64 |",
	"| 🪟 Windows amd64 |",
	"| 🍎 macOS Intel |",
	"| 🍎 macOS Apple Silicon |",
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: bodycheck <notes.md>")
		os.Exit(2)
	}

	body, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "bodycheck:", err)
		os.Exit(1)
	}

	problems := check(string(body))
	if len(problems) > 0 {
		fmt.Fprintf(os.Stderr, "%s does not follow docs/release-notes-template.md:\n", os.Args[1])
		for _, problem := range problems {
			fmt.Fprintf(os.Stderr, "  - %s\n", problem)
		}
		os.Exit(1)
	}

	fmt.Printf("%s: release body OK\n", os.Args[1])
}

// check reports what the body is missing or has out of place.
func check(body string) []string {
	var problems []string

	if !strings.Contains(body, "A terminal file manager and command-alias manager") {
		problems = append(problems, "missing the opening sentence")
	}

	if strings.Contains(body, "img.shields.io") {
		problems = append(problems, "shields.io badges are not part of the house style")
	}

	if strings.Contains(body, "## Crower ") {
		problems = append(problems, "the body must not repeat the release title")
	}

	last := -1
	for _, heading := range requiredSections {
		at := strings.Index(body, heading)
		if at < 0 {
			problems = append(problems, fmt.Sprintf("missing %q", heading))
			continue
		}
		if at <= last {
			problems = append(problems, fmt.Sprintf("%q is out of order", heading))
		}
		last = at
	}

	downloads := strings.Index(body, requiredSections[0])
	previous := -1
	for _, heading := range optionalSections {
		at := strings.Index(body, heading)
		if at < 0 {
			continue
		}
		if downloads >= 0 && at > downloads {
			problems = append(problems, fmt.Sprintf("%q must precede %q", heading, requiredSections[0]))
		}
		if at <= previous {
			problems = append(problems, fmt.Sprintf("%q is out of order", heading))
		}
		previous = at
	}

	for _, row := range downloadRows {
		if !strings.Contains(body, row) {
			problems = append(problems, fmt.Sprintf("missing the download row %q", row))
		}
	}

	if !strings.Contains(body, "vcs.modified=false") {
		problems = append(problems, "missing the provenance line")
	}

	return problems
}
