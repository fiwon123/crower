# Release notes template

The style used for `v0.8.0`. Use it for every release; keep the result under
~60 lines. Fill the placeholders, keep the section order.

Rules:

- **Title:** `Crower vX.Y.Z` — short, no emoji, no tagline.
- **Emoji:** only at the start of a section heading, and on the platform rows of
  the download table. Nothing else.
- **Tone:** professional, present tense, one idea per line. No exclamation
  marks, no "we are excited", no thanks paragraph.
- **Detail:** 4–6 highlights, a short fixes list, then the operational sections.
  The full PR list is a link, never pasted into the body.
- **Facts:** asset names, the commit and the checksum file must be copied from
  the actual release, not from memory.

---

```markdown
A terminal file manager and command-alias manager with undoable history and per-item notes.

## ✨ Highlights

- **<theme>.** One sentence on what changed and why it matters to the user.
- **<theme>.** …
- **<theme>.** …

## 🛠 Fixes

- <user-visible fix>.
- <user-visible fix>.

## 📦 Downloads

| Platform | CLI | GUI | Fyne |
| --- | --- | --- | --- |
| 🐧 Linux amd64 | `crower_vX.Y.Z_linux_amd64.tar.gz` | `crower_gui_vX.Y.Z_linux_amd64.tar.gz` | `crower_fyne_vX.Y.Z_linux_amd64.tar.gz` |
| 🐧 Linux arm64 | `crower_vX.Y.Z_linux_arm64.tar.gz` | `crower_gui_vX.Y.Z_linux_arm64.tar.gz` | — |
| 🪟 Windows amd64 | `crower_vX.Y.Z_windows_amd64.zip` | `crower_gui_vX.Y.Z_windows_amd64.zip` | `crower_fyne_vX.Y.Z_windows_amd64.zip` |
| 🪟 Windows arm64 | `crower_vX.Y.Z_windows_arm64.zip` | `crower_gui_vX.Y.Z_windows_arm64.zip` | — |
| 🍎 macOS Intel | `crower_vX.Y.Z_darwin_amd64.tar.gz` | `crower_gui_vX.Y.Z_darwin_amd64.tar.gz` | — |
| 🍎 macOS Apple Silicon | `crower_vX.Y.Z_darwin_arm64.tar.gz` | `crower_gui_vX.Y.Z_darwin_arm64.tar.gz` | `crower_fyne_vX.Y.Z_darwin_arm64.tar.gz` |

Each archive includes `crower`, the `cr` alias, README and licenses. The CLI archives also include the `add_to_path` installer. The Fyne archives are a separate CGO build for Linux/Windows amd64 and macOS arm64 and add `crower fyne`; the other archives keep the cgo-free desktop.

## ⚙️ Install

```sh
tar -xzf crower_vX.Y.Z_linux_amd64.tar.gz   # .zip on Windows
./add_to_path.sh                            # add_to_path.bat on Windows
crower --version
```

`crower gui` opens the terminal menu; `crower desktop` opens the cgo-free window;
`crower fyne` opens the Fyne window in the Fyne archive.

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
- Windows binaries are cross-compiled unless a Windows runner executed them.

## 📝 Changelog

[`vPREV...vX.Y.Z`](https://github.com/fiwon123/crower/compare/vPREV...vX.Y.Z) — one sentence on the main theme.
```
