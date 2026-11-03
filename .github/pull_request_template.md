<!--
Keep the sections below — they are read at the moment a claim is written, which
is the only moment it can be checked. Delete this comment block before you open
the PR.
-->

Closes #<issue-number>

## What

<!-- What changed, in one or two sentences, and which packages or files carry it. -->

## Why

<!-- The problem this solves. Link the issue's reasoning rather than restating it. -->

## Verification

<!--
What actually ran, and what it proved. Be specific: a command and its exit
status is evidence; "tests pass" is a claim.

State the gate as it is, not more:

- `make check` covers gofmt, vet, `go test ./...`, and `go build ./...`.
- A package printed as `[no test files]` was NOT exercised — `make check` lists
  those packages at the end of the run; say so.
- CI runs on `dev` → `main` pull requests only. A feature PR may trigger no run
  at all, and a job gated out comes back `skipped`, not `passed`.
- Workflows are not covered by CI. A `.github/workflows/` change is unverified
  until a run executes the job that exercises it.
- The release archives are cross-compiled and never executed by any gate: a
  darwin/arm64 or windows build is a build claim, not a run claim.
- A skipped job, a run that never happened, or an empty test target is
  **unverified**, not passed.
-->

- [ ] `make check` — exit status:
- Untested packages touched (from the `[no test files]` list):
- Anything else that was run, and what it showed:

## Notes

<!--
- Follow-ups deliberately left out of scope, and the issue that tracks them.
- Known caveats a reviewer should not have to rediscover.
- Anything in the diff that looks wrong but is intentional, and why.
-->
