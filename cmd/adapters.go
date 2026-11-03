package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/helper"
)

// legacySession adapts the pre-refactor app.Config to the ports the new use
// cases declare.
//
// It lives in cmd because cmd is the composition root: the only place allowed
// to see the legacy core and the new rings at once. The adapter ring cannot
// hold it -- it would have to import internal/app and the legacy list
// renderer, which the layer gate forbids for reasons that still hold.
//
// The type disappears in Phase 5, when app.Config becomes a session constructed
// here and the real storage and output adapters are supplied instead.
type legacySession struct {
	cfg   *app.Config
	table func() string
}

// newLegacySession wires the configuration and the renderer the legacy core
// still owns.
func newLegacySession(cfg *app.Config, table func() string) legacySession {
	return legacySession{cfg: cfg, table: table}
}

// Commands is the current command set.
func (s legacySession) Commands() catalog.Catalog {
	return s.cfg.Catalog()
}

// Commit replaces the command set and records the operation that explains it.
//
// It is still two writes -- the command file, then the history -- because it
// wraps the legacy core; making that one unit of work belongs to the storage
// adapter in Phase 2.
func (s legacySession) Commit(commands catalog.Catalog, op operation.Operation, replay string) error {
	s.cfg.ReplaceCatalog(commands)

	if err := s.cfg.SaveCommands(); err != nil {
		return err
	}

	s.cfg.History.Add(op, "", note.New().
		AddMainOperation(op).
		AddCrowerExec(replay, nil).
		Build())

	return s.cfg.Save()
}

// Info writes a diagnostic line.
func (s legacySession) Info(message string) {
	s.cfg.Logger.Info(message)
}

// Table renders the current command table.
func (s legacySession) Table(catalog.Catalog) string {
	return s.table()
}

// Confirm asks a yes/no question on the terminal.
func (s legacySession) Confirm(question string) (bool, error) {
	return helper.GetUserConfirmation(question, s.cfg)
}
