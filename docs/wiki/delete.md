# 🗑️ `delete`

**Delete a command, a file, a folder, or a history entry.**

```bash
crower delete [command] [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-f`, `--file` | Treat the arguments as file paths |
| `-o`, `--folder` | Treat the arguments as folder paths |
| `-c`, `--create` | Delete the most recent `create` |
| `-u`, `--update` | Delete the most recent `update` |
| `-x`, `--execute` | Delete the most recent `execute` |
| `--history` | Delete a history entry **and** its linked backup |
| `--syspath` | Delete a `PATH` entry |
| `--system` | Delete a system variable |

## 💡 Examples

```bash
crower delete                    # interactive
crower delete deploy             # by name or alias
crower delete --file old.log
crower delete --folder tmp build
crower delete --history          # one history entry, with its snapshot
```

> [!TIP]
> A delete is reversible: [`revert`](revert) restores the previous state, and
> [`restore`](restore) brings back a command you removed.

## 🔗 Related

[`revert`](revert) · [`restore`](restore) · [`reset`](reset) · **[📖 Commands](How-to-Use)**
