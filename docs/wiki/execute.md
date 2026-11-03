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
`;` stays one value:

```bash
crower create find "grep -rn {1} {2}"
cr find "TODO" "internal/"      # quoted correctly, even with spaces
```

## 🔗 Related

[`create`](create) · [🧩 Scripting](Scripting) · **[📖 Commands](How-to-Use)**
