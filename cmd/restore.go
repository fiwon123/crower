package cmd

import (
	"fmt"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/restorecommand"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore <command>",
	Short: "Restore a command from the history.",
	Long: `Restore a command from the history.

The command is looked up by name or alias in the snapshot of a past
operation, so a deleted or overwritten command can come back.

	crower restore greet
	crower restore hi`,
	Example: `	crower restore greet
	crower restore hi`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		request := restorecommand.Request{Args: args}
		if len(args) > 0 {
			request.Key = args[0]
		}

		session := newRestoreSession(cfg)

		return restorecommand.New(session, session, session, session).Run(request)
	},
}

// restoreSession adapts the configuration to the ports the restore use case
// declares.
type restoreSession struct {
	legacySession
}

// newRestoreSession wires the configuration.
func newRestoreSession(cfg *app.Config) restoreSession {
	return restoreSession{legacySession: newLegacySession(cfg)}
}

// HistoryRows returns the log newest first for the selection table.
func (s restoreSession) HistoryRows() []restorecommand.Entry {
	rows := make([]restorecommand.Entry, 0, len(s.cfg.History.AllData))

	for i := len(s.cfg.History.AllData) - 1; i >= 0; i-- {
		entry := s.cfg.History.AllData[i]
		rows = append(rows, restorecommand.Entry{
			Version: entry.Version,
			File:    entry.File,
			Note:    notes.Display(entry.Note),
		})
	}

	return rows
}

// FindCommand resolves the key inside the chosen snapshot.
func (s restoreSession) FindCommand(key, file string) (catalog.Command, error) {
	stored, err := s.cfg.FindCommand(key, history.Content{File: file})
	if err != nil {
		return catalog.Command{}, err
	}

	return catalog.Command{
		Name:    stored.Name,
		Aliases: stored.AllAlias,
		Exec:    stored.Exec,
		Steps:   stored.Steps,
	}, nil
}

// Save replaces the command set and writes the command file.
func (s restoreSession) Save(commands catalog.Catalog) error {
	s.cfg.ReplaceCatalog(commands)

	return s.cfg.SaveCommands()
}

// Warning writes a diagnostic line the user should see.
func (s restoreSession) Warning(message string, data ...any) {
	s.cfg.Logger.Warning(message, data...)
}

// HistoryTable renders the operation log.
func (s restoreSession) HistoryTable() string {
	return s.cfg.History.GetList()
}

// ChooseHistoryEntry asks which history row to restore.
func (s restoreSession) ChooseHistoryEntry() (restorecommand.Entry, error) {
	value, err := helper.GetUserInput("Select Row", helper.IsValidContentKey, s.cfg)
	if err != nil {
		return restorecommand.Entry{}, err
	}

	content, ok := value.(history.Content)
	if !ok {
		return restorecommand.Entry{}, fmt.Errorf("selected row is not a history entry")
	}

	return restorecommand.Entry{
		Version: content.Version,
		File:    content.File,
		Note:    notes.Display(content.Note),
	}, nil
}

// RestoreNote is the history note for the restore.
func (s restoreSession) RestoreNote(args []string) string {
	return notes.RestoreNote(args)
}
