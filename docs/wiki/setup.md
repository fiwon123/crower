# 🗂️ `setup`

**Keep separate command sets, each with its own history.**

```bash
crower setup list
crower setup new <name>
crower setup use <name>
crower setup show
crower setup delete <name> [--force]
```

Each setup lives under `~/crower/setups/<name>/` and owns its commands and its
history, so a revert in one never touches another. The default setup keeps the
original `~/crower/crower.yaml`; nothing migrates.

| Subcommand | Effect |
| --- | --- |
| `list` | the default setup and every named one, with the active one marked |
| `new <name>` | create an empty named setup; the name is a slug |
| `use <name>` | persist the default setup for future runs |
| `show` | the active setup, its configuration path and its command count |
| `delete <name>` | remove a named setup and its history; the active one needs `--force` |

## 🎯 Selecting a setup

`--setup <name>` selects one for a run; `CROWER_SETUP` sets the same default.
The order is `--config` > `--setup` > `CROWER_SETUP` > the persisted default
(`crower setup use`) > the default installation.

```bash
crower --setup streamer list
CROWER_SETUP=streamer crower list
```

## 🔗 Related

[⚙️ Configuration](Configuration) · [📖 Commands](How-to-Use)
