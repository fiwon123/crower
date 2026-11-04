<div align="center">

# 📖 Command reference

**Every command, every flag, and the things a script can rely on.**

</div>

`crower <command> --help` prints the same information for one command. This page
is the version you can read without a terminal.

---

## 🧭 The contract

Four rules hold for every command. A script may depend on all of them.

| | |
| --- | --- |
| **Results** go to **stdout** | the table of `list`, the rows of `search` |
| **Diagnostics** go to **stderr** | prompts, status lines, warnings, errors |
| **Exit 0** | the work happened — or the user declined, and nothing was the answer |
| **Exit 1** | the work did not happen |

```bash
# only the table reaches the file
crower list > commands.txt

# a failure is visible to the script
if ! crower execute deploy; then echo "deploy failed" >&2; fi
```

A prompt with **no input available** — a closed stdin, or `--no-input` — is an
error, not an empty answer. The run fails instead of waiting.

## ⚙️ Global flags

These work on every command.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--config <path>` | | configuration file (default `$HOME/crower/crower.yaml`) |
| `--no-input` | | never prompt; fail instead of waiting for a person |
| `--version` | `-v` | print the version |
| `--upgrade` | | upgrade to the latest release |
| `--check` | | check whether a newer release exists |

`--config` and `--no-input` are persistent: every subcommand inherits them.

## 🤖 Machine-readable results

`list` and `search` accept `--json`: the result becomes **one JSON object** on
stdout, diagnostics stay on stderr, and the exit status is unchanged.

```bash
crower list --json | jq -r '.commands[].name'
```

| Command | Object |
| --- | --- |
| `crower list` | `{"commands":[{"row":0,"name":"gs","aliases":["g"],"exec":"git status"}]}` |
| `crower list --history` | `{"history":[{"version":1,"file":"00001.yaml","timestamp":"2026-01-02T15:04:05Z","operation":"create","command":"gs","note":"…"}]}` |
| `crower list --folder <dir>` | `{"entries":[{"name":"a.txt","dir":false,"size":12,"mode":"-rw-r--r--","modified":"…"}]}` |
| `crower list --process` | `{"processes":[…]}` |
| `crower list --system` | `{"variables":[{"name":"PATH","values":["…"]}]}` |
| `crower list --syspath` | `{"variable":{"name":"PATH","values":["…"]}}` |
| `crower search` | `{"search":{"kind":"file","query":"main","path":".","results":["…"]}}` |

Any other command rejects `--json` as an unknown flag: a global switch a command
silently ignored would report success without honouring it.

---

## ➕ crower create

```text
crower create [name] [exec]
```

Create a command, a file or a folder. With no arguments, crower asks for each
field.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--alias <a,b>` | `-a` | add aliases (comma-separated, repeatable) |
| `--process <name>` | `-p` | create the command from a running process; the argument is its name or pid |
| `--file` | `-f` | treat every argument as a file path to create |
| `--folder` | `-o` | treat every argument as a folder path to create |
| `--script` | | create a script file that the command runs |
| `--system` | | create a system variable |
| `--syspath` | | add a directory to the system `PATH` |
| `--execute` | | create a command from the last command you ran |

```bash
crower create greet "echo hello"
crower create greet "echo hello" --alias hi
crower create --file notes.txt
crower create --folder reports
crower create -p server nginx
```

A command runs through `sh -c` on Linux and macOS, and `cmd /c` on Windows. A
Steps command runs its steps instead of its `exec`; `execute` refuses to run one
with parameters it cannot apply.

## ✏️ crower update

```text
crower update [command]
```

Change a command's name, aliases or exec. The command is found by name or alias.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--name <name>` | `-n` | new name |
| `--alias <a,b>` | `-a` | new alias list, replacing the old one |
| `--exec <line>` | `-e` | new command line |
| `--last` | `-l` | update the last executed command |
| `--create` | `-c` | update the last created command |
| `--execute` | `-x` | update the last updated command |

```bash
crower update greet --exec "echo hi"
crower update greet --name hello --alias hi
crower update --last
```

With no field given, crower asks for each one. A Steps command keeps its steps:
`update` has no flag for them, so it never replaces them.

## 🗑️ crower delete

```text
crower delete [command]
```

Delete a command, a file or a folder.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--create` | `-c` | delete the last created command |
| `--update` | `-u` | delete the last updated command |
| `--execute` | `-x` | delete the last executed command |
| `--history` | | delete one history entry and its snapshot |
| `--file` | `-f` | treat every argument as a file path |
| `--folder` | `-o` | treat every argument as a folder path |
| `--system` | | delete a system variable |
| `--syspath` | | remove a directory from the system `PATH` |

```bash
crower delete greet
crower delete hi
crower delete --file notes.txt
crower delete --create
```

With no arguments crower asks which command to delete. Deleting is recorded in
the history, so `revert` can bring the command back.

## 📋 crower list

```text
crower list [flags]
```

| Flag | Short | Meaning |
| --- | --- | --- |
| `--process` | `-p` | list running processes |
| `--history` | `-i` | list the history, newest first |
| `--folder <dir>` | `-o` | list the files in a folder |
| `--system` | | list system variables |
| `--syspath` | | list the system `PATH` entries |
| `--json` | | one JSON object on stdout instead of a table |

```bash
crower list
crower list --history
crower list --folder ~/projects
crower list --json
```

With no flag, crower lists the managed commands in the order they were saved.

## ▶️ crower execute

```text
crower execute <command> [params...]
```

Run a managed command by name or alias. `crower <name>` is the same thing.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--last` | `-l` | run the last executed command |
| `--create` | `-c` | run the last created command |
| `--update` | `-u` | run the last updated command |

```bash
crower execute greet
crower greet
crower execute --last
```

Parameters after the command fill `{1}`, `{2}`, … in its steps; `{}` takes the
rest. Each parameter is quoted before it reaches the shell.

## 📂 crower open

```text
crower open [path...]
```

| Flag | Short | Meaning |
| --- | --- | --- |
| `--file` | `-f` | open a file with the system opener |
| `--folder` | `-o` | open a folder with the system opener |
| `--system` | `-s` | open the system settings |

```bash
crower open                 # the configuration file
crower open --file notes.txt
crower open --folder .
```

## 📄 crower copy

```text
crower copy <source...> <destination>
```

Copy files or folders. The last argument is the destination.

```bash
crower copy notes.txt ~/backup
crower copy a.txt b.txt ~/backup
crower copy ./reports ~/backup
```

## 🚚 crower move

```text
crower move <source...> <destination>
```

Move files or folders. The last argument is the destination.

```bash
crower move notes.txt ~/archive
crower move ./reports ~/archive
```

## 🔍 crower search

```text
crower search <term> [path]
```

| Flag | Short | Meaning |
| --- | --- | --- |
| `--file` | `-f` | search for a file name |
| `--folder` | `-o` | search for a folder name |
| `--browser` | `-b` | search the web in the default browser |
| `--json` | | one JSON object on stdout instead of a table |

```bash
crower search report ./notes
crower search --folder projects ~
crower search --browser "golang context"
```

## 🗜️ crower extract

```text
crower extract <archive...>
```

Unpack `.zip`, `.tar`, `.tar.gz`, `.7z`, `.rar` and friends.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--out <dir>` | `-o` | extract into this folder instead of the archive's own |

```bash
crower extract files.zip
crower extract files.zip --out ./build
crower extract "*.zip"
```

## ↩️ crower revert

```text
crower revert <steps>
```

Undo operations, newest first. The argument is how many history entries to roll
back.

```bash
crower revert 1
crower revert 3
```

`crower list --history` shows what is there. Reverting writes the restored
command map back and drops the entries it passed.

## ♻️ crower restore

```text
crower restore <command>
```

Bring back a command from the snapshot of a past operation, by name or alias.
This is how a deleted or overwritten command comes back without undoing
everything after it.

```bash
crower restore greet
```

## 🔄 crower reset

```text
crower reset
```

Delete every managed command. The table is shown first and nothing happens
unless you confirm. History snapshots are left alone: `crower list --history`
still works and `crower revert` can bring the commands back.

```bash
crower reset
```

## 📤 crower export

```text
crower export [name]
```

Export one command, or all of them, to a TOML file for sharing.

| Flag | Short | Meaning |
| --- | --- | --- |
| `--out <path>` | `-o` | write here instead of the default file |

```bash
crower export
crower export greet
crower export --out ./commands.toml
```

## 📥 crower import

```text
crower import <file>
```

Import commands from a TOML file.

| Flag | Meaning |
| --- | --- |
| `--overwrite` | replace a command that already exists with the same name |

```bash
crower import ./commands.toml
crower import ./commands.toml --overwrite
```

Names and aliases share one namespace: an import that would collide is renamed
rather than silently shadowing a command.

## 🤝 crower share

List and pull command packs from a registry.

### crower share list

List the packs the registry publishes.

```bash
crower share list
```

### crower share pull

Pull a pack into the command set.

```bash
crower share pull tools
```

## 🎛️ crower gui

Open the interactive terminal UI.

```bash
crower gui
```

## 🪟 crower desktop

Open the native window UI (cgo-free).

```bash
crower desktop
```

---

## 🩺 crower doctor

```bash
crower doctor
crower doctor --json
```

Shows the configuration file and folder, whether the folder is writable, how
many commands and history entries exist, the update-check mode, the detected
shell and editor, whether the diagnostics stream is a terminal, and the
platform. `--json` prints the same as one object.

It reads only: nothing is created beyond the config folder the run needs.

## 🚀 crower init

```bash
crower init            # show the line for your shell
crower init --apply    # append it to the startup file
```

Prints the PATH line for bash, zsh or fish and where it goes. `--apply` appends
it once; a second run is a no-op. It also prints the completion command for the
detected shell.

## 🧾 Exit codes

| Status | Meaning |
| --- | --- |
| `0` | the work happened, help was printed, or the user declined a confirmation |
| `1` | the command failed; the reason is one line on stderr |

Nothing else is used: a script can test the status without decoding it.

## 🔤 A note on shorthands

`-o` means `--folder` on `create`, `delete`, `list`, `search` and `open`, and
`--out` on `extract` and `export`. It is the one shorthand that changes meaning
between commands; write the long flag in scripts.
