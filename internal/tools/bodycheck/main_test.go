package main

import (
	"strings"
	"testing"
)

const validBody = `A terminal file manager and command-alias manager with undoable history and per-item notes.

## ✨ Highlights

- **A feature** — [#1](https://example/1)

## 🛠 Fixes

- **A fix** — [#2](https://example/2)

## 📦 Downloads

| Platform | CLI | GUI |
| --- | --- | --- |
| 🐧 Linux amd64 | ` + "`crower_v1_linux_amd64.tar.gz`" + ` | ` + "`crower_gui_v1_linux_amd64.tar.gz`" + ` |
| 🐧 Linux arm64 | ` + "`crower_v1_linux_arm64.tar.gz`" + ` | — |
| 🪟 Windows amd64 | ` + "`crower_v1_windows_amd64.zip`" + ` | ` + "`crower_gui_v1_windows_amd64.zip`" + ` |
| 🪟 Windows arm64 | ` + "`crower_v1_windows_arm64.zip`" + ` | — |
| 🍎 macOS Intel | ` + "`crower_v1_darwin_amd64.tar.gz`" + ` | ` + "`crower_gui_v1_darwin_amd64.tar.gz`" + ` |
| 🍎 macOS Apple Silicon | ` + "`crower_v1_darwin_arm64.tar.gz`" + ` | ` + "`crower_gui_v1_darwin_arm64.tar.gz`" + ` |

## ⚙️ Install

` + "```sh" + `
tar -xzf crower_v1_linux_amd64.tar.gz
` + "```" + `

## ⚡ Upgrade

` + "```sh" + `
crower --upgrade
` + "```" + `

## 🔐 Verify

- Built from ` + "`abc123`" + `, ` + "`vcs.modified=false`" + `.

## 📝 Changelog

[` + "`v0...v1`" + `](https://example/compare)
`

func TestCheckAcceptsHouseStyle(t *testing.T) {
	if problems := check(validBody); len(problems) > 0 {
		t.Fatalf("check() = %v, want no problems", problems)
	}
}

func TestCheckAcceptsFixOnlyBody(t *testing.T) {
	body := strings.Replace(validBody, "## ✨ Highlights\n\n- **A feature** — [#1](https://example/1)\n\n", "", 1)
	if problems := check(body); len(problems) > 0 {
		t.Fatalf("check() = %v, want no problems: a fix-only release omits Highlights", problems)
	}
}

func TestCheckReportsDrift(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{
			name: "missing opening sentence",
			body: strings.Replace(validBody, "A terminal file manager and command-alias manager with undoable history and per-item notes.", "", 1),
			want: "opening sentence",
		},
		{
			name: "shields badge",
			body: strings.Replace(validBody, "A terminal file manager", "![Version](https://img.shields.io/badge/v1)\n\nA terminal file manager", 1),
			want: "shields.io",
		},
		{
			name: "repeated title",
			body: strings.Replace(validBody, "A terminal file manager", "## Crower v1\n\nA terminal file manager", 1),
			want: "repeat the release title",
		},
		{
			name: "missing upgrade",
			body: strings.Replace(validBody, "## ⚡ Upgrade\n\n```sh\ncrower --upgrade\n```\n\n", "", 1),
			want: "missing \"## ⚡ Upgrade\"",
		},
		{
			name: "sections out of order",
			body: swap(validBody, "## ⚙️ Install", "## ⚡ Upgrade"),
			want: "out of order",
		},
		{
			name: "highlights after downloads",
			body: appendSection(strings.Replace(validBody, "## ✨ Highlights\n\n- **A feature** — [#1](https://example/1)\n\n", "", 1), "## ✨ Highlights"),
			want: "must precede",
		},
		{
			name: "missing platform row",
			body: strings.Replace(validBody, "| 🍎 macOS Apple Silicon |", "| macOS |", 1),
			want: "download row",
		},
		{
			name: "missing provenance",
			body: strings.Replace(validBody, ", `vcs.modified=false`", "", 1),
			want: "provenance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			problems := check(tt.body)
			if len(problems) == 0 {
				t.Fatalf("check() = no problems, want one mentioning %q", tt.want)
			}
			if !strings.Contains(strings.Join(problems, "\n"), tt.want) {
				t.Errorf("check() = %v, want a problem mentioning %q", problems, tt.want)
			}
		})
	}
}

// swap exchanges two heading strings, so the first section now carries the
// heading of the second and the checker sees them out of order.
func swap(body, a, b string) string {
	body = strings.Replace(body, a, "\x00", 1)
	body = strings.Replace(body, b, a, 1)
	return strings.Replace(body, "\x00", b, 1)
}

func appendSection(body, heading string) string {
	return body + "\n" + heading + "\n\n- **Late** — [#9](https://example/9)\n"
}
