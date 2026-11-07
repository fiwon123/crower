# 🎛️ `gui` · `crower-gui`

**Two entries over the same configuration.**

| Command | Surface |
| --- | --- |
| `crower gui` | Interactive terminal UI |
| `crower-gui` | Fyne window |

```bash
crower gui
crower-gui
```

Both read the same `crower.yaml`, so a command created in the CLI appears in
the window, and anything done there lands in the history. The Fyne window can
create, edit and delete commands with the same catalog rules and history notes
as the CLI.

The window is split into tabs:

| Tab | What it holds |
| --- | --- |
| **Commands** | The command cards, a filter, a setup row and a group row; a card opens into its fields, parameters and **Run**, **Edit** and **Delete** |
| **History** | Every recorded operation, newest first, as expandable cards: an operation icon, command, time and version open into the labelled fields (with the full date), the raw line, and **Run**, **Edit** and **Restore this state** |
| **Tools** | Copy, move, search, extract and open; every path field has **Browse…**, and each card has its own result square that opens when an action runs |
| **About** | Logo, version, description, the repository link, and the update controls — padded to the tab width, with the actions stacked vertically |

A command can carry an **icon** — a short marker such as 🚀 — shown beside its
name, and a **group**; the filter lists every group in use, and an opened card
shows which one the command belongs to. Both are set from the create
and edit dialogs, or with `crower create --icon/--group` and
`crower update --icon/--group`. The Icon control is a large emoji button:
tapping it opens the full catalogue with a category selector and a search over
names, aliases, tags and the glyph; **Clear** removes it. The catalogue omits
multi-codepoint sequences (joined faces, skin-tone variants, flags) that the
platform font commonly renders as a blank cell. Picking from the window
always yields a catalogue emoji; a value set with `--icon` is shown as stored.
The **Group** field is a dropdown of the groups in use plus **New group**, which
asks for a name and selects it, so a group never has to be typed.

`New command` sits at the top right of the Commands tab, above the list it adds
to; the **Filter** row under it narrows the cards by name, alias, exec, group and
description. Rules separate the create and filter controls, the setup and group
rows, and the cards themselves.

On the Commands tab, each command is a card: a name with a chevron that opens
into the labelled fields (name, exec, aliases, group, description, params, cwd,
shell, timeout and numbered steps), a parameter entry, and **Run**, **Edit** and
**Delete**. Opening one closes the previous card; a collapsed card carries no
actions. **Run** passes the entry's text as parameters (an empty entry runs
without any), **Edit** opens the command dialog and **Delete** asks first.

Each command and history entry is a card panel: the header is the summary and the
opened body is indented from it, pairing each field label with its value in a
left-aligned two-column list, with the actions under a rule. Labels and values
are both selectable. The action row is spaced and coloured: **Run** green,
**Edit** blue, **Delete** red and **Restore this state** amber, so the risk of
each action reads before its label.

The edit dialog also changes the rest of the command: **Description**,
**Params** (comma-separated), **Cwd**, **Shell** and **Timeout**. Emptying a
field clears it. `Env` is not in the dialog — it holds values that look like
secrets, so it stays a configuration-file field, and the window shows keys only.
**Exec** and **Description** wrap and grow with their text (up to a cap); the
dialog's field list scrolls, so a long script never clips the buttons.

**Groups** on the group row opens one compact row per group — its emoji button,
its name and **Delete** — plus **New group**, which asks for a name and an icon
and stores them. Groups live in `groups.toml` next to the configuration (the
window is the only writer; the CLI ignores the file), so a group exists before a
command uses it and survives a restart. The filter then reads `🚀 deploy` instead
of a bare name. A command's Group dropdown lists the same groups and creates one
the same way. The dialog sizes itself from the window: it fits a small one and
does not stretch across a maximised one.

The **Setup** list beside it switches between the named setups the CLI creates
with `crower setup new`; **New setup** and **Delete setup** manage them without
leaving the window (the default and the active setup cannot be deleted). Each
entry says how many commands it holds — `work (12)` — so switching to an empty
setup is a choice rather than a surprise. Switching reloads
the commands, the history and the group icons without restarting, and saves the
choice the way `crower setup use` does — so the next `crower` run opens the same
setup.

Delete and restore both ask first. Restoring puts the command set back to what a
history entry recorded and drops the newer entries, through the same
`Session.RevertTo` path `crower revert` uses.

The History tab lists entries as cards, newest first. A card shows an operation
icon, the command, the time and the version; it opens into the labelled fields
(the full date and version included), the raw `key=value` note behind a toggle
with **Copy**, and **Run**, **Edit** and **Restore this state**. The raw note
sits on its own inset panel rather than as a stray line of text. Opening a card
expands it in place and only scrolls if it would otherwise leave the viewport.
Opening one closes the previous card. The newest hundred are rendered;
**Load older** appends the next hundred, so a long log opens quickly and only the
entries on screen are formatted. **Run** replays an execute entry with the
parameters it recorded; older entries without them ask when the command declares
parameters. **Edit** opens that command on the Commands tab.

Field labels and values in a command or history card, card descriptions, help
text, the raw note, tool results, and the output and error dialogs can be
selected and copied with Ctrl+C (⌘C on macOS), so a failure can be pasted into a
bug report. Information, picker and form dialogs close when you
click outside them; confirmation dialogs wait for an answer.

The About tab checks GitHub for a newer release when the window opens and
reports what it finds. When there is one, an **Install `<version>`** button
appears; nothing is installed without a confirmation, and the download goes
through the same checksum-verified path as `crower --upgrade`. Once a release is
on disk the install action gives way to **Restart now**, which closes the window
and replaces the running process with the new binary — on Linux and macOS it
takes over the same process rather than spawning a child a service manager could
reap with the old window. **Restore previous** appears only when the
installation recorded an earlier version. Neither the install nor the restart
button is a permanent control, so they never contradict the status line. If a
restart cannot start the new binary, `restart-error.log` is written beside
`crower.yaml`.

The Tools tab runs the file operations through the same use cases as the CLI
commands, and every run lands in the History tab. Path fields take a
comma-separated list, and **Browse…** beside each one fills it from a file or
folder picker — Fyne's dialog is single-select, so pressing Browse again appends
another path. Each card keeps its own result panel: it appears when that card
runs something and collapses with its toggle. The panel's left bar is green for
a success, amber for a stop and red for a failure, and its message is selectable. A copy, move or extract reports
how many items are done and offers **Cancel**, which stops the batch between
items — the platform command already running is not interrupted.

Dropping files or folders anywhere on the window appends their paths to the
**Sources** field of the copy/move card and opens the Tools tab. A drop only
fills the field; nothing runs until you press Copy or Move.

The window opens at the size it last closed with, recorded in `window.toml`
beside the configuration. A size too small to work in is ignored, so a corrupted
file cannot open an unusable window. Position is not recorded: Fyne's window API
does not expose it. A copy or move takes one or more
paths and a destination; a search takes a term, a folder (default `.`) and what
to look in; extract takes archive paths and an optional destination; open takes
paths and hands them to the platform's default application.

> [!NOTE]
> `crower-gui` needs a graphical session; in a terminal-only environment use
> `crower gui`. It ships in every platform archive next to `crower` and `cr`,
> not as a `crower` subcommand. The CLI binaries are cgo-free and need no
> display. On Windows, `crower-gui.exe` is the executable to double-click: it is
> a graphical program with the crow icon, while `crower.exe` and `cr.exe` stay
> terminal programs. It needs a driver that provides OpenGL 2.1; a start that
> cannot get one writes the reason to `fyne.log` beside `crower.yaml` and says so
> in a dialog. The window binary has no subcommands: it opens whatever the
> arguments are — `--config` and `--setup` are read, everything else is ignored —
> so a double-click can never land in a command parser.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
