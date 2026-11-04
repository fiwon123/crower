<div align="center">

# 🛠️ Development

**`make check` is the gate a change has to pass.**

</div>

## ⚙️ Commands

```bash
make check          # gofmt, vet, go test ./..., go build ./..., windows/darwin build
go test ./...       # tests only
make release        # cross-build the release archives into build/dist/
make help           # every target, one line each
```

`make check` fails when the toolchain is missing — a gate that cannot run reports
nothing, and reporting nothing must not look like success. It also lists the
packages that have **no test file**: `go test` prints `[no test files]` for them
and exits 0, so a green run is green for the packages that ran.

It builds for `windows/amd64` and `darwin/arm64` as well as the host. Several
packages pick their behaviour with build tags (`handler_windows.go`,
`handler_darwin.go`, …), and a host build on Linux never compiles those files, so
without the cross builds a Windows-only break would first surface in a release.

## 🐳 Dev sandbox

```bash
make dev-up         # the pinned toolchain, foreground
make dev-shell      # a shell inside it
make dev-opencode   # the agent inside it (Linux/WSL2)
make dev-down       # stop it
```

On any OS with Docker and no `make`: `docker compose up`. On native Windows:
`python scripts/dev-up.py dev-up`.

## 🏛️ Layout

Four rings, one dependency direction, and a gate in `architecture_test.go` that
fails the build when an import points the wrong way:

```text
cmd/  ->  internal/application/  ->  internal/domain/
                 ^
        internal/adapters/
```

**[Architecture](Architecture)** explains the rule, what each ring may import,
and where a new operation goes. Every operation has been ported (#215): there is
no `internal/legacy/` left to shrink.

## 🌿 Branches

```text
main          ← release merges (dev → main, a pull request)
  └── dev     ← integration branch and the repository's default branch
        ├── feat/42-note-export
        ├── fix/46-test-fixture-race
        └── chore/44-agent-layer
```

`<type>/<issue>-<slug>`, conventional commits, one feature per pull request.
`AGENTS.md` holds the full rules.

## 🚀 Releasing

```bash
gh workflow run release.yml --ref main -f version=v0.6.0
```

The workflow tags `main`, builds six platform archives, verifies each binary's
embedded revision, writes `checksums.txt`, and creates a **draft** release —
publish it by hand, or pass the `publish` input.

The changelog is built by `scripts/release-notes.sh` from the pull requests
merged into `dev`, **not** from `generate-notes`: feature commits are
squash-merged, so the commit graph the generator walks cannot see them and every
shipped change would be missing from the notes.

## 🧪 What a green run proves

Formatting, vet, the tested packages on Linux/amd64, and that every archive
contains a binary built from the release commit. It does **not** run the darwin,
windows or arm64 binaries, and it does not draw a window. See
**[Verifying claims](Verifying-claims)**.

---

🐛 [Issues](https://github.com/fiwon123/crower/issues) ·
🎯 [Roadmap board](https://github.com/users/fiwon123/projects/5) ·
🤝 [Contributing](How-to-Contribute)
