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
| **Commands** | The list, with a filter, a group filter and **Run**, **Edit** and **Delete** on every row; selecting one opens the **Details** panel below the list |
| **History** | Every recorded operation as a table, newest first; clicking a row selects the whole row and opens its details below the table, with **Restore this state** |
| **Tools** | Copy, move, search, extract and open; every path field has **Browse…**, and each card has its own result square that opens when an action runs |
| **About** | Logo, version, description, the repository link, and the update controls |

A command can carry an **icon** — a short marker such as 🚀 — shown beside its
name, and a **group**; the filter lists every group in use, and the details panel
shows which one the selected command belongs to. Both are set from the create
and edit dialogs, or with `crower create --icon/--group` and
`crower update --icon/--group`. The Icon field has a **Pick…** button with a grid
of about a hundred emoji; the field stays editable, so any other emoji can be
typed or pasted.

Selecting a row opens a **Details** panel below the list; the **Details** toggle
collapses it again. The panel shows the exec line and the command's metadata,
and carries a **params entry** with its own **Run**: the entry's text is passed
as parameters, and an empty entry runs the command without any. The row's ▶ and
the toolbar's **Run** are smart — a command that declares `Params` asks for them
first (an empty answer runs it without), and a command that declares none runs
straight away.

The edit dialog also changes the rest of the command: **Description**,
**Params** (comma-separated), **Cwd**, **Shell** and **Timeout**. Emptying a
field clears it. `Env` is not in the dialog — it holds values that look like
secrets, so it stays a configuration-file field, and the window shows keys only.

**Groups…** beside the filter gives each group an icon of its own, stored in
`groups.toml` next to the configuration (the window is the only writer; the CLI
ignores the file). The filter then reads `🚀 deploy` instead of a bare name.

The **Setup** list beside it switches between the named setups the CLI creates
with `crower setup new`. Switching reloads the commands, the history and the
group icons without restarting, and saves the choice the way
`crower setup use` does — so the next `crower` run opens the same setup.

Delete and restore both ask first. Restoring puts the command set back to what a
history entry recorded and drops the newer entries, through the same
`Session.RevertTo` path `crower revert` uses.

The About tab checks GitHub for a newer release when the window opens and
reports what it finds. **Install latest** and **Restore previous** always ask
first, then install through the same checksum-verified path as
`crower --upgrade`, and offer **Restart now** — the window starts the new binary
and closes.

The Tools tab runs the file operations through the same use cases as the CLI
commands, and every run lands in the History tab. Path fields take a
comma-separated list, and **Browse…** beside each one fills it from a file or
folder picker — Fyne's dialog is single-select, so pressing Browse again appends
another path. Each card keeps its own result square: it appears when that card
runs something and collapses with its toggle. A copy or move takes one or more
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
