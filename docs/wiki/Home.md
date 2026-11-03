# 🐦 crower

**A cross-platform terminal file manager and command-alias manager.**
Shortcut any command, automate workflows, and share them with your team.

[![CI](https://github.com/fiwon123/crower/actions/workflows/ci.go.yml/badge.svg?branch=main)](https://github.com/fiwon123/crower/actions/workflows/ci.go.yml)
[![Release](https://img.shields.io/github/v/release/fiwon123/crower)](https://github.com/fiwon123/crower/releases)
[![License](https://img.shields.io/github/license/fiwon123/crower)](https://github.com/fiwon123/crower/blob/main/LICENSE)

---

## 🚀 Start in thirty seconds

```bash
# create a shortcut for a long command
crower create "docker-ps" "docker ps --format 'table {{.Names}}\t{{.Status}}'"

# run it by name, or by its alias
crower execute docker-ps
cr docker-ps
```

Download the archive for your platform from
[Releases](https://github.com/fiwon123/crower/releases), extract it, and put
`crower` and `cr` on your `PATH` — the archive ships
`add_to_path.sh` / `add_to_path.bat` for that.

## 🎯 Pick your path

| I want to… | Go here |
| --- | --- |
| **install it** | [README — Installation](https://github.com/fiwon123/crower#-installation) |
| **see real recipes** | [README — Usage Recipes](https://github.com/fiwon123/crower#-usage-recipes) |
| **look up a command** | [README — Commands](https://github.com/fiwon123/crower#-commands) |
| **script against it** | [README — Scripting](https://github.com/fiwon123/crower#-scripting) |
| **build it** | [README — Development](https://github.com/fiwon123/crower#-development) |
| **report a bug** | [Issues](https://github.com/fiwon123/crower/issues) |
| **help prove a claim** | [Verifying claims](Verifying-claims) |

## ⚠️ What is not proven

This project says what it verified. The short version: CI compiles every
platform and tests on Linux/amd64, and no gate opens a GUI or runs an archive it
built. See **[Verifying claims](Verifying-claims)** for the list and the steps
that would settle each one.

---

📖 [README](https://github.com/fiwon123/crower) · 🐛 [Issues](https://github.com/fiwon123/crower/issues) · 🎯 [Milestones](https://github.com/fiwon123/crower/milestones) · 🤝 [Contributing](https://github.com/fiwon123/crower/blob/dev/CONTRIBUTING.md)
