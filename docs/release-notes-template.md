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
- **Assets.** `crower_vX.Y.Z_<os>_<arch>` is the cgo-free CLI; it ships
  `crower`, the `cr` alias, README, licenses and `add_to_path`.
  `crower_gui_vX.Y.Z_<os>_<arch>` is the Fyne app `crower-gui`, for Linux and
  Windows amd64 plus macOS Apple Silicon and Intel. The other two platforms get
  the CLI only.
- **Draft first.** The workflow creates a draft with the operational sections.
  Add Highlights and Fixes while it is a draft, then publish.

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

| Platform | CLI | GUI |
| --- | --- | --- |
| 🐧 Linux amd64 | `crower_vX.Y.Z_linux_amd64.tar.gz` | `crower_gui_vX.Y.Z_linux_amd64.tar.gz` |
| 🐧 Linux arm64 | `crower_vX.Y.Z_linux_arm64.tar.gz` | — |
| 🪟 Windows amd64 | `crower_vX.Y.Z_windows_amd64.zip` | `crower_gui_vX.Y.Z_windows_amd64.zip` |
| 🪟 Windows arm64 | `crower_vX.Y.Z_windows_arm64.zip` | — |
| 🍎 macOS Intel | `crower_vX.Y.Z_darwin_amd64.tar.gz` | `crower_gui_vX.Y.Z_darwin_amd64.tar.gz` |
| 🍎 macOS Apple Silicon | `crower_vX.Y.Z_darwin_arm64.tar.gz` | `crower_gui_vX.Y.Z_darwin_arm64.tar.gz` |

The CLI archive is cgo-free and ships `crower`, the `cr` alias, README, licenses and the `add_to_path` installer. The GUI archive ships `crower-gui`, the Fyne window over the same configuration, with README and licenses.

## ⚙️ Install

```sh
tar -xzf crower_vX.Y.Z_linux_amd64.tar.gz   # .zip on Windows
./add_to_path.sh                            # add_to_path.bat on Windows
crower --version
```

`crower gui` opens the terminal menu; `crower-gui` from the GUI archive opens the Fyne window.

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
