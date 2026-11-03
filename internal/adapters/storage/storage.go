// Package storage owns the bytes crower writes.
//
// Three files make up a user's installation: the TOML command map
// (`crower.yaml`, despite the name), the JSON history index
// (`historyhelper.json`), and one TOML snapshot per operation under `history/`.
// Their field names are a published format -- they are already in every user's
// home directory -- so the knowledge of them belongs in one place rather than
// spread over the loader, the saver and the revert path.
//
// During the migration `internal/app` is still the session object and calls
// this package, so the arrow points from the legacy core to the storage. In
// Phase 5 the session is constructed by `cmd` and receives these types, and the
// arrow reverses to the direction the rings describe.
package storage
