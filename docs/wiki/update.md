# ✏️ `update`

**Change a stored command — its name, aliases, or the command itself.**

```bash
crower update [command] [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-n`, `--name` | New name |
| `-a`, `--alias` | New aliases, comma-separated |
| `-e`, `--exec` | New command text |
| `--icon` | New icon; an empty value removes it |
| `--group` | New group; an empty value removes it |
| `-l`, `--last` | Update the most recent entry |
| `-c`, `--create` | Update the most recent `create` |
| `-x`, `--execute` | Update the most recent `execute` |

## 💡 Examples

```bash
crower update                        # interactive
crower update deploy                 # pick fields, interactively
crower update deploy -e "git pull && make"   # change the command
crower update deploy -n release -a rel       # rename and re-alias
crower update --last -e "make test"          # edit what you just ran
```

> [!NOTE]
> Every update appends to the history, so the previous text is still reachable
> with [`revert`](revert).

> [!NOTE]
> `Description`, `Params`, `Shell`, `Cwd`, `Env` and `Timeout` are part of the
> command schema, but `update` has no flags for them yet: edit them in the
> window's edit dialog, in the configuration file, or by re-importing a pack.
> See [`execute`](execute#-execution-target). `--icon` and `--group` are the
> exception: they are display metadata, so they have flags too.

## 🔗 Related

[`create`](create) · [`delete`](delete) · [`list`](list) · **[📖 Commands](How-to-Use)**
