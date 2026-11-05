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

> [!NOTE]
> `crower-gui` needs a graphical session; in a terminal-only environment use
> `crower gui`. It ships in every platform archive next to `crower` and `cr`,
> not as a `crower` subcommand. The CLI binaries are cgo-free and need no
> display.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
