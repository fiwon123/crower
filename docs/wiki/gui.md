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
| **Commands** | The list, with a filter, a group filter, a detail line, and **Run**, **Run with params**, **Edit** and **Delete** on every row |
| **History** | Every recorded operation as a table, newest first, with a **Restore this state** action |
| **Tools** | Copy, move, search, extract and open, with their result below |
| **About** | Logo, version, description, the repository link, and the update controls |

A command can carry an **icon** — a short marker such as 🚀 — shown beside its
name, and a **group**; the filter lists every group in use, and the detail line
shows which one the selected command belongs to. Both are set from the create
and edit dialogs, or with `crower create --icon/--group` and
`crower update --icon/--group`.

**Groups…** beside the filter gives each group an icon of its own, stored in
`groups.toml` next to the configuration (the window is the only writer; the CLI
ignores the file). The filter then reads `🚀 deploy` instead of a bare name.

Delete and restore both ask first. Restoring puts the command set back to what a
history entry recorded and drops the newer entries, through the same
`Session.RevertTo` path `crower revert` uses.

The About tab checks GitHub for a newer release when the window opens and
reports what it finds. **Install latest** and **Restore previous** always ask
first, then install through the same checksum-verified path as
`crower --upgrade`, and offer **Restart now** — the window starts the new binary
and closes.

The Tools tab runs the file operations through the same use cases as the CLI
commands, and every run lands in the History tab. A copy or move takes one path
per line and a destination; a search takes a term, a folder (default `.`) and
what to look in; extract takes archive paths and an optional destination; open
takes paths and hands them to the platform's default application.

> [!NOTE]
> `crower-gui` needs a graphical session; in a terminal-only environment use
> `crower gui`. It ships in every platform archive next to `crower` and `cr`,
> not as a `crower` subcommand. The CLI binaries are cgo-free and need no
> display.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
