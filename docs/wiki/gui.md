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

On the Commands tab, each command is a card: a name with a chevron that opens
into the labelled fields (name, exec, aliases, group, description, params, cwd,
shell, timeout and numbered steps), a parameter entry, and **Run**, **Edit** and
**Delete**. Opening one closes the previous card; a collapsed card carries no
actions. **Run** passes the entry's text as parameters (an empty entry runs
without any), **Edit** opens the command dialog and **Delete** asks first.
**New command** sits beside the filter.

Each command and history entry is a card panel: the header is the summary and the
opened body groups the fields in a label/value form, with the actions under a
rule. **Run** sits beside the parameters entry, so it is clear what it passes to
the command.

The edit dialog also changes the rest of the command: **Description**,
**Params** (comma-separated), **Cwd**, **Shell** and **Timeout**. Emptying a
field clears it. `Env` is not in the dialog — it holds values that look like
secrets, so it stays a configuration-file field, and the window shows keys only.
**Exec** and **Description** wrap and grow with their text (up to a cap); the
dialog's field list scrolls, so a long script never clips the buttons.

**Groups** on the group row gives each group an icon of its own through the
same catalogue picker, stored in `groups.toml` next to the configuration (the
window is the only writer; the CLI ignores the file). The filter then reads
`🚀 deploy` instead of a bare name. A group exists through the commands that use
it; create one from a command's Group dropdown.

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
(the full date and version included), the raw `key=value` line, and **Run**,
**Edit** and **Restore this state**. Opening a card expands it in place and only
scrolls if it would otherwise leave the viewport.
Opening one closes the previous card. The newest hundred are rendered;
**Load older** appends the next hundred, so a long log opens quickly and only the
entries on screen are formatted. **Run** replays an execute entry with the
parameters it recorded; older entries without them ask when the command declares
parameters. **Edit** opens that command on the Commands tab.

Values in a command or history card, card descriptions, help text, tool results,
and the output and error dialogs can be selected and copied with Ctrl+C (⌘C on
macOS), so a failure can be pasted into a bug report. Information, picker and form dialogs close when you
click outside them; confirmation dialogs wait for an answer.

The About tab checks GitHub for a newer release when the window opens and
reports what it finds. When there is one, an **Install `<version>`** button
appears; nothing is installed without a confirmation, and the download goes
through the same checksum-verified path as `crower --upgrade`. Once a release is
on disk the install action gives way to **Restart now**, which starts the new
binary and closes the window. **Restore previous** appears only when the
installation recorded an earlier version. Neither the install nor the restart
button is a permanent control, so they never contradict the status line.

The Tools tab runs the file operations through the same use cases as the CLI
commands, and every run lands in the History tab. Path fields take a
comma-separated list, and **Browse…** beside each one fills it from a file or
folder picker — Fyne's dialog is single-select, so pressing Browse again appends
another path. Each card keeps its own result square: it appears when that card
runs something and collapses with its toggle. A copy, move or extract reports
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
> display.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
