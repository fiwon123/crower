# ▶️ `execute`

**Run a stored command.**

```bash
crower execute <command> [params...]
```

`crower <command>` and `cr <command>` are the short forms. Use the explicit
`execute` when the name collides with a subcommand (`list`, `create`, …).

## 🚩 Flags

| Flag | Effect |
| --- | --- |
| `-l`, `--last` | Run the most recently executed command |
| `-c`, `--create` | Run the most recently created command |
| `-u`, `--update` | Run the most recently updated command |

## 💡 Examples

```bash
crower execute deploy
cr deploy                       # same thing
crower execute greet John dev   # {1} = John, {2} = dev
crower execute --last           # re-run what you just ran
```

## 🧠 Parameters

`{1}`, `{2}`, … are positional placeholders; `{}` takes the remaining arguments.
Every value is **quoted** before it reaches the shell, so a value with a space or
`;` stays one value.

A command with `Steps` gets the same substitution in every step. A step with no
placeholder is left as written, and passing parameters when no step has one is
an error rather than a silent drop.

```bash
crower create find "grep -rn {1} {2}"
cr find "TODO" "internal/"      # quoted correctly, even with spaces
```

## 🎯 Execution target

A command can say where and how it runs. The fields live in the configuration
file or in a pack; `create` and `update` do not have flags for them yet.

| Field | Meaning |
| --- | --- |
| `Shell` | `sh`, `pwsh` or `cmd`. Empty keeps the platform's own shell. |
| `Cwd` | Working directory; a leading `~` is the user's home. |
| `Env` | Entries merged over the process environment. |
| `Timeout` | A Go duration (`120s`, `2m`) applied to each step. |

An `Env` value that is exactly `${VAR}` reads the process environment; an unset
variable is an error rather than an empty value. A literal under a
secret-looking key (`*TOKEN*`, `*SECRET*`, `*PASSWORD*`, `*_KEY`, …) is refused,
because a shared pack carries placeholders, never secrets. An unknown shell, a
missing folder, a bad duration, or a secret literal fails before anything runs.

Every step of a `Steps` command inherits the whole target.

```toml
[deploy]
Description = "Ship the API"
Exec        = "rsync -az ./ {1}:/srv/api && ssh {1} 'systemctl restart api'"
Shell       = "sh"
Cwd         = "~/projects/api"
Env         = { STAGE = "prod", TOKEN = "${DEPLOY_TOKEN}" }
Timeout     = "120s"
Params      = ["host"]
AllAlias    = ["d"]
```

## 🔗 Related

[`create`](create) · [🧩 Scripting](Scripting) · **[📖 Commands](How-to-Use)**
