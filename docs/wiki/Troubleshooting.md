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

## The window does not open on Windows

`crower-gui` is a graphical program: it has no console, so a start that fails
inside Fyne used to look like a double-click that does nothing. It now writes
what happened to `fyne.log` beside `crower.yaml`
(`%USERPROFILE%\crower\fyne.log`), shows the reason in a dialog, and leaves the
file behind. A run that opens normally removes it.

The most common reason is that the display driver has no **OpenGL 2.1**. Windows
reports this as `window creation error` / `WGL: The driver does not appear to
support OpenGL`. Remote Desktop sessions and machines using the basic display
driver do this:

- connect to the machine locally, or over a console session rather than RDP;
- install the graphics driver for the GPU, or update it;
- meanwhile, `crower gui` opens the same tools in a terminal.

## Windows says "This app can't run on your PC"

That dialog is the Windows loader rejecting the file, before crower runs a single
line. For a release download the cause is the architecture: the release ships one
archive per architecture, and an arm64 archive does not run on an Intel/AMD PC
(or the reverse).

Check which machine the PC is — Settings → System → About → System type:

- "64-bit operating system, x64-based processor" → `crower_<version>_windows_amd64.zip`
- "64-bit operating system, ARM-based processor" → `crower_<version>_windows_arm64.zip`

Every Windows archive contains `WINDOWS.txt`, which names the architecture it
was built for. A file that was damaged on the way down fails the same way: check
the download against the release's `checksums.txt` before extracting.

```powershell
certutil -hashfile crower_<version>_windows_amd64.zip SHA256
```

## Still stuck

Open an [issue](https://github.com/fiwon123/crower/issues/new) with:

- `crower --version`
- your OS and architecture
- the exact command and the exact output
- whether it reproduces from a **release build**
