<p align="center">
  <img src="assets/logo.svg?v=2" alt="crower logo — a crow perched on a terminal cursor" width="128" height="128">
</p>

<h1 align="center">crower</h1>

<p align="center">
  <strong>A cross-platform terminal file manager and command-alias manager</strong><br>
  <em>Shortcut any command, automate workflows, and share them with your team</em>
</p>

<p align="center">
  <a href="https://github.com/fiwon123/crower/actions/workflows/ci.go.yml"><img src="https://github.com/fiwon123/crower/actions/workflows/ci.go.yml/badge.svg" alt="CI"></a>
  <a href="https://github.com/fiwon123/crower/releases"><img src="https://img.shields.io/github/v/release/fiwon123/crower" alt="Release"></a>
  <a href="https://github.com/fiwon123/crower/blob/main/LICENSE"><img src="https://img.shields.io/github/license/fiwon123/crower" alt="License"></a>
</p>

---

## Table of Contents

- [🔖 Why crower?](#-why-crower)
- [🚀 Quick start](#-quick-start)
- [📦 Installation](#-installation)
- [🍳 Usage recipes](#-usage-recipes)
- [📖 Commands](#-commands)
- [🧩 Scripting](#-scripting)
- [🤝 Sharing crower files](#-sharing-crower-files)
- [⚙️ Configuration](#-configuration)
- [🛠️ Development](#-development)
- [📚 Documentation](#-documentation)
- [💬 Contributing](#-contributing)
- [⚠️ What is not proven](#-what-is-not-proven)
- [📄 License](#-license)

## 🔖 Why crower?

Tired of retyping long, complex commands? Need to remember obscure CLI incantations? **crower** lets you:

- **🔖 Shortcut anything** - Turn any shell command into a memorable alias
- **⚡ Run multi-step workflows** - Chain multiple commands together with automatic error handling
- **🔗 Nest commands** - Call crower commands from within other crower commands with recursion protection
- **📝 Add context** - Attach notes to your commands and track history
- **↩️ Undo anything** - Full history with revert capability
- **📁 Manage files** - Create, copy, move, open files and folders with simple commands
- **🔍 Search** - Find files/folders or search the web instantly
- **📦 Share workflows** - Export and import crower files, or download from a registry
- **🧠 Smart execution** - Proper parameter quoting, placeholders (`{1}`, `{2}`, `{}`), multi-step chains, and safe execution

## 🚀 Quick start

```bash
# Create a shortcut for a complex command
crower create deploy "ssh user@server 'cd /app && git pull && systemctl restart app'"

# Run it by name or alias
crower deploy
# or
cr deploy

# Create with aliases
crower create backup -a bkp "tar -czf backup.tar.gz ~/important-data"
cr bkp

# Run with parameters
crower create greet "echo Hello {1}, welcome to {2}!"
cr greet John Developer
```

## 📦 Installation

### Download a release (recommended)

Download the latest release for your platform from the [releases page](https://github.com/fiwon123/crower/releases).

Archives are named `crower_<tag>_<os>_<arch>` — for example
`crower_v0.4.0_linux_amd64.tar.gz`. The binaries inside are always `crower` and
`cr`, with no architecture suffix. Replace `v0.4.0` below with the tag you are
installing.

**Linux** (swap `linux_amd64` for `linux_arm64`, `darwin_amd64` or
`darwin_arm64` on another platform):

```bash
VERSION=v0.4.0
curl -LO "https://github.com/fiwon123/crower/releases/download/$VERSION/crower_${VERSION}_linux_amd64.tar.gz"
tar -xzf "crower_${VERSION}_linux_amd64.tar.gz"
chmod +x crower cr
sudo mv crower cr /usr/local/bin/
```

Without `sudo`: the archive also ships `add_to_path.sh`, which adds the
directory it lives in to the config file of the shell you are running. Extract
it next to the binaries and run `./add_to_path.sh`.

**Windows:** extract `crower_<tag>_windows_amd64.zip` (or `..._arm64.zip`), then
either run `add_to_path.bat` from the extracted folder or add the folder
containing `crower.exe` and `cr.exe` to your PATH.

Verify either install with `crower --version`.

### Build from source

```bash
git clone https://github.com/fiwon123/crower.git
cd crower
go build -o crower .
./crower --help

# short alias, then put both binaries on your PATH
cp crower cr
sudo mv crower cr /usr/local/bin/
```

`make help` lists every target: `make release` builds the Windows and Linux
release archives into `build/`, and `make check` runs the local gate — gofmt,
vet, test and build.

**Requirements:** Go 1.25.4 or later

## 🍳 Usage recipes

### 🛠️ Development workflow
```bash
# Create shortcuts for common dev tasks
crower create test "go test ./..."
crower create build "go build ./..."
crower create lint "golangci-lint run"
crower create deploy "git push && ssh prod 'cd app && git pull && systemctl restart app'"

# Run your workflow
cr test && cr build
```

### Multi-step automation
```bash
# Create a command with multiple steps (separated by newlines in exec)
crower create release "echo 'Building...'\ngo build -o app .\necho 'Testing...'\ngo test ./...\necho 'Done!'"
cr release
```

### File management made easy
```bash
# Create files and folders
crower create project -f  # creates file
crower create src -o     # creates folder

# Copy and move
crower copy file.txt ./backup/
crower move old.txt ./archive/

# Open files/folders
crower open -f README.md
crower open -o ./src
```

### History and undo
```bash
# View history
cr list -i

# Revert last N operations
cr revert 2

# Restore a deleted command
cr restore <key>
```

## 📖 Commands

| Command | Description | Example |
|---------|-------------|---------|
| `gui` | Launch interactive TUI GUI | `crower gui` |
| `desktop` | Launch native desktop window UI (cgo-free) | `crower desktop` |
| `create` | Create a command, file, folder, or script | `crower create name "command"` |
| `update` | Update existing commands or last operation | `crower update -n oldname -n newname` |
| `delete` | Delete commands, files, folders, or history | `crower delete name` |
| `list` | List commands, history, processes, or files | `crower list` or `cr list -i` |
| `execute` | Execute stored commands | `crower execute name` or just `cr name` |
| `open` | Open files, folders, or system settings | `crower open -f file.txt` |
| `copy`/`move` | Copy or move files and folders | `crower copy src dst` |
| `search` | Search files or web | `crower search -f pattern` |
| `extract` | Extract archives | `crower extract archive.tar.gz` |
| `revert` | Revert history | `crower revert 1` |
| `restore` | Restore from history | `crower restore <key>` |
| `reset` | Reset all commands | `crower reset` |
| `export` | Export commands to share | `crower export` or `crower export name` |
| `import` | Import commands from file | `crower import pack.toml` |
| `share` | Browse and pull from registry | `crower share list` or `cr share pull pack` |

Use `crower <command> --help` for detailed options.

## 🧩 Scripting

crower is meant to be called from a script as well as typed by hand.

- **Results go to stdout, diagnostics to stderr.** `crower list` writes the
  table and nothing else to stdout; prompts, status lines, warnings and errors
  go to stderr, so piping the results never picks them up.
- **The exit status says whether the work happened.** `0` means it did — or that
  you answered *no* to a confirmation, because declining is a choice. `1` means
  it did not, and the reason is one line on stderr: `error: ...`.
- **A prompt never waits forever.** With `--no-input`, or with stdin closed, a
  command that needs an answer fails with a message instead of asking again.

```bash
# Take the result and leave the status lines behind
greeting="$(crower execute greet)"

# Fail the script when the delete did not happen, without a prompt to block it
crower --no-input delete old-command || exit 1
```

`--config <path>` points a run at another configuration file, which is how a
script uses its own set of commands without touching your personal one.

## 🤝 Sharing crower files

**Build once, share everywhere.** Perfect for onboarding, CI setups, or sharing workflows with your team.

### Export your workflows
```bash
# Export all commands
crower export -o my-workflows.toml

# Export a specific command
crower export deploy -o deploy.toml
```

### Import workflows
```bash
# Import from a file
crower import my-workflows.toml

# Import and overwrite conflicts
crower import my-workflows.toml --overwrite
```

### Pull from a registry
```bash
# Browse available packs
crower share list

# Pull and install a pack
crower share pull docker-setup
```

**Registry Protocol:** Any HTTP server can host crower packs using the simple manifest protocol (`index.json` + `packs/<name>.toml`). Check the [crower-resources](https://github.com/fiwon123/crower-resources) repo for examples.

## ⚙️ Configuration

Crower stores data in `~/crower/` by default:
- `crower.yaml` - Your stored commands (TOML format)
- `historyhelper.json` - Operation history
- `history/` - Command snapshots for revert/restore
- `crower.log` - Execution logs

Override with `--config /path/to/config.yaml`

## 🛠️ Development

```bash
# The local gate: gofmt, vet, go test, go build
make check

# Run tests only
go test ./...

# Build for all platforms
make release

# Run from source
go run ./ --help
```

See [AGENTS.md](AGENTS.md) and [PROJECT_CONTEXT.md](PROJECT_CONTEXT.md) for architecture details.

The dev sandbox runs the same pinned Go toolchain on every machine:
`make dev-up` (or `docker compose up` on any OS with Docker). Native Windows
uses `python scripts/dev-up.py dev-up`. Run `make help` for every target.

## 📚 Documentation

| I want to… | Read |
|---|---|
| see real recipes | [Usage recipes](#-usage-recipes) |
| look up a command or flag | [Commands](#-commands) |
| script against crower | [Scripting](#-scripting) |
| browse the wiki | [crower wiki](https://github.com/fiwon123/crower/wiki) |
| see what is planned | [Roadmap board](https://github.com/users/fiwon123/projects/5) |
| understand the internals | [AGENTS.md](AGENTS.md) · [PROJECT_CONTEXT.md](PROJECT_CONTEXT.md) |
| help prove a claim | [Verifying claims](https://github.com/fiwon123/crower/wiki/Verifying-claims) |

## 💬 Contributing

Issues, ideas, questions and pull requests are welcome — see
[CONTRIBUTING.md](CONTRIBUTING.md). One feature per pull request. Say what ran
and what did not, and do not describe a skipped check as a pass.

The most useful contribution right now is **running the unproven parts on real
hardware** — see [What is not proven](#-what-is-not-proven).

## ⚠️ What is not proven

This project says what it verified. CI formats, vets, tests and builds on
Linux/amd64, then cross-compiles six archives on the same runner — no gate
executes what it built.

- **The macOS build runs** — the darwin archives are never executed, and shell-backed operations are reported broken ([#129](https://github.com/fiwon123/crower/issues/129)).
- **The Windows build runs** — cross-compiled on Linux, never started.
- **The arm64 builds run** — built, never executed on ARM hardware.
- **The TUI and desktop surfaces render** — no gate draws a window; the tests assert the model, not pixels.
- **History and revert across a real session** — the suite uses isolated temp directories, not a day of real use.

Running the program and reporting exactly what happened — including "it worked"
— is the most useful contribution right now. The steps for each claim, and what
to record, are on the wiki's
[Verifying claims](https://github.com/fiwon123/crower/wiki/Verifying-claims) page.

## 📄 License

**crower is MIT licensed** — see [LICENSE](LICENSE).

The release archives also carry two files that are not crower's licence:

- [NOTICE](NOTICE) — attribution for the Apache-2.0 dependencies compiled into
  the binary (`spf13/cobra`, `minio/selfupdate`, `hajimehoshi/ebiten`);
- [LICENSE-APACHE](LICENSE-APACHE) — the Apache License 2.0 text those
  dependencies require to travel with a distribution.

A dependency's licence does not change this project's: your MIT grant is the only
one you need for crower's own code, and the Apache files exist so the bundled
libraries arrive with the notices they are owed.
