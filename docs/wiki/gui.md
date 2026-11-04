# 🎛️ `gui` · `desktop` · `fyne`

**Three surfaces over the same configuration.**

| Command | Surface |
| --- | --- |
| `crower gui` | Interactive terminal UI |
| `crower desktop` | Native window (cgo-free) |
| `crower fyne` | Fyne window (separate CGO archive) |

```bash
crower gui
crower desktop
crower fyne
```

All three read the same `crower.yaml`, so a command created in the CLI appears
in each, and anything done there lands in the history. The Fyne window can
create, edit and delete commands with the same catalog rules and history notes
as the CLI.

> [!NOTE]
> `desktop` needs a graphical session; in a terminal-only environment use `gui`.
> `fyne` exists only in the `crower_fyne_<version>_linux_amd64` and
> `crower_fyne_<version>_windows_amd64` archives; the default archive keeps the
> cgo-free `desktop`. The windows are not exercised by CI — see
> **[🧪 Verifying claims](Verifying-claims)**.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
