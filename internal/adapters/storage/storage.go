// Package storage owns the bytes crower writes.
//
// Three files make up a user's installation: the TOML command map
// (`crower.yaml`, despite the name), the JSON history index
// (`historyhelper.json`), and one TOML snapshot per operation under `history/`.
// Their field names are a published format -- they are already in every user's
// home directory -- so the knowledge of them belongs in one place rather than
// spread over the loader, the saver and the revert path.
//
// The composition root (`internal/session`) constructs these types and calls
// them. An adapter never imports the root: the architecture gate enforces that
// direction, so the knowledge of the files stays below the session.
package storage
