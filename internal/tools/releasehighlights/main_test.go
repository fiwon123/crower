package main

import (
	"testing"
	"time"
)

func at(day int) time.Time {
	return time.Date(2026, 11, day, 12, 0, 0, 0, time.UTC)
}

func TestRenderSections(t *testing.T) {
	prs := []pullRequest{
		{Number: 3, Title: "feat: make the window editable", URL: "https://example/3", Labels: []label{{Name: "enhancement"}}, MergedAt: at(2)},
		{Number: 1, Title: "fix: package the archive", URL: "https://example/1", Labels: []label{{Name: "bug"}}, MergedAt: at(2)},
		{Number: 2, Title: "feat: too early", URL: "https://example/2", Labels: []label{{Name: "enhancement"}}, MergedAt: at(1)},
		{Number: 4, Title: "chore: maintenance", URL: "https://example/4", Labels: []label{{Name: "chore"}}, MergedAt: at(2)},
	}

	got := render(prs, at(1), at(3))

	want := "## ✨ Highlights\n\n" +
		"- **Make the window editable** — [#3](https://example/3)\n\n" +
		"## 🛠 Fixes\n\n" +
		"- **Package the archive** — [#1](https://example/1)\n\n"
	if got != want {
		t.Errorf("render() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderOmitsEmptySections(t *testing.T) {
	prs := []pullRequest{
		{Number: 1, Title: "fix: only a fix", URL: "https://example/1", Labels: []label{{Name: "bug"}}, MergedAt: at(2)},
	}

	got := render(prs, at(1), at(3))

	want := "## 🛠 Fixes\n\n- **Only a fix** — [#1](https://example/1)\n\n"
	if got != want {
		t.Errorf("render() = %q, want %q", got, want)
	}
}

func TestRenderNothingInWindow(t *testing.T) {
	prs := []pullRequest{
		{Number: 1, Title: "feat: out of window", URL: "https://example/1", Labels: []label{{Name: "enhancement"}}, MergedAt: at(1)},
	}

	if got := render(prs, at(2), at(3)); got != "" {
		t.Errorf("render() = %q, want empty", got)
	}
}

func TestRenderWindowBoundaries(t *testing.T) {
	prs := []pullRequest{
		{Number: 1, Title: "feat: at since", URL: "https://example/1", Labels: []label{{Name: "enhancement"}}, MergedAt: at(2)},
		{Number: 2, Title: "feat: at until", URL: "https://example/2", Labels: []label{{Name: "enhancement"}}, MergedAt: at(3)},
	}

	got := render(prs, at(2), at(3))

	// The window is (since, until]: a PR merged exactly at since belongs to the
	// previous release, one merged exactly at until belongs to this one.
	want := "## ✨ Highlights\n\n- **At until** — [#2](https://example/2)\n\n"
	if got != want {
		t.Errorf("render() = %q, want %q", got, want)
	}
}

func TestCleanTitle(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"feat: make it work", "Make it work"},
		{"fix(api): return an error", "Return an error"},
		{"docs: explain the flag", "Explain the flag"},
		{"refactor: rename the package", "Rename the package"},
		{"Already capitalised", "Already capitalised"},
		{"no prefix", "No prefix"},
		{"", ""},
	}

	for _, tt := range tests {
		if got := cleanTitle(tt.title); got != tt.want {
			t.Errorf("cleanTitle(%q) = %q, want %q", tt.title, got, tt.want)
		}
	}
}
