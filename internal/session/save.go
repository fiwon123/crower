package session

import (
	"fmt"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/shell"
)

// validateEnv refuses a command set the executor would refuse. The write path
// is where a pack or a hand-edited command would otherwise persist a literal
// under a secret-looking key; the runner's check alone runs too late.
func validateEnv(commands catalog.Catalog) error {
	for _, command := range commands.Commands() {
		if err := shell.ValidateEnv(command.Env); err != nil {
			return fmt.Errorf("command %q: %w", command.Name, err)
		}
	}

	return nil
}

// Save writes the current command map as a history snapshot and the history
// index next to the config file.
//
// The error is returned rather than logged: a run whose save failed did not
// persist anything, and the caller has to be able to say so with a non-zero
// exit status.
func (s *Session) Save() error {
	entry := s.History.GetLast()
	if entry == nil {
		return s.SaveOnlyHistory()
	}

	if err := s.historyFiles.WriteSnapshot(entry.File, s.allCommandsByName); err != nil {
		return err
	}

	return s.historyFiles.SaveIndex(s.History)
}

// Commit writes a new command set and the operation that produced it as one
// unit: the snapshot first, then the history index, then the command file.
//
// The order is the contract. A failure after the snapshot leaves an orphan
// snapshot and an index that still describes the old state; a failure after the
// index leaves the history ahead of the command file, which `revert` and
// `restore` can bring back. The command file is written last because it is the
// one copy the history can reconstruct.
//
// Each file is replaced through a temporary file and a rename, so a crash
// cannot leave a half-written file behind.
func (s *Session) Commit(commands catalog.Catalog, op operation.Operation, commandName, note string) error {
	if err := validateEnv(commands); err != nil {
		return err
	}

	s.ReplaceCatalog(commands)
	s.History.Add(op, commandName, note)

	entry := s.History.GetLast()
	if err := s.historyFiles.WriteSnapshot(entry.File, s.allCommandsByName); err != nil {
		return err
	}

	if err := s.SaveOnlyHistory(); err != nil {
		return err
	}

	return s.SaveCommands()
}

// SaveCommands writes the command map to the configuration file.
//
// Together with RevertTo (which writes the map it restores) this is the only
// writer of the configuration file; both live in this package, so no domain has
// to know that the file is TOML or where it is.
func (s *Session) SaveCommands() error {
	if err := validateEnv(s.Catalog()); err != nil {
		return err
	}

	return s.commands.Save(s.allCommandsByName)
}

// DeleteHistoryBackup removes the snapshot a history entry points at.
//
// file is a name from a history entry, not a path: the storage layer checks
// that it cannot escape the history folder, because a history file can be
// hand-edited and is not trusted.
func (s *Session) DeleteHistoryBackup(file string) error {
	return s.historyFiles.DeleteSnapshot(file)
}

// SaveOnlyHistory writes the history index without taking a snapshot.
func (s *Session) SaveOnlyHistory() error {
	return s.historyFiles.SaveIndex(s.History)
}
