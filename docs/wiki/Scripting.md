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
