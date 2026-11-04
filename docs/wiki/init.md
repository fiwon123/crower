# 🚀 `init`

**Put `crower` and `cr` on `PATH`, and show the completion command.**

```bash
crower init            # show the line for your shell
crower init --apply    # append it to the startup file
```

`init` detects bash, zsh, fish or PowerShell, prints the startup file and the
line it needs, and prints the matching `crower completion <shell>` command.
`--apply` appends the line once; a second run is a no-op.

```bash
crower init
# Shell:   bash
# Startup: /home/me/.bashrc
# Add:     export PATH="$PATH:/home/me/.local/bin"
# Completion: crower completion bash
```

> [!NOTE]
> `--apply` changes a startup file. Read the printed line first if you keep your
> shell configuration under version control.

## 🔗 Related

[📦 Installation](How-to-Install) · [📖 Commands](How-to-Use)
