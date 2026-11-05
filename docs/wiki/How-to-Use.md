<div align="center">

# 🗂️ Command map

**Every command, one line each.**

**Flags, examples, exit codes and the scripting contract** are in the
**[command reference](Command-Reference)**.

</div>

## 🗂️ Command map

| Command | What it does | Example |
| --- | --- | --- |
| 🎛️ [`gui`](gui) | Terminal UI | `crower gui` |
| 🖼️ [`crower-gui`](gui) | Fyne window (same archive) | `crower-gui` |
| ➕ [`create`](create) | Command, file, folder, or script | `crower create name "cmd"` |
| ✏️ [`update`](update) | Change a command, or the last operation | `crower update name -n new` |
| 🗑️ [`delete`](delete) | Command, file, folder, or history entry | `crower delete name` |
| 📋 [`list`](list) | Commands, history, processes, files | `crower list -i` |
| ▶️ [`execute`](execute) | Run a stored command | `crower execute name` |
| 📂 [`open`](open) | File, folder, or system settings | `crower open -f file.txt` |
| 📄 [`copy`](copy) | Copy files or folders | `crower copy src dst` |
| 🚚 [`move`](move) | Move files or folders | `crower move src dst` |
| 🔍 [`search`](search) | Files, folders, or the web | `crower search -f pattern` |
| 🗜️ [`extract`](extract) | Unpack archives | `crower extract a.tar.gz` |
| ↩️ [`revert`](revert) | Undo N history steps | `crower revert 1` |
| ♻️ [`restore`](restore) | Bring back a deleted command | `crower restore name` |
| 🧹 [`reset`](reset) | Remove every command | `crower reset` |
| 🤝 [`share` / `export` / `import`](share) | Move commands between machines | `crower export -o team.toml` |

## 🌐 Global flags

Every command accepts these:

| Flag | Effect |
| --- | --- |
| `--config <path>` | Use another configuration file (default `$HOME/crower/crower.yaml`) |
| `--no-input` | Never prompt — fail with a message instead of waiting |
| `-h`, `--help` | Help for the command |
| `-v`, `--version` | Print the version |
| `--check` | Ask whether a newer release exists |
| `--upgrade` | Self-update to the latest release |

## 🔑 Prefix shortcuts

| You type | You get |
| --- | --- |
| `crower <name>` | Runs the stored command `name` |
| `cr <name>` | Same thing, shorter |
| `crower execute <name>` | Same thing, explicit |

> [!NOTE]
> `crower <name>` falls back to executing a stored command, so a name that
> collides with a subcommand (`list`, `create`, …) cannot be run that way — use
> `crower execute <name>` for those.

## 📌 Common patterns

```bash
crower create --alias b,bkp "backup"    # multiple aliases, comma-separated
crower list -i                          # history
crower update --last -e "new command"   # edit the most recent entry
crower delete --history                 # drop a history entry and its backup
crower create --file a.txt b.txt        # several files at once
```

## 🧩 From a script

Exit codes, stdout/stderr separation and `--no-input` are documented in
**[Scripting](Scripting)**.
