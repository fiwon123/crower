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
| `-p`, `--process` | New command name; pass the process name or pid as an argument |
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
# (one name; --alias and the other create flags still apply)
crower create --execute greet2 -a g2

# Create a shortcut for a running process:
# -p is the new command name, the process comes last
crower create -p editor code -a ed
```

> [!TIP]
> Quote the command as one argument. `crower create greet "echo hi"` stores
> `echo hi`; without quotes it sees three words.

> [!NOTE]
> `create --execute` cannot rebuild a command that uses `Steps`: that command
> never runs its `--exec`, and a list of steps is not one command line. It
> fails with an error instead of creating something different.

> [!NOTE]
> `Description`, `Params`, `Shell`, `Cwd`, `Env` and `Timeout` are part of the
> command schema, but `create` has no flags for them yet: set them in the
> configuration file or import a pack. See
> [`execute`](execute#-execution-target).

## 🔗 Related

[`update`](update) · [`execute`](execute) · [`delete`](delete) · **[📖 Commands](How-to-Use)**
