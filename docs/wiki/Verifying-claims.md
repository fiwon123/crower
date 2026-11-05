<div align="center">

# 🧪 Verifying claims

**Pick a claim you can test. Paste the result on its issue — "it worked" counts.**

</div>

Some claims in the release notes are **unproven**, because the machine that wrote
the code cannot run them: CI has no macOS box, no Windows box, no ARM machine,
and nothing that draws a window. Each has an issue, and each issue says what
would settle it.

## ✅ What CI proves

`make check` and the `CI (Go)` job run `gofmt`, `go vet`, `go test ./...` and
`go build ./...` on **ubuntu-latest, amd64**. The release workflow cross-builds
six archives with `CGO_ENABLED=0` on the same runner and checks each binary's
embedded `vcs.revision`.

So a green run means: formatted, vet-clean, the tested packages pass on
Linux/amd64, and each archive came from the release commit. It does **not** mean
any of the following.

## 📋 Before you start

- Use a **release build** from [Releases](https://github.com/fiwon123/crower/releases)
  — a build from source proves something about your tree, not about the download.
- Record the version: `crower --version`.
- Note the platform and architecture: the archive name has both.

## 🍎 1. The macOS build runs — [#129](https://github.com/fiwon123/crower/issues/129)

The darwin archives are cross-compiled on Linux and never executed. There is an
open report that shell-backed operations silently do nothing on macOS.

**You need:** a Mac (Intel or Apple Silicon).

1. Download `crower_<v>_darwin_<arch>.tar.gz`, extract, add to `PATH`.
2. `crower --version`
3. `crower create hello "echo hi"` then `crower execute hello`
4. Try one shell-backed operation: `crower find .`, `crower open .`

**Record:** the output of step 4, and whether step 3 printed `hi`.

## 🪟 2. The Windows build runs

The windows archives are cross-compiled on Linux and never started.

**You need:** Windows on amd64 or arm64.

1. Download `crower_<v>_windows_<arch>.zip`, extract, add to `PATH`.
2. `crower --version`
3. `crower create hello "echo hi"` then `crower execute hello`
4. `add_to_path.bat` — does a new terminal find `crower` and `cr`?

**Record:** whether the binary starts, and the output of step 3.

## 🧮 3. The arm64 builds run

`linux_arm64` and `windows_arm64` are built and never executed.

**You need:** an ARM machine (a Pi, an ARM server, an ARM Windows device).

Run the steps from claim 1 or 2 against the arm64 archive. **Record:**
`uname -m` (or `echo %PROCESSOR_ARCHITECTURE%`) and the output.

## 🖼️ 4. The TUI and Fyne GUI render

No gate draws a window; `go test ./...` asserts the model, not pixels.

**You need:** a terminal (`crower gui`) or a graphical session and the GUI archive (`crower-gui`).

1. `crower gui` — does the list render, do the keys move, does it exit cleanly?
2. `crower-gui` — does a window open, and does the list populate?

**Record:** a screenshot, plus anything clipped, misaligned or unreadable.

## 🕰️ 5. History and revert across a real session

The tests use isolated temp directories; nobody has replayed a long real session.

1. Use crower for a day: `create`, `copy`, `move`, `delete`, notes.
2. `crower list -i`
3. `crower revert 3` — did the last three operations undo?
4. `crower restore <name>` for something deleted earlier.

**Record:** what `list -i` showed, and what `revert` changed on disk.

## 📝 What to report

- the **version** and **platform/arch**;
- the **steps**, verbatim;
- **what happened** — the exact message, not a summary;
- **"it worked"**, if it did.

> [!IMPORTANT]
> A skipped step is not a verified step. If you stopped early, say which step —
> do not report the claim as proven.

---

🤝 [Contributing](How-to-Contribute) ·
🐛 [Issues](https://github.com/fiwon123/crower/issues)
