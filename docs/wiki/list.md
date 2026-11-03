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

## 💡 Examples

```bash
crower list                          # every stored command
crower list -i                       # what ran, and when
crower list -p                       # processes
crower list --folder ~/projects
crower list --system
```

> [!NOTE]
> The table goes to **stdout**; status lines go to stderr. `crower list > file`
> captures the list and nothing else.

## 🔗 Related

[`execute`](execute) · [`revert`](revert) · [🧩 Scripting](Scripting) · **[📖 Commands](How-to-Use)**
