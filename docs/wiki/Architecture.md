<div align="center">

# 🏛️ Architecture

**Four rings, one direction, and a gate that fails when it is broken.**

</div>

This page is for someone changing the code. A user does not need it; the
[command reference](Command-Reference) is for them.

---

## 🗺️ The shape

```text
cmd/                          driver adapter — cobra parses and delegates
  ↓
internal/session/             composition root — paths, loaded state, persistence
  ↓
internal/application/         use cases — the flow, one package per verb
  ↓
internal/domain/              the model — the catalog aggregate, the operation enum
  ↑
internal/adapters/            technology — storage, output, (shell, prompt)
```

The arrow points inward. A package in a ring may import the rings inside it,
never the ones outside.

| Ring | May import | Must not import |
| --- | --- | --- |
| `internal/domain` | its own subtree, the standard library | anything else |
| `internal/application` | `internal/domain` | `internal/adapters`, `internal/session`, `cmd` |
| `internal/adapters` | `internal/domain`, `internal/application`, the standard library, third-party code | `internal/session`, `cmd` |
| `internal/session` | `internal/adapters`, `internal/application`, `internal/domain`, the standard library | `cmd`, `internal/tui`, `internal/desktop` |
| `cmd` | everything | — |

`internal/session` is the composition root: it owns the configuration paths,
the loaded command and history state, and the save/commit/revert calls. `cmd`
is its only production importer; `internal/testutil` constructs one for tests.
The pre-refactor `internal/app`, `internal/data`, `internal/helper` and
`internal/interfaces` are gone (#215).

## 🚧 The gate

`architecture_test.go` parses the imports of every non-test file and applies one
rule per ring. It is not a convention in a document — a violation fails
`make check`:

```text
internal/application/… imports internal/adapters/…: the dependency arrow points the other way
```

It also reports which rings it examined, and fails a rule that is marked
required but matched no package — a gate that checks nothing must not read
green.

## 🧱 What is where

| Path | Holds |
| --- | --- |
| `internal/domain/catalog` | `Command`, `Catalog`: order, alias index, warnings, `Resolve`, `Update` |
| `internal/domain/operation` | the operation enum — a file format, not a label |
| `internal/application/<verb>commands` | one use case: the flow, with the ports it declares |
| `internal/adapters/storage` | the TOML command map, the JSON history index, the snapshots |
| `internal/adapters/output` | the streams, the result format, the JSON encoding |
| `internal/shell` | the platform's shell invocation and quoting |
| `internal/testutil` | the fixture harness every test starts from |
| `pkg/` | generic helpers: logging, filesystem, codecs |

## 🧩 Adding an operation

1. **Model.** If the operation changes a rule the aggregate owns (what a name
   may be, which alias wins), put it in `internal/domain/catalog` with tests. It
   takes no I/O and prints nothing.
2. **Use case.** `internal/application/<verb>commands/` — a struct with a `Run`,
   and the interfaces it needs declared in the same package:

   ```go
   type State interface {
       Commands() catalog.Catalog
       Commit(catalog.Catalog, operation.Operation, string) error
   }

   type Output interface{ Info(message string) }
   type Prompter interface{ Confirm(question string) (bool, error) }
   ```

   The consumer owns the interface, so the port set can only describe what this
   use case does.
3. **Adapter.** Implement those interfaces. If it wraps the legacy core, wire it
   in `cmd/adapters.go`; if it is real technology, it belongs in
   `internal/adapters/`.
4. **Command.** `cmd/<verb>.go` parses flags and delegates — it stays around 20
   lines.
5. **Tests.** Fakes for the use case, an end-to-end run through the cobra tree,
   and a mutation proof that the test fails when the behaviour is removed.
6. **Docs.** A `### crower <verb>` section in
   [Command-Reference](Command-Reference); the guard test fails without it.

## 💾 The files on disk

| File | Format | Contents |
| --- | --- | --- |
| `crower.yaml` | TOML | the command map; the name is historical |
| `historyhelper.json` | JSON | one entry per operation: version, file, timestamp, operation, command, note |
| `history/NNNNN.yaml` | TOML | the command map as it was after that operation |
| `crower.log` | text | a debug copy of every result and diagnostic |

The field names in these files are a contract with every existing
installation. `internal/adapters/storage` is the only package that knows them,
and `internal/session/persist_test.go` pins them.

## 🔐 Two decisions worth knowing

- **`cmd` is not behind a layer.** The release injects the version at build time
  with `-X github.com/fiwon123/crower/cmd.Version`; moving the package would
  break `Makefile`, `release.yml` and `internal/tools/archiveverify` together.
- **Ports are declared by the consumer**, not in a shared `ports` package. Go
  connects them structurally, so a use case can only describe what it actually
  needs — and a test can pass a fake with three methods.

## 🧪 Verifying a claim about this page

```bash
go test -run TestLayerDependencies ./   # the rings, and which were checked
go test ./...                            # behaviour
make check                               # gofmt, vet, tests, three platforms
```

A page like this can go stale. The gate is the part that cannot: if a rule here
stops being true, `make check` fails.
