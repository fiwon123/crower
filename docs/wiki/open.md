# 📂 `open`

**Open a file, a folder, or your system settings.**

```bash
crower open [path...] [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| *(none)* | Open the configuration file |
| `-f`, `--file` | Open the configuration file, or the files you pass |
| `-o`, `--folder` | Open the configuration folder, or the folders you pass |
| `-s`, `--system` | Open system settings |

## 💡 Examples

```bash
crower open                                  # the configuration file
crower open --file notes.md
crower open --folder ~/projects
crower open --system
```

> [!TIP]
> `crower open` is the quickest way to see what a stored command actually
> contains — the file is plain TOML.

## 🔗 Related

[⚙️ Configuration](Configuration) · [`update`](update) · **[📖 Commands](How-to-Use)**
