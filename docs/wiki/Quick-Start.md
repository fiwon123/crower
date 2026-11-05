<div align="center">

# ⚡ Quick start

**Four commands, one minute.**

</div>

## 1️⃣ Create a shortcut

```bash
crower create deploy "ssh user@server 'cd /app && git pull && systemctl restart app'"
```

`create` takes a name and the command to run later.

## 2️⃣ Run it

```bash
crower execute deploy
cr deploy
```

Both forms work — `cr` is the same binary under a short name, and `execute` is
optional when the first argument is a stored name.

## 3️⃣ Add an alias

```bash
crower create backup -a bkp "tar -czf backup.tar.gz ~/important-data"
cr bkp
```

## 4️⃣ Pass parameters

```bash
crower create greet "echo Hello {1}, welcome to {2}!"
cr greet John Developer
```

`{1}`, `{2}`, … are positional; `{}` takes the rest. Arguments are quoted
before they reach the shell, so a value with a space stays one value.

---

## 🧭 What next

| | |
| --- | --- |
| 💬 **Don't want flags?** | `crower create` with no arguments asks, step by step |
| 📜 **See what you did** | `crower list -i` — history, newest first |
| ↩️ **Undo the last operation** | `crower revert 1` |
| 🗑️ **Bring something back** | `crower restore <name>` |
| 🎛️ **Prefer a UI?** | `crower gui` (terminal) or `crower-gui` (window) |
| 🤝 **Share your set** | `crower export -o team.toml` |

> [!TIP]
> Every command that touches the filesystem or your command list appends to the
> history, and every entry can be reverted. You do not have to be careful to be
> productive.
