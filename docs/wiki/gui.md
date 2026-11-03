# 🎛️ `gui` · `desktop`

**Two surfaces over the same configuration.**

| Command | Surface |
| --- | --- |
| `crower gui` | Interactive terminal UI |
| `crower desktop` | Native window (cgo-free) |

```bash
crower gui
crower desktop
```

Both read the same `crower.yaml`, so a command created in the CLI appears in
both, and anything done there lands in the history.

> [!NOTE]
> `desktop` needs a graphical session; in a terminal-only environment use `gui`.
> Neither surface is exercised by CI — see
> **[🧪 Verifying claims](Verifying-claims)**.

## 🔗 Related

[⚙️ Configuration](Configuration) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
