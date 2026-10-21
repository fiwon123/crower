# Project Context

> Structured identity for this project. `AGENTS.md` is the working document an
> agent reads first; this file holds the facts behind it. Update both when the
> project changes.

---

## Identity

- **Project Name**: crower
- **One-liner**: A cross-platform terminal file manager and command-alias manager that keeps an undoable history and per-item notes.
- **GitHub repository**: https://github.com/fiwon123/crower

## Problem

Doing routine file and shell work from a terminal means retyping the same
commands, losing the trail of what was changed, and having no way to attach a
note to a file. Existing file managers are GUI-first; existing shell history is
not searchable, not reversible, and not annotatable.

## Goal

A user can run one short command to create, find, move, copy, open, annotate, or
revert a file or a command alias — and can always go back to the state before
the last operation.

## Stack

- **Language**: Go 1.25.4 (see `go.mod`)
- **CLI framework**: spf13/cobra
- **Configuration**: TOML (pelletier/go-toml), stored under the user's home directory
- **Logging**: custom logger in `pkg/logger`
- **Self-update**: minio/selfupdate
- **Dev environment**: a local Go toolchain and `make` — no containers, no external services

## Architecture

A cobra command is a thin shell over a *domain* package. Each domain owns its
own data, its operation history, and its notes, so a new operation is a new
package rather than a new branch in a shared handler.

```text
cmd/<operation>.go            cobra command: flags and argument parsing
  ↓
internal/domain/<op>dom/      core.go (the operation) + handler.go (the rules) + notes.go
  ↓
internal/data/                history, notes, operations, command map, state
  ↓
internal/app/                 configuration and application wiring
```

Each `internal/domain/<op>dom` package follows the same shape: `core.go` does
the work, `handler.go` validates and orchestrates it, and `notes.go` carries the
per-item documentation.

## Design Decisions

- **A domain package per operation** — each operation owns its own history, notes, and rules, so they can evolve independently. (2026-02)
- **TOML configuration under the user's home directory** — the config follows the user between machines and never lands in the repository. (2026-02)
- **Cobra with a flat command set** — `crower <operation>`, no nested verb groups, so every operation is discoverable from `--help`. (2026-02)
- **The `cr` short alias ships next to the `crower` binary** — a name typed many times a day should be short. (2026-02)
- **The release archive carries the README and the licenses** — a download has to be usable and lawful without the repository. (2026-02)

## Topology

- kind: cli, path: ., purpose: the single binary — commands, domains, and data
- kind: library, path: pkg/, purpose: logger and filesystem utilities

There is one component: the repository root *is* the Go module.

## Directory Layout

```text
cmd/                 cobra commands, one file per operation
internal/domain/     one package per operation (core, handler, notes)
internal/data/       history, notes, operations, command map, state
internal/app/        configuration and application wiring
internal/errors/     error types shared across layers
internal/helper/     shared helpers
pkg/logger/          logging
pkg/utils/           filesystem and process helpers
scripts/             add_to_path installers
```

## Key Commands

- **Build**: `go build ./...`
- **Run locally**: `go run ./ --help`
- **Tests**: `go test ./...`
- **Format check**: `gofmt -l .`
- **Vet**: `go vet ./...`
- **Full gate**: `make check`
- **Release archives**: `make release` (writes to `build/`)

## Status

- **Done**: command-alias management (add, update, delete, reset, list, search, extract); file operations (create, copy, move, open, restore, revert); notes and history; self-update; release archives for Windows and Linux.
- **In progress**: the test suite fails in three packages when it runs in parallel — the shared test fixture directory is wiped by every package that uses it (issue #46).

## Backlog

1. Fix the shared test fixture directory so the suite passes when the packages run in parallel (issue #46).
2. Add tests for the domains that have none — `executedom`, `extractdom`, `resetdom`, `restoredom`, `revertdom`, `upgradedom`, `checkdom`.
3. Extend the release matrix to macOS and Linux arm64.

## Notes

- The `cr` alias is the same binary under a second name; both are shipped in
  the release archive.
- `internal/helper/tests.go` creates and removes `$HOME/crower/test` when the
  test suite runs. It is outside the repository, and it is why running the
  suite touches the user's home directory — and why the packages race each
  other when they run in parallel (issue #46).
- Windows uses `cmd` to execute commands, Linux uses `sh`.
- The config file is TOML and lives under the user's home directory. It is
  never committed.
