# ↩️ `revert`

**Undo history, one step or several.**

```bash
crower revert <steps>
```

## 💡 Examples

```bash
crower list -i        # see what is there, newest first
crower revert 1       # undo the most recent operation
crower revert 3       # undo the last three
```

Each entry holds a snapshot of the configuration taken **before** the operation,
so reverting restores the previous state rather than guessing.

> [!WARNING]
> Reverting is itself an operation. Revert the wrong number of steps and you can
> undo the undo — check `crower list -i` first.

## 🔗 Related

[`restore`](restore) · [`list`](list) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
