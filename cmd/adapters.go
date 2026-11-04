package cmd

import (
	"fmt"

	"github.com/fiwon123/crower/internal/adapters/github"
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/output"
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
	cfg *app.Config
}

// newLegacySession wires the configuration.
func newLegacySession(cfg *app.Config) legacySession {
	return legacySession{cfg: cfg}
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

// Result writes data a caller may pipe -- a command's own output.
func (s legacySession) Result(text string) {
	s.cfg.Logger.Out(text)
}

// Debug writes a line only the debug log keeps.
func (s legacySession) Debug(message string, data ...any) {
	s.cfg.Logger.Debug(message, data...)
}

// Save replaces the command set and writes the command file.
//
// A use case that writes in two steps -- the command file, then the history --
// calls this and Record separately.
func (s legacySession) Save(commands catalog.Catalog) error {
	s.cfg.ReplaceCatalog(commands)

	return s.cfg.SaveCommands()
}

// ConfigFile is the path of the configuration file.
func (s legacySession) ConfigFile() string {
	return s.cfg.ConfigFile()
}

// ConfigFolder is the folder the configuration lives in.
func (s legacySession) ConfigFolder() string {
	return s.cfg.ConfigFolder()
}

// Table renders the command table from the catalog it is given.
func (s legacySession) Table(commands catalog.Catalog) string {
	return output.CommandsTable(commands)
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

// Pending renders the history entries a revert would drop and the ones that
// would remain.
func (s legacySession) Pending(steps int) (string, string, error) {
	if steps < 0 {
		return "", "", fmt.Errorf("steps is negative number")
	}

	stopIndex := s.cfg.History.GetIndexFromLastTo(steps)
	if stopIndex < 0 {
		return "", "", fmt.Errorf("can't revert: steps is greater than quantity of history registry")
	}

	return s.cfg.History.GetListLastHistory(steps), s.cfg.History.GetListFirstHistory(stopIndex), nil
}

// RevertTo restores the command set steps entries back and reports the version
// it landed on.
func (s legacySession) RevertTo(steps int) (int, error) {
	entry, err := s.cfg.History.GetBeforeLast(steps)
	if err != nil {
		return 0, err
	}

	if err := s.cfg.RevertTo(entry); err != nil {
		return 0, err
	}

	return entry.Version, nil
}

// Revert is the history note for a revert.
func (s legacySession) Revert(args []string) string {
	return notes.Revert(args)
}

// Check is the history note for --check.
func (s legacySession) Check() string {
	return notes.Check()
}

// UpgradeNote is the history note for a self-update.
func (s legacySession) UpgradeNote() string {
	return notes.UpgradeNote()
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
