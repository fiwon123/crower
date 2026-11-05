# Release notes template

The house style, first used by `v0.8.0`. Every release body follows it; keep the
result under ~60 lines and keep the section order.

## Rules

- **Title.** The release title is `Crower vX.Y.Z` — short, no emoji. The body
  does not repeat it.
- **Colour.** GitHub does not render custom text colour. Sections use emoji
  headings and the download table uses platform icons; never shields.io badges.
- **Icons.** 🐧 Linux, 🪟 Windows, 🍎 macOS. Section headings: ✨ Highlights,
  🛠 Fixes, 📦 Downloads, ⚙️ Install, ⚡ Upgrade, 🔐 Verify, ⚠️ Known
  limitations, 📝 Changelog.
- **Tone.** Professional, present tense, one idea per line. No exclamation
  marks, no "we are excited", no thanks paragraph.
- **Detail.** 3–6 highlights and a short fixes list, then the operational
  sections. The full PR list is a compare link, never pasted into the body.
- **Facts.** Asset names, the commit and the checksum file are copied from the
  published release, not from memory.
- **One archive per platform.** `crower_vX.Y.Z_<os>_<arch>` contains `crower`,
  the `cr` alias and `crower-gui` (the Fyne window), plus README, licenses and
  the `add_to_path` installer. There is no separate GUI archive.
- **Draft first.** The workflow creates a draft with Highlights, Fixes and the
  operational sections already in place. Highlights and Fixes come from the
  PRs merged into `dev` (labels `enhancement` and `bug`); rewrite a line if a
  title needs a sentence, then publish. `internal/tools/bodycheck` refuses a
  body that drifts from this order, so the dispatch fails before the draft
  exists.

## Template

````markdown
A terminal file manager and command-alias manager with undoable history and per-item notes.

## ✨ Highlights

- **<theme>.** One sentence on what changed and why it matters to the user.
- **<theme>.** …
- **<theme>.** …

## 🛠 Fixes

- <user-visible fix>.
- <user-visible fix>.

## 📦 Downloads

| Platform | Archive |
| --- | --- |
| 🐧 Linux amd64 | `crower_vX.Y.Z_linux_amd64.tar.gz` |
| 🐧 Linux arm64 | `crower_vX.Y.Z_linux_arm64.tar.gz` |
| 🪟 Windows amd64 | `crower_vX.Y.Z_windows_amd64.zip` |
| 🪟 Windows arm64 | `crower_vX.Y.Z_windows_arm64.zip` |
| 🍎 macOS Intel | `crower_vX.Y.Z_darwin_amd64.tar.gz` |
| 🍎 macOS Apple Silicon | `crower_vX.Y.Z_darwin_arm64.tar.gz` |

Every archive contains `crower`, the `cr` alias and `crower-gui`, plus README, licenses and the `add_to_path` installer. The CLI is cgo-free; `crower-gui` is the Fyne window over the same configuration.

## ⚙️ Install

```sh
tar -xzf crower_vX.Y.Z_linux_amd64.tar.gz   # .zip on Windows
./add_to_path.sh                            # add_to_path.bat on Windows
crower --version
```

`crower gui` opens the terminal menu; `crower-gui` opens the Fyne window.

## ⚡ Upgrade

```sh
crower --upgrade
```

Verifies the archive against the release's `checksums.txt` before replacing the binary.

## 🔐 Verify

- SHA-256 for every asset: `checksums.txt`.
- Built from `<commit>` (tag `vX.Y.Z`), `vcs.modified=false`:

```sh
go version -m crower | grep vcs.revision
```

## ⚠️ Known limitations

- <open issue that affects this release, or omit the section entirely>.

## 📝 Changelog

[`vPREV...vX.Y.Z`](https://github.com/fiwon123/crower/compare/vPREV...vX.Y.Z) — one sentence on the main theme.
````
