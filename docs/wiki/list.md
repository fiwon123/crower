# 📋 `list`

**List commands (the default), history, processes, files, or variables.**

```bash
crower list [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-i`, `--history` | Operation history, newest first |
| `-p`, `--process` | Running processes |
| `-o`, `--folder` | Files in a folder (pass the path) |
| `--system` | System variables |
| `--syspath` | `PATH` entries |
| `--json` | One JSON object instead of the table ([shape](Scripting#-machine-readable-results---json)) |

## 💡 Examples

```bash
crower list                          # every stored command
crower list -i                       # what ran, and when
crower list -p                       # processes
crower list --folder ~/projects
crower list --system
crower list --json                     # the same list, as one JSON object
crower list -i --json | jq '.history[0]'
```

> [!NOTE]
> The table goes to **stdout**; status lines go to stderr. `crower list > file`
> captures the list and nothing else.

The history table has a `Command` column: the command an entry touched, or a
dash for an operation that touched a file or a variable.

## 🔗 Related

[`execute`](execute) · [`revert`](revert) · [🧩 Scripting](Scripting) · **[📖 Commands](How-to-Use)**
