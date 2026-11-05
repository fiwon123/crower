# Release notes template

The style used for `v0.12.0`. Use it for every release; keep the result under
~60 lines. Fill the placeholders, keep the section order.

Rules:

- **Title:** `Crower vX.Y.Z` — short, no tagline.
- **Colour:** GitHub does not render custom text colour. Use shields.io badges
  at the top and platform emoji in the download table; that is the colour.
- **Platform icons:** 🐧 Linux, 🪟 Windows, 🍎 macOS.
- **Tone:** professional, present tense, one idea per line. No exclamation
  marks, no "we are excited", no thanks paragraph.
- **Detail:** 4–6 highlights, a short fixes list, then the operational sections.
  The full PR list is a link, never pasted into the body.
- **Facts:** asset names, the commit and the checksum file must be copied from
  the actual release, not from memory.
- **Assets:** the CLI archive is cgo-free; the GUI archive is the Fyne app
  `crower-gui`. GUI exists for Linux/Windows amd64 and macOS arm64/Intel only.

---

```markdown
![Version](https://img.shields.io/badge/version-vX.Y.Z-2563eb?style=flat-square)
![Platforms](https://img.shields.io/badge/platforms-Linux%20%7C%20Windows%20%7C%20macOS-059669?style=flat-square)
![CLI](https://img.shields.io/badge/CLI-cgo--free-059669?style=flat-square)
![GUI](https://img.shields.io/badge/GUI-Fyne-7c3aed?style=flat-square)

## Crower vX.Y.Z

A terminal file manager and command-alias manager with undoable history and per-item notes.

### ✨ Highlights

- **<theme>.** One sentence on what changed and why it matters to the user.
- **<theme>.** …
- **<theme>.** …

### 🛠 Fixes

- <user-visible fix>.
- <user-visible fix>.

### 📦 Downloads

| Platform | CLI | GUI (Fyne) |
| :-- | :--: | :--: |
| 🐧 Linux amd64 | [crower_vX.Y.Z_linux_amd64.tar.gz](https://github.com/fiwon123/crower/releases/download/vX.Y.Z/crower_vX.Y.Z_linux_amd64.tar.gz) | [crower_gui_vX.Y.Z_linux_amd64.tar.gz](https://github.com/fiwon123/crower/releases/download/vX.Y.Z/crower_gui_vX.Y.Z_linux_amd64.tar.gz) |
| 🐧 Linux arm64 | [crower_vX.Y.Z_linux_arm64.tar.gz](…) | — |
| 🪟 Windows amd64 | [crower_vX.Y.Z_windows_amd64.zip](…) | [crower_gui_vX.Y.Z_windows_amd64.zip](…) |
| 🪟 Windows arm64 | [crower_vX.Y.Z_windows_arm64.zip](…) | — |
| 🍎 macOS Apple Silicon | [crower_vX.Y.Z_darwin_arm64.tar.gz](…) | [crower_gui_vX.Y.Z_darwin_arm64.tar.gz](…) |
| 🍎 macOS Intel | [crower_vX.Y.Z_darwin_amd64.tar.gz](…) | [crower_gui_vX.Y.Z_darwin_amd64.tar.gz](…) |

> **CLI** is cgo-free and has no window. **GUI** is the Fyne app `crower-gui`.

### ⚡ Install

```sh
tar -xzf crower_vX.Y.Z_linux_amd64.tar.gz   # .zip on Windows
./add_to_path.sh                            # add_to_path.bat on Windows
crower --version
```

Run `crower gui` for the terminal menu, or launch `crower-gui` from the GUI archive.

### ⚡ Upgrade

```sh
crower --upgrade
```

Verifies the archive against the release's `checksums.txt` before replacing the binary.

### 🔐 Verify

- SHA-256 for every asset: `checksums.txt`.
- Built from `<commit>` (tag `vX.Y.Z`), `vcs.modified=false`:

```sh
go version -m crower | grep vcs.revision
```

### ⚠️ Known limitations

- <open issue that affects this release, or omit the section entirely>.
- Windows CLI binaries are cross-compiled unless a Windows runner executed them.

### 📝 Changelog

[`vPREV...vX.Y.Z`](https://github.com/fiwon123/crower/compare/vPREV...vX.Y.Z) — one sentence on the main theme.
```
