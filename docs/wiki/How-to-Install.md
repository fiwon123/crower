<div align="center">

# 📦 Installation

**One archive, three executables: `crower`, `cr` and `crower-gui`.**

</div>

## ⬇️ Download a release (recommended)

Grab the latest archive from **[Releases](https://github.com/fiwon123/crower/releases)**.
Archives are named `crower_<version>_<os>_<arch>`. Every one contains the same
three executables — `crower`, `cr` and `crower-gui` — with no architecture
suffix.

| Platform | Archive |
| --- | --- |
| 🐧 Linux amd64 | `crower_<v>_linux_amd64.tar.gz` |
| 🐧 Linux arm64 | `crower_<v>_linux_arm64.tar.gz` |
| 🪟 Windows amd64 | `crower_<v>_windows_amd64.zip` |
| 🪟 Windows arm64 | `crower_<v>_windows_arm64.zip` |
| 🍎 macOS Intel | `crower_<v>_darwin_amd64.tar.gz` |
| 🍎 macOS Apple Silicon | `crower_<v>_darwin_arm64.tar.gz` |

`crower` and `cr` are cgo-free and need no display. `crower-gui` is the Fyne
window; it needs a graphical session.

```bash
VERSION=v0.5.2   # replace with the tag you are installing
curl -LO "https://github.com/fiwon123/crower/releases/download/$VERSION/crower_${VERSION}_linux_amd64.tar.gz"
tar -xzf "crower_${VERSION}_linux_amd64.tar.gz"
chmod +x crower cr crower-gui
sudo mv crower cr crower-gui /usr/local/bin/
```

> [!TIP]
> **No `sudo`?** The archive ships `add_to_path.sh` — extract it next to the
> binaries and run `./add_to_path.sh`. It adds that folder to the config file of
> the shell you are running (`~/.bashrc`, `~/.zshrc`, …).

**Windows:** extract the `..._windows_amd64.zip` (or `..._arm64.zip`), then run
`add_to_path.bat` from the extracted folder — or add the folder containing
`crower.exe`, `cr.exe` and `crower-gui.exe` to your `PATH`.

`crower.exe` and `cr.exe` are terminal programs. **Double-click
`crower-gui.exe`** — the one carrying the crow icon — to open the window; it is
built as a graphical program, so no console window opens with it.

Verify either install:

```bash
crower --version
```

## 🧪 Build from source

```bash
git clone https://github.com/fiwon123/crower.git
cd crower
go build -o crower .
./crower --help
```

Go **1.25.4** is pinned in `go.mod`; the dev sandbox and CI use the same version.

> [!WARNING]
> A build from source proves something about *your tree*, not about the archive
> users download. If you are reporting a bug, reproduce it on a release build
> first, and say which version you used.

## 🐳 Try it without installing

```bash
git clone https://github.com/fiwon123/crower.git
cd crower
docker compose run --rm dev crower --help
```

Or run the agent sandbox with `make dev-up` — see **[Development](Development)**.

## ✅ Next

**[⚡ Quick start →](Quick-Start)** — your first shortcut in four commands.
