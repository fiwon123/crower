# 🔍 `search`

**Search files, folders, or the web.**

```bash
crower search <term> [path] [flags]
```

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-f`, `--file` | Match file names |
| `-o`, `--folder` | Match folder names |
| `-b`, `--browser` | Search the web in your default browser |

## 💡 Examples

```bash
crower search report ./            # files or folders matching "report"
crower search -f "*.log" ~/projects
crower search -o vendor .
crower search -b "golang context cancellation"
```

## 🔗 Related

[`list`](list) · [`open`](open) · **[📖 Commands](How-to-Use)**
