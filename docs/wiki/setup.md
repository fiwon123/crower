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

## 🖥️ Switching from the window

`crower-gui` lists the same setups in a **Setup** control on the Commands tab.
Switching reloads the commands, history and group icons without a restart, and
writes the same persisted default `crower setup use` writes — the next terminal
command opens the setup the window chose.

`default` always means the original `~/crower/crower.yaml`, not a folder called
`default`.

Each entry in the control shows how many commands that setup holds, and a count
that cannot be read shows as `(?)` instead of hiding the setup. **New…** beside
the control creates an empty setup, and **Delete…** removes the one the control
points at, both with the rules the CLI uses: the name is a slug,
the default setup is not a folder, and a setup can only be deleted when it is not
the active one — switch first.

## 🔗 Related

[⚙️ Configuration](Configuration) · [📖 Commands](How-to-Use)
