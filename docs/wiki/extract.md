# 🗜️ `extract`

**Unpack archives — `.zip`, `.tar`, `.tar.gz`, `.7z`, `.rar`, …**

```bash
crower extract <archive...> [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-o`, `--out` | Output folder (default: each archive's own folder) |

## 💡 Examples

```bash
crower extract files.zip
crower extract files.zip -o ~/out/
crower extract "downloads/*.zip"            # several at once
crower extract "downloads/*.zip" -o ~/out/
```

> [!TIP]
> A quoted glob is expanded by crower, so `"*.zip"` works where your shell would
> otherwise error on an unmatched pattern.

## 🔗 Related

[`create`](create) · [`copy`](copy) · [🚑 Troubleshooting](Troubleshooting) · **[📖 Commands](How-to-Use)**
