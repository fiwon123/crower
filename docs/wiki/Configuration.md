<div align="center">

# ⚙️ Configuration

**Everything lives under `~/crower/` unless you say otherwise.**

</div>

## 🗃️ Files

| Path | Contents |
| --- | --- |
| `~/crower/crower.yaml` | Your commands (TOML format) |
| `~/crower/historyhelper.json` | Operation history |
| `~/crower/history/` | Command snapshots, used by `revert` and `restore` |
| `~/crower/crower.log` | Execution log |

## 🔀 Another configuration

```bash
crower --config /path/to/crower.yaml list
```

Useful for a project-local set of commands, a CI job, or a second profile.

## 🖥️ Open instead of remembering

```bash
crower open              # the configuration file
crower open --folder     # the folder that holds it
crower open --system     # system variables
```

## 🔄 Upgrade

```bash
crower --check      # is there a newer release?
crower --upgrade    # self-update
```

`--upgrade` downloads the archive built for this OS and architecture, extracts
the `crower` binary and replaces the running one. If the release has no build
for this platform, it fails with the URL and the reason instead of applying
something else. Only the running binary is replaced: a separate `cr` copy keeps
the old version until it runs `--upgrade` itself.

> [!NOTE]
> Reverting an operation restores the *configuration snapshot* taken before it.
> History is append-only; `revert` walks it backwards, one step per call.
