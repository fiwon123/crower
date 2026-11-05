# 🎛️ `gui` · `crower-gui`

**Two entries over the same configuration.**

| Command | Surface |
| --- | --- |
| `crower gui` | Interactive terminal UI (CLI archive) |
| `crower-gui` | Fyne window (GUI archive) |

```bash
crower gui
crower-gui
```

Both read the same `crower.yaml`, so a command created in the CLI appears in
the window, and anything done there lands in the history. The Fyne window can
create, edit and delete commands with the same catalog rules and history notes
as the CLI.

> [!NOTE]
> `crower-gui` needs a graphical session; in a terminal-only environment use
> `crower gui`. It is a separate executable in the `crower_gui_*` archive
> (Linux/Windows amd64, macOS arm64/Intel), not a `crower` subcommand. The CLI
> archive is cgo-free and has no window.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
