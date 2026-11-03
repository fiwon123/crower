# Verifying claims

## ⚡ In a hurry

Pick a claim you can test, follow its steps, and paste the result on its issue —
even "it worked" is the answer. Everything you need is a release archive.

Some claims in the release notes are **unproven**, because the machine that
wrote the code cannot run them: CI has no macOS box, no Windows box, no ARM
machine, and nothing that draws a window. Each one has an issue, and each issue
says in prose what would settle it.

This page is the missing half: the steps, the commands, and what to write down.

## Before you start

- Use a **release build** from
  [Releases](https://github.com/fiwon123/crower/releases), not a local one — a
  build from source proves something about your tree, not about what users
  download.
- Record the **version**: `crower --version`. Without it a report cannot be
  matched to a build.
- Note the **platform and architecture**: the archive name says both
  (`linux_amd64`, `darwin_arm64`, `windows_arm64`, …).

## What CI actually proves

`make check` (and the `CI (Go)` job) run `gofmt`, `go vet`,
`go test ./...` and `go build ./...` on **ubuntu-latest, amd64**. The release
workflow cross-compiles six archives with `CGO_ENABLED=0` on the same runner and
verifies each binary's embedded `vcs.revision`.

So a green build means: the code is formatted, vets clean, the tested packages
pass on Linux/amd64, and every archive contains a binary built from the release
commit. It does **not** mean any of these:

## The unproven claims

### 1. The macOS build runs — [#129](https://github.com/fiwon123/crower/issues/129)

The darwin archives are cross-compiled on Linux and never executed. There is an
open report that shell-backed operations silently do nothing on macOS.

**You need:** a Mac (Intel or Apple Silicon).

1. Download `crower_<version>_darwin_<arch>.tar.gz`, extract, add to `PATH`.
2. `crower --version`
3. `crower create hello "echo hi"` then `crower execute hello`
4. Try one operation that touches the shell: `crower find .`, `crower open .`

**Record:** the exact output of step 4, and whether `hello` printed `hi`.

### 2. The Windows build runs

The windows archives are cross-compiled on Linux and never executed.

**You need:** a Windows machine (amd64 or arm64).

1. Download `crower_<version>_windows_<arch>.zip`, extract, add to `PATH`.
2. `crower --version`
3. `crower create hello "echo hi"` then `crower execute hello`
4. `add_to_path.bat` — does it put `crower` and `cr` on `PATH` in a new terminal?

**Record:** whether the binary starts at all, and the output of step 3.

### 3. The arm64 builds run

`linux_arm64` and `windows_arm64` are built but never run.

**You need:** an ARM machine (a Pi, an ARM server, an ARM Windows device).

1. Run steps from claim 1 or 2 against the arm64 archive.
2. `crower --version`, then `crower create hello "echo hi"` + `crower execute hello`.

**Record:** `uname -m` (or `echo %PROCESSOR_ARCHITECTURE%`) and the output.

### 4. The TUI and desktop surfaces render

No gate draws a window. `go test ./...` asserts the model, not pixels.

**You need:** any machine with a terminal (`crower gui`) or a desktop session
(`crower desktop`).

1. `crower gui` — does the list render, do the keys move, does it exit cleanly?
2. `crower desktop` — does a window open, and does the list populate?

**Record:** a screenshot, plus anything clipped, misaligned, or unreadable.

### 5. History and revert across a real session

The tests use isolated temp directories; nobody has replayed a long real session.

1. Use crower for a day's file work: `create`, `copy`, `move`, `delete`, `note`.
2. `crower history`
3. `crower revert 3` — did the last three operations actually undo?
4. `crower restore <name>` for something deleted earlier.

**Record:** what `history` showed, and what `revert` changed on disk.

---

## What to report

Paste on the issue:

- the **version** and **platform/arch**;
- the **steps** you ran, verbatim;
- **what happened** — the exact message, not a summary;
- **"it worked"** if it worked. That turns an unproven claim into a proven one
  and is worth as much as reporting a defect.

A skipped step is not a verified step. If you could not finish the list, say
which step you stopped at rather than reporting the claim as proven.
