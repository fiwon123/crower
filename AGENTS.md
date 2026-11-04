# AGENTS.md

> Working instructions for agents in this repository. `PROJECT_CONTEXT.md` holds
> the structured project identity behind this document; keep the two in step.

### Product

## Project Objective

crower — a cross-platform terminal file manager and command-alias manager that
keeps an undoable history and per-item notes.

### Problem

Doing routine file and shell work from a terminal means retyping the same
commands, losing the trail of what was changed, and having no way to attach a
note to a file. GUI-first file managers do not fit a terminal, and shell history
is neither searchable nor reversible nor annotatable.

### Goal

A user can run one short command to create, find, move, copy, open, annotate, or
revert a file or a command alias — and can always go back to the state before
the last operation.

### Main Technology

- **Language**: Go 1.25.4 (pinned in `go.mod`)
- **CLI framework**: spf13/cobra
- **Configuration**: TOML, under the user's home directory
- **Self-update**: minio/selfupdate
- **Development**: a local Go toolchain and `make`, or the dev sandbox — a pinned Go 1.25.4 container that runs the same on Linux, macOS, WSL2 and Windows. No external services.

### Core Architecture

```text
cmd/<operation>.go            cobra command: flags and argument parsing
  ↓
internal/application/<verb>/  use cases: the flow, with ports it declares itself
  ↓
internal/domain/              the model: catalog aggregate, operation enum
  ↓
internal/adapters/            storage, output, shell, prompt (the technology)

internal/data/ + internal/app/  storage shapes and wiring (Phase 5 work)
```

A use case is a package under `internal/application/`, and it declares the
interfaces it needs. The model under `internal/domain/` depends on the standard
library only. Every operation has been ported (#215): `internal/legacy/` is gone,
and the remaining `internal/data`/`internal/app` shapes are Phase 5 work.

`internal/shell` is the one thing the domains share below that line: the
platform's shell invocation (`sh -c` / `cmd /c`). A domain that needs to run a
shell line depends on it, never on another operation.

## Project Structure

```text
cmd/                 cobra commands, one file per operation
internal/domain/     the model: catalog aggregate, operation enum (stdlib only)
internal/application/ use cases, one package per verb, with their own ports
internal/adapters/   storage, output, shell, prompt — the technology side
internal/data/       history, notes, operations, command map, state
internal/app/        configuration and application wiring
internal/errors/     error types shared across layers
internal/helper/     shared helpers
internal/shell/      the platform's shell invocation, shared by every domain
pkg/logger/          logging
pkg/utils/           filesystem and process helpers
scripts/             add_to_path installers, the dev sandbox runtime, release-notes.sh
docs/wiki/           wiki source — synced by hand to the GitHub wiki
.github/             CI, release, branch gate, dependabot, PR and issue templates
CONTRIBUTING.md      the contributor contract (issue first, make check, evidence)
```

## Development Commands

```bash
make check           # the local gate: gofmt, vet, test, build
make release         # build the Windows and Linux release archives into build/
make clean           # remove build/
go run ./ --help     # run from source
go test ./...        # tests only
```

`make check` is the gate every agent runs before committing. It fails when the
toolchain is missing — a gate that cannot run reports nothing, and reporting
nothing must not look like reporting success.

### Dev sandbox — the same toolchain on every machine

```bash
make dev-up          # start the sandbox (foreground; Ctrl+C stops it)
make dev-shell       # a shell inside the sandbox
make dev-opencode    # run the agent inside the sandbox (Linux/WSL2 only)
make dev-down        # stop the sandbox
make help            # every target, with a one-line description
```

Without `make`, on any OS with Docker:

```bash
docker compose up    # the same sandbox, no Make, no Python, no host setup
```

On native Windows, where `make` does not exist: `python scripts/dev-up.py dev-up`.

The container pins Go 1.25.4 — the version in `go.mod` — so a build inside the
sandbox is the build CI runs. The sandbox is **optional**: a host toolchain
still works, and `make check` does not require Docker.

The image creates the directories the compose file mounts into, owned by the
runtime user: Docker creates a mount destination the image does not contain as
**root**, and the container runs as `${HOST_UID}`, so a sandbox whose image
omits them cannot write its own caches. A sandbox created before that fix
keeps its old root-owned volumes — Docker does not re-initialise a volume that
already exists — so it needs `make reset` once.

### Test coverage — read this before claiming a change is verified

The suite is small: many packages have no test file, and `go test` prints
`[no test files]` for them while still exiting 0. `make check` lists those
packages at the end of its output; anything named there has unverified
behaviour. A green run proves gofmt, vet, and build — and nothing about
behaviour.

There is **no known-failing baseline.** Every test gets its own fixture
directory from `internal/testutil.New(t)` (`t.TempDir`), so parallel
packages cannot race over a shared path. The failure that used to live in
`opendom`, `deletedom` and `movedom` was fixed in #47, which closed #46. A red
run is a regression until proven otherwise — do not discount it, and never
silence it by excluding a package from `go test` or from CI.

## Runtime Environment

The AI agent (opencode) runs at the repository root. There are two modes, and
they are not equivalent — say which one a claim was made in.

### Host mode (the default)

The agent runs directly on the host, against the tools installed on the
machine. This is what `make check` and `git` operate on.

### Sandbox mode (optional)

`make dev-opencode` runs the agent inside the dev container, where the
toolchain is pinned to Go 1.25.4 — the version in `go.mod`, the same one CI
uses. The workspace is bind-mounted at `/workspace`; edits made inside are the
same files on the host.

The container has **no Docker socket**: the agent inside it cannot manage
containers, and that is deliberate. `docker` commands fail there, and a claim
about the sandbox made from inside it cannot be about the host's Docker.

The in-container agent needs a **Linux** opencode binary (the container is
Linux). On macOS and native Windows, `make dev-opencode` refuses with a clear
message and the agent stays on the host; use WSL2, where the host binary *is*
Linux.

### What the agent CAN do

- Run Go commands (`go build`, `go test`, `go vet`, `go run`, `gofmt`)
- Run `make` targets
- Run git and GitHub CLI commands
- Read and edit any file in the repository except the secret files listed below
- Start and use the dev sandbox (`make dev-up`, `make dev-shell`) from the host

### What the agent CANNOT do

- Access `.env`, `.env.*`, `.envrc`, or any other secret file without explicit permission
- Run commands that dump environment values (`printenv`, `env`)
- Run `docker` from **inside** the sandbox (no socket is mounted)
- Push directly to `dev` or `main`
- Merge a pull request without explicit instruction

## Key Conventions

- **Assignee**: `fiwon123` — `workflow.md` tells agents to read `ASSIGNEE` from here before assigning issues and PRs
- **GitHub credentials**: `GH_TOKEN` or `GITHUB_TOKEN` — either name works, and neither is defined here. The values belong in the host environment; nothing is written into the repository.
- **Model assignments**: `opencode.json` assigns `opencode/big-pickle` to `build` and `opencode/mimo-v2.6-flash-free` to `plan`. That is a starting point, not a recommendation — no prompt depends on which model an agent runs. Verify a new id with `opencode models` first: a catalog entry on models.dev does not guarantee the installed build offers it, and an unresolvable id fails at spawn rather than at config load.
- **Repository link**: the README links to the GitHub repository, and the release page is linked from the README — a reader is always one click from the source and from the download.
- **Formatting**: `gofmt` — `make check` fails on an unformatted file
- **Testing**: `go test` — standard library only, no assertion framework
- **CI**: GitHub Actions runs gofmt, vet, test, and build on `dev` → `main` pull requests only — never on a feature-branch pull request
- **Local testing**: run `make check` before pushing a feature branch. CI is not the gate for feature work; the local gate is.
- **Version**: the binary reports the tag it was built from, injected with `-X` into `github.com/fiwon123/crower/cmd.Version`. A build without that variable fails.

## Environment and Secret-File Rules

- The application reads a TOML configuration from the user's home directory. It is not in the repository and must never be committed.
- Never access `.env`, `.env.*`, `.envrc`, or any other environment/secret file without explicit permission.
- Never open, read, print, summarize, quote, or send the contents of those files. `.env.example` is the one exception: it holds placeholders and is safe to read.
- Read `.env.example` with the **Read tool**, not with `cat`. Bash denies cover every `.env*` path on purpose.
- Never run commands that reveal environment values — `printenv`, `env`. `docker compose config` is on the deny list for the same reason: it renders the *interpolated* values.
- Never display, repeat, log, store, or include secret values in responses, code changes, or commits.
- To confirm a variable is set, check for presence, never print it: `test -n "$VAR" && echo set`.

### Enforcement and its limits

`opencode.json` enforces the rules above at the tool level. The top-level
`permission` block denies `read` and `edit` of `*.env`, `*.env.*`, and
`*.envrc` while allowing `*.env.example`, and denies `bash` commands matching
those paths or dumping the environment. Every agent inherits it, and each agent
repeats the deny list locally so the project stays guarded when one file is
edited in isolation.

Treat it as a backstop, not as the permission. The rules match **command
strings and file paths**, not intent, so they have known gaps you must not rely
on your way past:

- An allowed interpreter can still reach a file. `go run ./script.go` matches
  no deny pattern, because the pattern is matched against the command string,
  not what the command does.
- A repo-wide `grep -rn SECRET .` ignores `.gitignore`, so it will read a
  `.env` that is correctly excluded from version control. Prefer the Grep
  tool, which is gitignore-aware.
- Rules are evaluated with the **last matching pattern winning**, and
  agent-level rules take precedence over global ones. A broad allow placed
  after a narrow deny silently reinstates the leak — which is why every agent
  that allows `grep *` also denies `grep *.env`, `grep *.env.*`, and
  `grep *.envrc` explicitly, positioned *after* the allow.

## Development Workflow

### Git Workflow

- Branch from `dev` (not `main`)
- Never push directly to `dev` or `main`
- Branch naming: `<type>/<issue-number>-<slug>` (e.g. `feat/42-note-export`)
- Branch types: `feat/`, `fix/`, `refactor/`, `docs/`, `test/`, `chore/`, `ci/`
- Use conventional commits: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`, `ci:`
- All testing is local for feature branches — CI only runs on `dev` → `main` pull requests
- Do NOT merge pull requests unless explicitly instructed
- Do NOT automatically create release pull requests or merge to `main` — the user must explicitly request it
- Always return to `dev` after completing any merge
- Every change is tracked on GitHub: **issue → branch → PR → merge**

### Linking

- Every branch includes the issue number: `feat/42-note-export`
- Every pull request references its issue: `Closes #<number>`
- Every issue is assigned to a milestone before work begins
- Every pull request is assigned to a milestone when created
- No orphaned branches, pull requests, or issues

> A linked issue is closed automatically when its pull request merges into the
> repository's *default* branch. That branch is **`dev`** here, so a feature pull
> request carrying `Closes #N` closes its issue on merge. `gh issue close` is the
> fallback for when that did not happen — the PR omitted the clause, or the link
> was edited out — and reports `already closed` if it did. A `dev` → `main`
> release pull request closes nothing, because `main` is not the default branch.

### Branch Strategy

```text
main          ← release merges (CI runs here — dev→main PRs only)
  └── dev     ← integration branch, and the repository's default branch
        ├── feat/42-note-export        (all feature branches merge into dev)
        ├── fix/46-parallel-test-fixture-race
        └── chore/44-opencode-agent-layer
```

- Feature branches: `<type>/<issue>-<slug>` → merge to `dev`
- Release: `dev` → `main` (a pull request — CI runs there)
- Hotfixes: `fix/<issue>-<slug>` → merge to `dev`, cherry-pick to `main` if urgent

### Branch gate

`.github/workflows/branch-gate.yml` fails a pull request whose branch does not
follow the model above, so the author finds out before a reviewer does. It is a
**guard rail, not a lock**: an administrator can override a failed required
check, and a check nobody blocks on is noise. Branch protection (below) is the
control; the gate is the feedback. `dependabot/*` branches targeting `dev` are
exempt (they have no issue number), while `dependabot/*` targeting `main` and
all non-`dev`/non-`fix/<number>` branches targeting `main` remain blocked.

### Repository settings (configure by hand)

These live in the repository settings, not in a file, so no agent can set them.
Verify them once after a transfer.

**Branch protection** (Settings → Branches → Rules):

- `dev`: require a pull request, and require the **Branch gate** check. Do
  *not* require the CI check here — CI runs on release pull requests only, so a
  required check that never runs leaves every pull request waiting forever.
- `main`: require a pull request, and require both **CI (Go)** and **Branch
  gate**. The name must match the job's `name:` exactly — a required check that
  names a job id (`go`) instead of the job's name never matches, and every
  release pull request waits forever with no failing check to show for it.
  Do not allow deletions: `dev` is the head of the release pull request, so
  deleting it breaks every future pull request's base.

**Other settings worth turning on:**

- Secret scanning and push protection (Settings → Code security).
- Tag protection for `v*`, so a published release tag cannot be moved.
- Dependabot targeting `dev` — see `.github/dependabot.yml`.

### Milestones

- Every issue MUST be assigned to a milestone before work begins
- Milestones represent releases or sprint iterations
- Use `gh issue edit <issue-number> --milestone "<milestone-name>"`
- Track milestone progress on the GitHub Milestones page

### Project board

- **This repository's board**: <https://github.com/users/fiwon123/projects/5>.
- **Issues are the source of truth; the board is a view of them.** Every open
  issue belongs on the board, with its milestone and its status.
- Status values are **Todo**, **In Progress**, **Done** — the board's own
  `Status` field. Add columns only if the workflow genuinely needs them; a board
  that mirrors four states nothing updates is noise.
- Keep cards moving with the steps above: issue created → Todo, branch opened →
  In Progress, PR merged → Done. `gh project item-add 5 --owner fiwon123 --url <issue-url>`
  adds an issue; `gh project item-list 5 --owner fiwon123` reads the board.
- Do not duplicate the link graph in comments: the issue's branch, PR and
  milestone are already on it, and the board renders them.

### Wiki

- The wiki is sourced from `docs/wiki/` in this repository — edit the page there,
  not in the wiki UI, or the next sync overwrites it.
- The sync is a **human step** (the wiki is a second git repository; CI cannot
  push to it): see `docs/wiki/README.md` for the three commands.
- A page that ships in `docs/wiki/` and is not synced is a page readers cannot
  see. Run the sync in the same change that edits a page.

### Agents

Agent definitions live in `opencode.json` and `.opencode/agents/`:

- **build** (primary): Full development work with all tools enabled
- **plan** (primary): Analysis and planning without making changes — restricted to read-only subagents (explore, reviewer only)
- **reviewer** (subagent): Reviews code for bugs, security, and regressions

Each agent's instructions live in exactly one place: the body of
`.opencode/agents/<name>.md`. Model, mode, and permissions live in
`opencode.json`. Do not add a `prompt` key to `opencode.json` — it creates a
second prompt source that drifts from the markdown and leaves it undefined which
one wins.

An agent named in `opencode.json` without a matching markdown file has no
instructions at all. If you add an agent, add its `.md` in the same commit.

### Skills

- `github-workflow` ships in `.opencode/skills/github-workflow/SKILL.md` — the command cookbook for git, issues, and pull requests. `workflow.md` remains the policy authority; the skill is the quick reference.

### Session Memory (4-tier)

- At the START of every new chat: read `.opencode/memory/global.md` then
  `.opencode/memory/current.md` and resume from "Next action" — this is how
  long-running work survives between sessions. If the state references an
  issue/branch/PR, continue that work unless the user says otherwise.
- At meaningful milestones (issue created, branch renamed, code done, tests
  run, PR opened, merged) update `global.md`; update `current.md` after each
  step and at session end rotate the oldest archive out so only 3 remain.
- Durable, reusable lessons (a footgun, the command that works, a decision's
  rationale) go to `lessons.md`; offload them there when `global.md` approaches
  its cap rather than letting tracked work and lessons share one file.
- Keep sizes bounded: `global.md` ≤ ~200 lines, `lessons.md` ≤ ~100 lines,
  `current.md` ≤ ~100 lines, archive files ≤ ~50 lines each; never store
  secrets/tokens/credentials.
- The whole `.opencode/memory/` directory is gitignored and must never be
  committed. Full rules: `.opencode/instructions/memory.md`.

### Rules for Agents

All workflow rules that agents must follow are documented in this file:

- branch strategy and naming conventions;
- CI strategy — `.github/workflows/ci.go.yml` runs on release pull requests (`dev` → `main`) and manual dispatch, not on feature-branch pull requests, which are gated by local `make check`;
- milestone and issue conventions;
- development commands — `make check`, `make release`, `make clean`;
- environment and secret-file rules;
- the test-coverage notes above, and the rule that a skipped or empty run is not a pass.

### CI triggers (read this before claiming a change is verified)

CI runs on release pull requests (`dev` → `main`) and manual dispatch — never on
a feature-branch push. Two things can make a run verify nothing:

- a `paths:` filter that the change does not touch means the workflow does not
  trigger at all; and
- a job gated behind a label or `if:` condition comes back `skipped`, not
  `passed`.

**Skipped is not passed.** A triggered run whose jobs were all gated out, or a
run that never happened, is not coverage. Do not describe a change as "CI
passed" unless a run actually executed the job that covers it; if it did not,
say the change is unverified and test it locally.

The same rule applies locally: `go test` reporting `[no test files]` for a
package is not a pass for that package. `make check` prints those packages at
the end of every run; a green run that names them is green for the packages that
ran, not for those.

### Release workflow (read this before dispatching)

`.github/workflows/release.yml` is manual-dispatch and `main`-only. It tags
`main`, builds the six platform archives on `ubuntu-latest`, verifies each
binary's embedded `vcs.revision`, writes `checksums.txt`, and creates a
**draft** release — or publishes it when the `publish` input is set.

```bash
gh workflow run release.yml --ref main -f version=v0.6.0
```

- It refuses to build when `dev` holds commits `main` does not: a release ships
  what `main` contains, so merge `dev` into `main` first.
- **The changelog is built from the PRs merged into `dev`, not from
  `releases/generate-notes`.** Every feature PR is squash-merged, so its commits
  never become ancestors of the release commit and the generated notes see only
  the release merge. `scripts/release-notes.sh` bounds the list by the previous
  release and this commit, and groups by label.
- A published release is never rewritten; re-running the same version refreshes
  the draft. Bump the version instead of moving a published tag.
- The workflow itself is not covered by CI. Its first real proof is a run, and a
  run whose jobs never executed proves nothing — which is also why a build job
  stuck on a retired runner label used to leave no draft at all (#121).
