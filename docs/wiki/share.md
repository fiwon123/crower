# 🤝 `share` · `export` · `import`

**Move commands between machines, or pull a pack from a registry.**

```bash
crower export [name] [-o file]
crower import <file> [--overwrite]
crower share list
crower share pull <pack>
```

## 🚩 Flags

| Command | Flag | Effect |
| --- | --- | --- |
| `export` | `-o`, `--out` | Output file path |
| `import` | `--overwrite` | Replace commands that already exist |

## 💡 Examples

```bash
crower export -o my-workflows.toml     # everything
crower export deploy -o deploy.toml    # one command
crower import my-workflows.toml
crower import team.toml --overwrite
crower share list
crower share pull docker-setup
```

## 🛰️ Registry protocol

Any HTTP server can host packs: `index.json` plus `packs/<name>.toml`. See
[crower-resources](https://github.com/fiwon123/crower-resources) for examples.

> [!WARNING]
> An exported file contains what you typed — hosts, paths, sometimes secrets.
> Read it before you commit or send it, and prefer `{1}`-style parameters over
> baked-in addresses.

## 🔗 Related

[🤝 Sharing](Sharing) · [⚙️ Configuration](Configuration) · **[📖 Commands](How-to-Use)**
