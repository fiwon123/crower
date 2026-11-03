# Contributing

Thanks for wanting to help with crower. This project is small, and it needs more
kinds of help than patches — reading the list below is genuinely useful.

## Report an issue

Open an [issue](https://github.com/fiwon123/crower/issues/new) with:

- **the version** — `crower --version`;
- **the platform** — Linux, Windows, macOS, and the architecture (amd64/arm64);
- **what you expected, and what happened** — the exact command and the exact
  output helps more than a summary;
- **whether it reproduces from a release build** — a build from source proves
  something about your tree, not about what users download.

A screenshot is welcome for the TUI (`crower gui`) and desktop (`crower
desktop`) surfaces, and so is the text.

## Discuss

- **Issues** are for bugs and concrete proposals.
- **[Discussions](https://github.com/fiwon123/crower/discussions)** are for
  questions, ideas that are not yet a plan, and "would this be useful?".

## Help without writing code

Several claims in the release notes are **unproven**, because the machine that
wrote the code has none of the hardware they need. The README's
**[What is not proven](README.md#-what-is-not-proven)** section lists them, and
the wiki's **[Verifying claims](https://github.com/fiwon123/crower/wiki/Verifying-claims)**
page has the steps for each one and what to record.

Running the program and reporting exactly what happened — including "it worked"
— turns an unproven claim into a proven one. That is worth as much as a patch.

## Send code

- Open an issue first, so the change and the reason are recorded together.
- Branch from `dev` as `<type>/<issue>-<slug>` (`fix/42-scan-timeout`), and use
  conventional commits (`feat:`, `fix:`, `docs:`, `test:`, `chore:`, `ci:`).
- Run `make check` before pushing — gofmt, vet, tests and the build. A change
  that was not verified is not finished.
- One feature per pull request. Say in the body what ran and what did not, and
  do not describe a skipped check as a pass.
- `AGENTS.md` holds the detailed rules (branching, commits, the gate, release
  steps). It is written for the agent, but it is the same bar a human patch is
  held to.

## About the AI

This project is developed with an AI agent doing much of the implementation, and
that is not hidden: the commits say what was verified, the pull requests say
what was not, and a human decides what merges. What that means for you:

- **the review bar is the same** — a patch is judged by its tests and its
  evidence, not by who or what wrote it;
- **"it compiles" is not "it works"** — the repository's rule is that a skipped
  run is not a pass, and that applies to agent-written code first of all;
- **a human is on the other end** — issues, discussions and reviews are read by
  a person, and nothing is merged without one.
