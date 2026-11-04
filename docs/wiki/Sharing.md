<div align="center">

# 🤝 Sharing

**Build once, share everywhere — onboarding, CI, or a teammate's laptop.**

</div>

## 📤 Export

```bash
crower export -o my-workflows.toml      # every command
crower export deploy -o deploy.toml     # one command
```

The file is TOML: readable, diffable, and safe to commit. A command may carry
`Description` and `Params` — what it is for and what its positional values mean;
both are metadata and never change how it runs.

It may also declare `Shell` (`sh`, `pwsh` or `cmd`; empty is the platform's
own) and `Cwd` (the folder it runs in, with `~` for the home directory). A pack
that says `Shell = "pwsh"` runs through PowerShell 7 on every platform, so the
reader needs `pwsh` installed; the declaration removes the guess about which
shell the command meant, not the requirement that it exists.

`Env` adds environment entries, and `Timeout` (`120s`) limits each step. An
`Env` value should be a `${VAR}` reference when it is a credential: a literal
under a secret-looking key is refused, because a pack is a file people commit
and share.

## 📥 Import

```bash
crower import my-workflows.toml
crower import my-workflows.toml --overwrite    # replace conflicts
```

Conflicts are reported by name; nothing is overwritten unless you ask.

## 🛰️ Pull from a registry

```bash
crower share list                 # browse available packs
crower share pull docker-setup    # install one
```

A registry is any HTTP server serving a manifest (`index.json`) plus
`packs/<name>.toml` — see
[crower-resources](https://github.com/fiwon123/crower-resources) for examples.

## 🔒 Before you share

> [!WARNING]
> An exported command contains what you typed — server names, paths, sometimes
> tokens, if the command carried one. Read the `.toml` before you commit it or
> send it.

- Prefer `{1}`-style parameters over baked-in hosts: `ssh {1}` shares, `ssh
  deploy@10.0.0.5` documents your network.
- Never share a command whose text contains a secret; rotate it if you already
  have.
