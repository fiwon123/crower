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
| `import` | `--overwrite` | Replace a command with the same name; drop aliases that conflict |

## 💡 Examples

```bash
crower export -o my-workflows.toml     # everything
crower export deploy -o deploy.toml    # one command
crower import my-workflows.toml
crower import team.toml --overwrite
crower share list
crower share pull docker-setup
```

## 🔀 Name and alias collisions

Import never lets two commands share a name or an alias. If a pack entry's name
or alias is already taken:

- **without `--overwrite`**, crower asks once per command; accepting renames the
  conflicting name or alias to `x_imported` (numbered if needed), declining
  skips that command;
- **with `--overwrite`**, a command with the same *name* is replaced, and a
  conflicting *alias* is dropped — taking it would silently edit the command
  that owns it.

Entries are processed in name order. Within one pack, a later entry that
collides with an earlier one is renamed the same way, even with `--overwrite`,
so two entries can never share a key.

`share pull` uses the non-overwrite path after its own confirmation.

## 🛰️ Registry protocol

Any HTTP server can host packs: `index.json` plus `packs/<name>.toml`. See
[crower-resources](https://github.com/fiwon123/crower-resources) for examples.

> [!WARNING]
> An exported file contains what you typed — hosts, paths, sometimes secrets.
> Read it before you commit or send it, and prefer `{1}`-style parameters over
> baked-in addresses.

## 🔗 Related

[🤝 Sharing](Sharing) · [⚙️ Configuration](Configuration) · **[📖 Commands](How-to-Use)**
