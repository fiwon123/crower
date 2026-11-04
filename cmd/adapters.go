package cmd

import (
	"fmt"

	"github.com/fiwon123/crower/internal/adapters/github"
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/output"
	"github.com/fiwon123/crower/internal/adapters/prompt"
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/session"
)

// sessionPorts adapts the session to the ports the use cases declare.
//
// It lives in cmd because cmd is the composition root: the only place allowed
// to see the loaded session and the use-case ports at once. The adapter ring
// cannot hold it -- it would have to import internal/session, which the layer
// gate forbids.
//
// One type implements every port because every port reads the same session; a
// type per port would only repeat the constructor.
type sessionPorts struct {
	cfg *session.Session
}

// newSessionPorts wires the configuration.
func newSessionPorts(cfg *session.Session) sessionPorts {
	return sessionPorts{cfg: cfg}
}

// Commands is the current command set.
func (s sessionPorts) Commands() catalog.Catalog {
	return s.cfg.Catalog()
}

// LastCommand is the command the last operation of a kind touched.
func (s sessionPorts) LastCommand(op operation.Operation) (string, bool) {
	entry := s.cfg.GetLast(op)
	if entry == nil || entry.CommandName == "" {
		return "", false
	}

	return entry.CommandName, true
}

// Commit replaces the command set and records the operation with its note.
//
// It is one unit of work: the snapshot, the index, then the command file (see
// Session.Commit for why that order).
func (s sessionPorts) Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error {
	return s.cfg.Commit(commands, op, commandName, note)
}

// Info writes a diagnostic line.
func (s sessionPorts) Info(message string, data ...any) {
	s.cfg.Logger.Info(message, data...)
}

// Result writes data a caller may pipe -- a command's own output.
func (s sessionPorts) Result(text string) {
	s.cfg.Logger.Out(text)
}

// Debug writes a line only the debug log keeps.
func (s sessionPorts) Debug(message string, data ...any) {
	s.cfg.Logger.Debug(message, data...)
}

// Save replaces the command set and writes the command file.
//
// A use case that writes in two steps -- the command file, then the history --
// calls this and Record separately.
func (s sessionPorts) Save(commands catalog.Catalog) error {
	s.cfg.ReplaceCatalog(commands)

	return s.cfg.SaveCommands()
}

// ConfigFile is the path of the configuration file.
func (s sessionPorts) ConfigFile() string {
	return s.cfg.ConfigFile()
}

// ConfigFolder is the folder the configuration lives in.
func (s sessionPorts) ConfigFolder() string {
	return s.cfg.ConfigFolder()
}

// Table renders the command table from the catalog it is given.
func (s sessionPorts) Table(commands catalog.Catalog) string {
	return output.CommandsTable(commands)
}

// prompt builds the terminal prompter over the current state.
func (s sessionPorts) prompt() *prompt.Prompt {
	return prompt.New(s.cfg.ReadLine, s.cfg.Logger.Info, s.cfg.NoInput, s.cfg.Catalog(), s.cfg.History)
}

// Confirm asks a yes/no question on the terminal.
func (s sessionPorts) Confirm(question string) (bool, error) {
	return s.prompt().Confirm(question)
}

// ChooseCommand asks which command to act on, by row key.
func (s sessionPorts) ChooseCommand() (string, error) {
	return s.prompt().ChooseCommand()
}

// NewName asks for a required name.
func (s sessionPorts) NewName() (string, error) {
	return s.prompt().Ask("Add new name: ")
}

// NewAliases asks for aliases until the user enters an empty line.
func (s sessionPorts) NewAliases() ([]string, error) {
	return s.prompt().Aliases()
}

// NewExec asks for a required command line.
func (s sessionPorts) NewExec() (string, error) {
	return s.prompt().Ask("Add new exec: ")
}

// Record appends a history entry and persists the index.
//
// A read-only operation -- a check, an open -- changes no command, so it must
// not rewrite the command file; that is why this is separate from Commit.
func (s sessionPorts) Record(op operation.Operation, commandName, note string) error {
	s.cfg.History.Add(op, commandName, note)
	return s.cfg.Save()
}

// Pending renders the history entries a revert would drop and the ones that
// would remain.
func (s sessionPorts) Pending(steps int) (string, string, error) {
	if steps < 0 {
		return "", "", fmt.Errorf("steps is negative number")
	}

	stopIndex := s.cfg.History.GetIndexFromLastTo(steps)
	if stopIndex < 0 {
		return "", "", fmt.Errorf("can't revert: steps is greater than quantity of history registry")
	}

	return output.HistoryLast(s.cfg.History, steps), output.HistoryFirst(s.cfg.History, stopIndex), nil
}

// RevertTo restores the command set steps entries back and reports the version
// it landed on.
func (s sessionPorts) RevertTo(steps int) (int, error) {
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
func (s sessionPorts) Revert(args []string) string {
	return notes.Revert(args)
}

// Check is the history note for --check.
func (s sessionPorts) Check() string {
	return notes.Check()
}

// UpgradeNote is the history note for a self-update.
func (s sessionPorts) UpgradeNote() string {
	return notes.UpgradeNote()
}

// Reset is the history note for a reset.
func (s sessionPorts) Reset() string {
	return notes.Reset()
}

// Update is the history note for an update.
func (s sessionPorts) Update(before, after catalog.Command, last operation.Operation) string {
	return notes.UpdateNote(before, after, last)
}

// updateSession adds the installed-version list to the shared session.
//
// It exists because the version history is its own file, not the command set:
// the update use case asks it what ran before, and records what it installed.
type updateSession struct {
	sessionPorts

	versions *storage.Versions
}

// newUpdateSession wires the configuration and the version list.
func newUpdateSession(cfg *session.Session) updateSession {
	return updateSession{
		sessionPorts: newSessionPorts(cfg),
		versions:     storage.NewVersions(cfg.ConfigFolder()),
	}
}

// Previous returns the most recently installed version other than current.
func (s updateSession) Previous(current string) (string, bool) {
	return s.versions.Previous(current)
}

// RecordInstalled puts version at the front of the installed list.
func (s updateSession) RecordInstalled(version string) error {
	return s.versions.RecordInstalled(version)
}

// latestTag adapts the github adapter's function to the Releases port, so cmd
// stays wiring and the port stays an interface.
type latestTag struct{}

func (latestTag) LatestTag() (string, error) {
	return github.LatestTag()
}
