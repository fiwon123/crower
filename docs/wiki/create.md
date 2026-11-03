# ➕ `create`

**Create a command, a file, a folder, or a script.**

```bash
crower create [name] [exec] [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-a`, `--alias` | Aliases, comma-separated (`--alias 'a1,a2'`) |
| `-f`, `--file` | Treat the arguments as file paths |
| `-o`, `--folder` | Treat the arguments as folder paths |
| `--execute` | Create from the last executed command |
| `-p`, `--process` | Create from a process name or pid |
| `--script` | Create a helper script for a complex command |
| `--syspath` | Create a `PATH` entry |
| `--system` | Create a system variable |

## 💡 Examples

```bash
# Interactive: asks for name, command and aliases
crower create

# Name + command
crower create deploy "ssh user@server 'cd /app && git pull'"

# With aliases
crower create backup -a bkp,bk "tar -czf ~/b.tgz ~/data"

# Files and folders
crower create --file notes.md todo.txt
crower create --folder tmp logs

# Re-create the command you last executed under a new name
crower create --execute greet2
```

> [!TIP]
> Quote the command as one argument. `crower create greet "echo hi"` stores
> `echo hi`; without quotes it sees three words.

## 🔗 Related

[`update`](update) · [`execute`](execute) · [`delete`](delete) · **[📖 Commands](How-to-Use)**
