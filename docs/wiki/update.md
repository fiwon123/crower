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
| `--description` | What the command is for; an empty value removes it |
| `--params` | Positional parameters, comma-separated or repeated; an empty value removes them |
| `--cwd` | Working directory the command runs in; an empty value removes it |
| `--shell` | Shell that runs the command (`sh`, `pwsh`, `cmd`); an empty value removes it |
| `--timeout` | Longest the command may run, as a Go duration (`120s`); an empty value removes it |
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
> A flag that is not given leaves its field alone; giving it an empty value
> clears the field. `Env` has no flag — it holds values that look like secrets,
> so it stays a configuration-file field. See
> [`execute`](execute#-execution-target) for what each field does at run time.

## 🔗 Related

[`create`](create) · [`delete`](delete) · [`list`](list) · **[📖 Commands](How-to-Use)**
