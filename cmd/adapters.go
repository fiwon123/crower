package cmd

import (
	"fmt"

	"github.com/fiwon123/crower/internal/adapters/github"
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/app"
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
// One type implements every port because every port reads the same session; a
// type per port would only repeat the constructor. It disappears in Phase 5,
// when app.Config becomes a session constructed here and the real storage and
// output adapters are supplied instead.
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

// LastCommand is the command the last operation of a kind touched.
func (s legacySession) LastCommand(op operation.Operation) (string, bool) {
	entry := s.cfg.GetLast(op)
	if entry == nil || entry.CommandName == "" {
		return "", false
	}

	return entry.CommandName, true
}

// Commit replaces the command set and records the operation with its note.
//
// It is still two writes -- the command file, then the history -- because it
// wraps the legacy core; making that one unit of work belongs to the storage
// adapter in Phase 2.
func (s legacySession) Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error {
	s.cfg.ReplaceCatalog(commands)

	if err := s.cfg.SaveCommands(); err != nil {
		return err
	}

	s.cfg.History.Add(op, commandName, note)

	return s.cfg.Save()
}

// Info writes a diagnostic line.
func (s legacySession) Info(message string, data ...any) {
	s.cfg.Logger.Info(message, data...)
}

// Table renders the current command table.
func (s legacySession) Table(catalog.Catalog) string {
	return s.table()
}

// Confirm asks a yes/no question on the terminal.
func (s legacySession) Confirm(question string) (bool, error) {
	return helper.GetUserConfirmation(question, s.cfg)
}

// ChooseCommand asks which command to act on, by row key.
func (s legacySession) ChooseCommand() (string, error) {
	value, err := helper.GetUserInput("Select Row", helper.IsValidInputKey, s.cfg)
	if err != nil {
		return "", err
	}

	key, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("select row: expected a name, got %T", value)
	}

	return key, nil
}

// NewName asks for a required name.
func (s legacySession) NewName() (string, error) {
	return helper.InputName(s.cfg)
}

// NewAliases asks for aliases until the user enters an empty line.
func (s legacySession) NewAliases() ([]string, error) {
	return helper.InputAlias(s.cfg)
}

// NewExec asks for a required command line.
func (s legacySession) NewExec() (string, error) {
	return helper.InputExec(s.cfg)
}

// Record appends a history entry and persists the index.
//
// A read-only operation -- a check, an open -- changes no command, so it must
// not rewrite the command file; that is why this is separate from Commit.
func (s legacySession) Record(op operation.Operation, commandName, note string) error {
	s.cfg.History.Add(op, commandName, note)
	return s.cfg.Save()
}

// Check is the history note for --check.
func (s legacySession) Check() string {
	return notes.Check()
}

// Reset is the history note for a reset.
func (s legacySession) Reset() string {
	return notes.Reset()
}

// Update is the history note for an update.
func (s legacySession) Update(before, after catalog.Command, last operation.Operation) string {
	return notes.UpdateNote(before, after, last)
}

// latestTag adapts the github adapter's function to the Releases port, so cmd
// stays wiring and the port stays an interface.
type latestTag struct{}

func (latestTag) LatestTag() (string, error) {
	return github.LatestTag()
}
