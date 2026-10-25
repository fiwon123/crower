<p align="center">
  <img src="assets/logo.svg" alt="crower logo" width="128" height="128">
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

- [Why crower?](#why-crower)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Usage Recipes](#usage-recipes)
- [Commands](#commands)
- [Sharing Crower Files](#sharing-crower-files)
- [Configuration](#configuration)
- [Development](#development)
- [License](#license)

## Why crower?

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

## Quick Start

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

## Installation

### Download Release (Recommended)

Download the latest release for your platform from the [releases page](https://github.com/fiwon123/crower/releases).

**Linux:**
```bash
tar -xzf crower_*.tar.gz
chmod +x crower cr
sudo mv crower cr /usr/local/bin/
```

**Windows:**
Extract the ZIP and add the directory containing `crower.exe` and `cr.exe` to your PATH.

### Build from Source

```bash
git clone https://github.com/fiwon123/crower.git
cd crower
make build
sudo make install  # or add build/ to PATH
```

**Requirements:** Go 1.25.4 or later

## Usage Recipes

### Development Workflow
```bash
# Create shortcuts for common dev tasks
crower create test "go test ./..."
crower create build "go build ./..."
crower create lint "golangci-lint run"
crower create deploy "git push && ssh prod 'cd app && git pull && systemctl restart app'"

# Run your workflow
cr test && cr build
```

### Multi-step Automation
```bash
# Create a command with multiple steps (separated by newlines in exec)
crower create release "echo 'Building...'\ngo build -o app .\necho 'Testing...'\ngo test ./...\necho 'Done!'"
cr release
```

### File Management Made Easy
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

### History & Undo
```bash
# View history
cr list -i

# Revert last N operations
cr revert 2

# Restore a deleted command
cr restore <key>
```

## Commands

| Command | Description | Example |
|---------|-------------|---------|
| `create` | Create a command, file, folder, or script | `crower create name "command"` |
| `update` | Update existing commands or last operation | `crower update -n oldname -n newname` |
| `delete` | Delete commands, files, folders, or history | `crower delete name` |
| `list` | List commands, history, processes, or files | `crower list` or `cr list -i` |
| `execute`/`run` | Execute stored commands | `crower execute name` or just `cr name` |
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

## Sharing Crower Files

**Build once, share everywhere.** Perfect for onboarding, CI setups, or sharing workflows with your team.

### Export Your Workflows
```bash
# Export all commands
crower export -o my-workflows.toml

# Export a specific command
crower export deploy -o deploy.toml
```

### Import Workflows
```bash
# Import from a file
crower import my-workflows.toml

# Import and overwrite conflicts
crower import my-workflows.toml --overwrite
```

### Pull from a Registry
```bash
# Browse available packs
crower share list

# Pull and install a pack
crower share pull docker-setup
```

**Registry Protocol:** Any HTTP server can host crower packs using the simple manifest protocol (`index.json` + `packs/<name>.toml`). Check the [crower-resources](https://github.com/fiwon123/crower-resources) repo for examples.

## Configuration

Crower stores data in `~/crower/` by default:
- `crower.yaml` - Your stored commands (TOML format)
- `historyhelper.json` - Operation history
- `history/` - Command snapshots for revert/restore
- `crower.log` - Execution logs

Override with `--config /path/to/config.yaml`

## Development

```bash
# Run tests
make check

# Build for all platforms
make release

# Run from source
go run ./ --help
```

See [AGENTS.md](AGENTS.md) and [PROJECT_CONTEXT.md](PROJECT_CONTEXT.md) for architecture details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
