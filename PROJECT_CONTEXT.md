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
internal/session/             composition root: paths, loaded state, persistence
  ├─→ internal/application/<verb>/  use cases: the flow, with ports it declares
  └─→ internal/adapters/            storage, output, shell, prompt (the technology)
        ↓
internal/domain/              the model: catalog aggregate, operation enum
```

A use case is a package under `internal/application/`, and it declares the
interfaces it needs (consumer-owned ports). The model under `internal/domain/`
depends on the standard library only. Every operation has been ported (#215):
`internal/legacy/`, `internal/app`, `internal/data` and `internal/helper` are
gone; `internal/session` is the composition root.

`internal/shell` is the shared floor: the platform's shell invocation. A domain
depends on it to run `ls`, `cp` or `xdg-open`, never on another operation.

## Design Decisions

- **A domain package per operation** — each operation owns its own history, notes, and rules, so they can evolve independently. (2026-02)
- **TOML configuration under the user's home directory** — the config follows the user between machines and never lands in the repository. (2026-02)
- **Cobra with a flat command set** — `crower <operation>`, no nested verb groups, so every operation is discoverable from `--help`. (2026-02)
- **The `cr` short alias ships next to the `crower` binary** — a name typed many times a day should be short. (2026-02)
- **The release archive carries the README and the licenses** — a download has to be usable and lawful without the repository. (2026-02)
- **The shell runner is infrastructure, not an operation** — `internal/shell` owns `sh -c`/`cmd /c`, so `list`, `copy` or `open` do not import the `execute` operation to run a command. It stays the shared floor below every adapter rather than being absorbed into `adapters/system`: quoting and running a line is the one thing storage, output, notes, github, system and update would all otherwise duplicate. (2026-11)

## Topology

- kind: cli, path: ., purpose: the single binary — commands, domains, and data
- kind: library, path: pkg/, purpose: logger and filesystem utilities

There is one component: the repository root *is* the Go module.

## Directory Layout

```text
cmd/                 cobra commands, one file per operation
internal/domain/     the model: catalog aggregate, operation enum (stdlib only)
internal/application/ use cases, one package per verb, with their own ports
internal/adapters/   storage, output, shell, prompt — the technology side
internal/session/    composition root: paths, loaded state, persistence
internal/errors/     error types shared across layers
internal/shell/      the platform's shell invocation, shared by every domain
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

- **Done**: the hexagonal layering (#215) — `cmd` → `internal/session` →
  `internal/application` use cases → `internal/adapters`, with `internal/domain`
  on the standard library only. Every operation from the old `*dom` packages is
  ported, and `internal/legacy`, `internal/app`, `internal/data` and
  `internal/helper` are gone.
- **Done**: command-alias management (add, update, delete, reset, list, search,
  extract); file operations (create, copy, move, open, restore, revert); notes
  and history; setups; share/import; self-update; the CLI, TUI and Fyne window.
- **Done**: the Fyne window grew Commands, History, Tools and About tabs, icons
  and groups, setup switching, drag-and-drop, a remembered window size, the
  update flow, History Run/Edit, and a responsive detail/layout pass.
- **Done**: every package has a test, and the release workflow builds six
  targets — `linux/{amd64,arm64}`, `windows/{amd64,arm64}`,
  `darwin/{amd64,arm64}`.
- **Released**: v0.20.6 is the newest release; `dev` is the integration branch
  and the repository's default branch.

## Backlog

1. **The wiki sync is manual** — `docs/wiki/` is the source and the wiki is a
   second git repository, so a page edit needs the three commands in
   `docs/wiki/README.md` in the same change.
2. **`internal/fyneui` is outside `make check`** — every file is
   `//go:build fyne`, so the local gate neither compiles nor lists the package;
   the dispatched `fyne.yml` run is its only gate.
3. **Windows behaviour has no executing gate** — `cmd /c` quoting, paths and
   the process/script launches are build claims, never runs, on every CI job.

## Notes

- The `cr` alias is the same binary under a second name; both are shipped in
  the release archive.
- `internal/testutil.New(t)` gives each test its own directory under the
  system temporary directory (`t.TempDir`, removed when the test ends), so
  parallel packages cannot race over a shared fixture path — the failure that
  used to affect `opendom`, `deletedom` and `movedom` (#46, fixed in #47).
  It replaced `internal/helper.InitCrowerTests`, which used `os.MkdirTemp` and
  never cleaned up.
- Windows uses `cmd` to execute commands, Linux uses `sh`.
- The config file is TOML and lives under the user's home directory. It is
  never committed.
