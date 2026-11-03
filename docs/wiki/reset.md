# 🧹 `reset`

**Remove every stored command.**

```bash
crower reset
```

Asks for confirmation first. History and its snapshots are left alone, so a
mistake is still recoverable with [`revert`](revert) — but take an
[`export`](share) first if the commands matter.

```bash
crower export -o before-reset.toml
crower reset
```

> [!WARNING]
> This removes **all** commands, not one. For a single one use
> [`delete`](delete).

## 🔗 Related

[`delete`](delete) · [`revert`](revert) · [`share` / export](share) · **[📖 Commands](How-to-Use)**
