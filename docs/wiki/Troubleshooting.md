<div align="center">

# 🚑 Troubleshooting

**Start here; most of these have a one-line fix.**

</div>

## `crower: command not found`

The archive was extracted but its folder is not on `PATH`. Either move both
binaries to a directory that is (`/usr/local/bin`), or run the bundled
`add_to_path.sh` / `add_to_path.bat` from the folder you extracted them into, and
open a new terminal.

## A shortcut runs the wrong thing

A stored name that collides with a subcommand is not reachable as
`crower <name>`. Use the explicit form:

```bash
crower execute list      # your command named "list", not the subcommand
```

## `error: …` on stderr and exit code 1

That is the designed failure path: the work did not happen and the reason is the
one line printed. Common causes: a path that does not exist, a name that
conflicts, or a prompt declined under `--no-input`.

## It hangs waiting for input

Something needs an answer. If you are in a script, pass `--no-input` so it fails
fast instead of waiting; interactively, just answer the prompt.

## `revert` did nothing

`revert` walks history backwards, one step per call:

```bash
crower list -i        # what is there, newest first
crower revert 1       # undo the most recent entry
crower revert 3       # or the last three
```

If history is empty, the operation never appended — check that you used the
`crower`/`cr` binary and not a shell alias.

## Restored the wrong thing

```bash
crower list -i                       # find the entry
crower restore <name>                # bring back a deleted command
crower revert 1                      # or undo the last operation
```

## A command behaves differently on Windows

A command line is one line. `cmd` is line-oriented: a CR/LF inside a value
starts a new command, and no quoting prevents it. `%VAR%` expands on the
command line and cannot be quoted away either. Keep `Exec` single-line and use
`Env` with `${VAR}` for values that must not be expanded.

## The GUI window does not open

`crower-gui` is the Fyne app included in every platform archive and needs a
graphical session. In a terminal-only environment use `crower gui`, which
renders in the terminal itself.

On Windows, run `crower-gui.exe` — the one with the crow icon. It is built as a
graphical program, so a console window does not open and closing one cannot take
the window with it. If a build without the icon or the graphical subsystem is in
front of you, reinstall from the release archive: the release checks both.

## Still stuck

Open an [issue](https://github.com/fiwon123/crower/issues/new) with:

- `crower --version`
- your OS and architecture
- the exact command and the exact output
- whether it reproduces from a **release build**
