// Command releasehighlights turns the merged pull requests of a release window
// into the Highlights and Fixes sections of the release body.
//
// It reads the JSON printed by `gh pr list --json
// number,title,url,labels,mergedAt` on stdin. The workflow resolves the
// window's commit dates, so the shaping stays testable without a network call.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// label is the part of a pull request's label the sectioning reads.
type label struct {
	Name string `json:"name"`
}

// pullRequest is the part of `gh pr list` output the sectioning reads.
type pullRequest struct {
	Number   int       `json:"number"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Labels   []label   `json:"labels"`
	MergedAt time.Time `json:"mergedAt"`
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: releasehighlights <since> <until> < prs.json")
		os.Exit(2)
	}

	since, err := time.Parse(time.RFC3339, os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "releasehighlights: since:", err)
		os.Exit(2)
	}

	until, err := time.Parse(time.RFC3339, os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "releasehighlights: until:", err)
		os.Exit(2)
	}

	var prs []pullRequest
	if err := json.NewDecoder(os.Stdin).Decode(&prs); err != nil {
		fmt.Fprintln(os.Stderr, "releasehighlights: decode:", err)
		os.Exit(1)
	}

	fmt.Print(render(prs, since, until))
}

// render formats the sections. It is separate from main so the shaping can be
// tested with fixture JSON and a fixed clock.
func render(prs []pullRequest, since, until time.Time) string {
	sort.Slice(prs, func(i, j int) bool { return prs[i].Number < prs[j].Number })

	var highlights, fixes []string
	for _, pr := range prs {
		if !pr.MergedAt.After(since) || pr.MergedAt.After(until) {
			continue
		}

		line := fmt.Sprintf("- **%s** — [#%d](%s)", cleanTitle(pr.Title), pr.Number, pr.URL)
		switch {
		case hasLabel(pr, "enhancement"):
			highlights = append(highlights, line)
		case hasLabel(pr, "bug"):
			fixes = append(fixes, line)
		}
	}

	var body strings.Builder
	writeSection(&body, "✨ Highlights", highlights)
	writeSection(&body, "🛠 Fixes", fixes)
	return body.String()
}

func writeSection(body *strings.Builder, heading string, lines []string) {
	if len(lines) == 0 {
		return
	}

	fmt.Fprintf(body, "## %s\n\n%s\n\n", heading, strings.Join(lines, "\n"))
}

func hasLabel(pr pullRequest, name string) bool {
	for _, l := range pr.Labels {
		if l.Name == name {
			return true
		}
	}
	return false
}

// conventionalPrefix matches the commit-type prefixes the repository uses in
// pull request titles, such as `feat: ` and `fix(api): `.
var conventionalPrefix = regexp.MustCompile(`^[a-z]+(\([^)]*\))?: `)

// cleanTitle drops a conventional-commit prefix and capitalises the first
// letter, so a PR title reads as a sentence in the release body.
func cleanTitle(title string) string {
	title = conventionalPrefix.ReplaceAllString(title, "")
	if title == "" {
		return title
	}

	runes := []rune(title)
	if runes[0] >= 'a' && runes[0] <= 'z' {
		runes[0] -= 'a' - 'A'
	}
	return string(runes)
}
