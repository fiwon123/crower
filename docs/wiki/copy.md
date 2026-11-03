# 📄 `copy`

**Copy one or more files or folders to a destination.**

```bash
crower copy <source...> <destination>
```

## 💡 Examples

```bash
crower copy file.txt ~/out/                 # one file
crower copy a.txt b.txt ~/out/              # several at once
crower copy ~/projects/demo ~/backup/       # a folder
```

The last argument is always the destination; everything before it is a source.

> [!NOTE]
> A copy is an operation like any other: it appears in
> [`list -i`](list) and can be undone with [`revert`](revert).

## 🔗 Related

[`move`](move) · [`delete`](delete) · [`extract`](extract) · **[📖 Commands](How-to-Use)**
