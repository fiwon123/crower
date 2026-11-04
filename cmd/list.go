package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/output"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/listcontent"

	"github.com/spf13/cobra"
)

var historyFlag bool

// Cmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "List commands, history, processes, folders or variables.",
	Long: `List commands, history, processes, folders or variables.

With no flag, crower lists the managed commands in the order they were
saved.

	crower list
	crower list --history
	crower list --process
	crower list --folder ~/projects
	crower list --system
	crower list --syspath

--json prints one JSON object on stdout instead of a table.`,
	Example: `	crower list
	crower list --history
	crower list --folder ~/projects
	crower list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		request := listcontent.Request{Args: args}

		switch {
		case processFlag:
			request.Mode = listcontent.ModeProcess
		case historyFlag:
			request.Mode = listcontent.ModeHistory
		case folderFlag:
			request.Mode = listcontent.ModeFolder
		case systemFlag:
			request.Mode = listcontent.ModeSystem
		case sysPathFlag:
			request.Mode = listcontent.ModePath
		default:
			request.Mode = listcontent.ModeCommands
		}

		session := newListSession(cfg)

		return listcontent.New(session, session, session, session, session, session, session).Run(request)
	},
}

// listSession adapts the configuration to the ports the list use case declares.
//
// It embeds legacySession for the catalog, the command table, the diagnostics
// stream and the history entry, and adds what only list needs: the result
// stream and its format, the history rows, and the platform listings.
type listSession struct {
	legacySession

	listings *system.Listings
}

// newListSession wires the configuration and the platform listings.
func newListSession(cfg *app.Config) listSession {
	return listSession{
		legacySession: newLegacySession(cfg),
		listings:      system.NewListings(),
	}
}

// Emit writes the text, or the payload when --json was asked for.
func (s listSession) Emit(text string, payload any) error {
	return s.cfg.Result(text, payload)
}

// Structured reports whether the result is requested as JSON.
func (s listSession) Structured() bool {
	return s.cfg.Format == app.FormatJSON
}

// HistoryTable renders the operation log.
func (s listSession) HistoryTable() string {
	return output.HistoryTable(s.cfg.History)
}

// HistoryRows returns the log newest first for a structured result.
//
// The note is decoded here: the escaping is a storage detail.
func (s listSession) HistoryRows() []listcontent.HistoryRow {
	rows := make([]listcontent.HistoryRow, 0, len(s.cfg.History.AllData))

	for i := len(s.cfg.History.AllData) - 1; i >= 0; i-- {
		entry := s.cfg.History.AllData[i]
		rows = append(rows, listcontent.HistoryRow{
			Version:   entry.Version,
			File:      entry.File,
			Timestamp: entry.Timestemp,
			Operation: string(entry.Operation),
			Command:   entry.CommandName,
			Note:      notes.Display(entry.Note),
		})
	}

	return rows
}

// List reports the running processes, optionally filtered.
func (s listSession) List(partName string) ([]listcontent.ProcessRow, error) {
	return s.listings.Processes(partName)
}

// FolderText runs the platform's folder listing.
func (s listSession) FolderText(path string) (string, error) {
	return s.listings.FolderText(path)
}

// FolderEntries reads the directory behind a structured listing.
func (s listSession) FolderEntries(path string) ([]listcontent.FolderEntry, error) {
	return s.listings.FolderEntries(path)
}

// Variables returns the raw environment listing.
func (s listSession) Variables() (string, error) {
	return s.listings.Variables()
}

// Path returns the raw PATH value.
func (s listSession) Path() (string, error) {
	return s.listings.Path()
}

// ListCommandsNote is the history note for `list`.
func (s listSession) ListCommandsNote() string {
	return notes.ListCommandsNote()
}

// ListProcessNote is the history note for `list --process`.
func (s listSession) ListProcessNote(args []string) string {
	return notes.ListProcessNote(args)
}

// ListHistoryNote is the history note for `list --history`.
func (s listSession) ListHistoryNote() string {
	return notes.ListHistoryNote()
}

// ListFolderNote is the history note for `list --folder`.
func (s listSession) ListFolderNote(args []string) string {
	return notes.ListFolderNote(args)
}

// ListSystemNote is the history note for `list --system`.
func (s listSession) ListSystemNote() string {
	return notes.ListSystemNote()
}

// ListPathNote is the history note for `list --syspath`.
func (s listSession) ListPathNote() string {
	return notes.ListPathNote()
}

func init() {
	listCmd.Flags().BoolVarP(&processFlag, "process", "p", false, "list all process")
	listCmd.Flags().BoolVarP(&historyFlag, "history", "i", false, "list history")
	listCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "list folder files")
	listCmd.Flags().BoolVar(&systemFlag, "system", false, "list all system variables")
	listCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "list path system variable")
	listCmd.Flags().Bool(flagJSON, false, "print the result as one JSON object")
}
