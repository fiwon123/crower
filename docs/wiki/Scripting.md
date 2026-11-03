<div align="center">

# 🧩 Scripting

**crower is meant to be called from a script as well as typed by hand.**

</div>

## 📤 Results and diagnostics are separate

- **Results go to stdout.** `crower list` writes the table and nothing else.
- **Diagnostics go to stderr.** Prompts, status lines, warnings and errors.

```bash
# Only the list reaches the file
crower list > commands.txt

# The result, with the status lines left behind
greeting="$(crower execute greet)"
```

## 🤖 Machine-readable results (`--json`)

`list` and `search` accept `--json`: the result becomes **one JSON object** on
stdout instead of a table. Diagnostics stay on stderr, and the exit status is
unchanged.

```bash
crower list --json | jq -r '.commands[].name'
crower list --folder ~/projects --json | jq -r '.entries[].name'
```

| Command | Object |
| --- | --- |
| `list` | `{"commands":[{"row":0,"name":"gs","aliases":["g"],"exec":"git status"}]}` |
| `list --history` | `{"history":[{"version":1,"file":"00001.yaml","timestamp":"2026-01-02T15:04:05Z","operation":"create","note":"…"}]}` |
| `list --process` | `{"processes":[{"pid":1,"name":"init","path":"/sbin/init"}]}` |
| `list --folder DIR` | `{"folder":"DIR","entries":[{"name":"a.txt","dir":false,"size":3,"mode":"-rw-r--r--","modified":"2026-01-02T15:04:05Z"}]}` |
| `list --system` | `{"variables":[{"name":"PATH","values":["/usr/local/bin","/usr/bin"]}]}` |
| `list --syspath` | `{"variable":{"name":"PATH","values":["/usr/local/bin","/usr/bin"]}}` |
| `search --file/--folder/(default)` | `{"search":{"kind":"file","query":"x","path":".","results":["./x.txt"]}}` |
| `search --browser` | `{"search":{"kind":"browser","query":"x","path":"","results":[]}}` |

- History is listed **newest first**, like the table.
- An empty collection is `[]`, never `null`.
- `--json` exists only on `list` and `search`. Any other command **rejects** it
  as an unknown flag instead of ignoring it — a flag that quietly does nothing
  is worse than a flag that is missing.
- `list --folder --json` reads the directory with Go instead of parsing
  `ls -la`, so it needs no shell arm. It is a real directory read, not a
  translation of the text output: the path is used literally (no `$VAR`
  expansion, no globbing), a symlink is followed, and a path that is not a
  directory is an error where `ls` would print the file.
- `search --json` splits the shell's output into lines: a filename containing a
  newline cannot come back as one result.
- `execute`, `extract` and `open` keep their output byte-for-byte; wrapping a
  command's own output in an envelope would break
  `crower execute build | grep error`.

## ✅ The exit status says whether the work happened

| Exit | Meaning |
| :--: | --- |
| `0` | It happened — or you answered **no** to a confirmation, because declining is a choice |
| `1` | It did not happen; the reason is one line on stderr, prefixed `error: ` |

```bash
crower --no-input delete old-command || exit 1
```

## 🤐 Never block on a prompt

With `--no-input`, or with stdin closed, a command that needs an answer **fails
with a message** instead of asking again:

```bash
if ! crower --no-input import team.toml; then
  echo "import needs a decision; run it by hand" >&2
  exit 1
fi
```

## 🗃️ A script's own command set

`--config` points a run at another configuration file, so a script does not touch
your personal list:

```bash
export CROWER_CI_CONFIG="$PWD/.ci-crower/crower.yaml"
crower --config "$CROWER_CI_CONFIG" create lint "golangci-lint run"
crower --config "$CROWER_CI_CONFIG" execute lint
```

## 🔁 In CI

```yaml
- name: Provision project commands
  run: |
    crower --config .ci/crower.yaml import .ci/commands.toml --overwrite
    crower --config .ci/crower.yaml --no-input execute build
```

> [!WARNING]
> A command that prompts will not wait under `--no-input` — it exits `1` with a
> message on stderr. Check the exit status; a command that "did nothing" is
> usually a declined prompt.

## 🧠 Parameter quoting

`{1}`, `{2}`, … are substituted and **quoted** before they reach the shell, so a
value containing spaces or `;` stays a value:

```bash
crower create find "grep -rn {1} {2}"
cr find "TODO" "internal/"
```
