# 🚚 `move`

**Move one or more files or folders to a destination.**

```bash
crower move <source...> <destination>
```

## 💡 Examples

```bash
crower move old.log ~/trash/
crower move a.txt b.txt ~/archive/
crower move ~/projects/demo ~/archive/
```

The last argument is the destination. Unlike [`copy`](copy), the source no longer
exists afterwards — and like every operation, the move is reversible.

## 🔗 Related

[`copy`](copy) · [`delete`](delete) · [`revert`](revert) · **[📖 Commands](How-to-Use)**
