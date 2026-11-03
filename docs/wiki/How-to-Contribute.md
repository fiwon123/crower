<div align="center">

# 🤝 Contributing

**Patches are welcome. So is running the unproven parts.**

</div>

## 🐛 Report an issue

Open an [issue](https://github.com/fiwon123/crower/issues/new) with:

- **the version** — `crower --version`
- **the platform** — Linux, Windows, macOS, and the architecture
- **what you expected, and what happened** — the exact command and output beat a
  summary
- **whether it reproduces from a release build**

## 💬 Discuss

- **Issues** are for bugs and concrete proposals.
- **[Discussions](https://github.com/fiwon123/crower/discussions)** are for
  questions, half-formed ideas, and "would this be useful?".

## 🧪 Help without writing code

The most useful contribution right now is **running the unproven claims**:
**[🧪 Verifying claims](Verifying-claims)** lists each one, the steps, and what to
record. Reporting "it worked" turns an unproven claim into a proven one.

Documentation is the other gap: a wrong sentence here or in the README is a bug.

## 💻 Send code

- Open an issue first, so the change and the reason live together.
- Branch from `dev` as `<type>/<issue>-<slug>` (`fix/42-scan-timeout`), and use
  conventional commits.
- Run `make check` before pushing — it is also listed in **`CONTRIBUTING.md`**,
  which is the authoritative version.

**`CONTRIBUTING.md`** in the repository is the authoritative version of this page;
the **[README](https://github.com/fiwon123/crower)** has the user-facing docs.

## 🤖 About the AI

This project is developed with an AI agent doing much of the implementation, and
that is not hidden: commits say what was verified, pull requests say what was
not, and a human decides what merges. The review bar is the same — a patch is
judged by its tests and evidence, not by who wrote it.
