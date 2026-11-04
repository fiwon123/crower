package cmd

import (
	"fmt"
	"github.com/fiwon123/crower/internal/adapters/output"
	"github.com/fiwon123/crower/internal/domain/timeline"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/application/deletecommand"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/session"
	"github.com/fiwon123/crower/pkg/utils"

	"github.com/spf13/cobra"
)

var historyFlagDelete bool

// Cmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [command]",
	Short: "Delete a command, a file or a folder.",
	Long: `Delete a command, a file or a folder.

A command is found by name or by alias:
	crower delete greet
	crower delete hi

The last operation can be undone instead:
	crower delete --create
	crower delete --update
	crower delete --execute
	crower delete --history

A file or a folder is deleted at the path you give:
	crower delete --file notes.txt
	crower delete --folder reports

With no arguments crower asks which command to delete.`,
	Example: `	crower delete
	crower delete greet
	crower delete --file notes.txt
	crower delete --folder reports
	crower delete --create`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		request := deletecommand.Request{Mode: deletecommand.ModeCommand, Args: args}

		// The flags decide which job this run is; the use case owns what each
		// job does. A positional path keeps the old convenience of deleting a
		// file or folder without a flag.
		switch {
		case createFlag:
			request.Mode, request.Last = deletecommand.ModeLast, operation.Create
		case updateFlag:
			request.Mode, request.Last = deletecommand.ModeLast, operation.Update
		case executeFlag:
			request.Mode, request.Last = deletecommand.ModeLast, operation.Execute
		case fileFlag:
			request.Mode = deletecommand.ModeFile
		case folderFlag:
			request.Mode = deletecommand.ModeFolder
		case systemFlag:
			request.Mode = deletecommand.ModeSystemVariable
		case sysPathFlag:
			request.Mode = deletecommand.ModeSystemPath
		case historyFlagDelete:
			request.Mode = deletecommand.ModeHistory
		case len(args) == 0:
			request.Mode = deletecommand.ModeCommand
		case utils.IsValidFilePath(args[0]):
			request.Mode = deletecommand.ModeFile
		case utils.IsValidFolderPath(args[0]):
			request.Mode = deletecommand.ModeFolder
		default:
			request.Key = args[0]
		}

		session := newDeleteSession(cfg)

		return deletecommand.New(session, session, session, session, session, session, session).Run(request)
	},
}

// deleteSession adapts the configuration to the ports the delete use case
// declares.
//
// It embeds sessionPorts for the ports delete shares with every other use case
// -- the catalog, the diagnostics stream, the two prompts -- and adds the four
// this operation needs: the platform remover, the environment editor, the
// history log, and the delete notes.
type deleteSession struct {
	sessionPorts

	remover     *system.Remover
	environment *system.Environment
}

// newDeleteSession wires the configuration, the renderer and the two platform
// adapters delete needs.
func newDeleteSession(cfg *session.Session) deleteSession {
	return deleteSession{
		sessionPorts: newSessionPorts(cfg),
		remover:      system.NewRemover(),
		environment:  system.NewEnvironment(),
	}
}

// File deletes one file with the host's own tool.
func (s deleteSession) File(path string) error {
	return s.remover.File(path)
}

// Folder deletes one folder with the host's own tool.
func (s deleteSession) Folder(path string) error {
	return s.remover.Folder(path)
}

// Variable drops one variable from the file the shell reads at start.
func (s deleteSession) Variable(name string) (string, error) {
	return s.environment.Variable(name)
}

// VariablePath drops one entry from PATH.
func (s deleteSession) VariablePath(path string) (string, error) {
	return s.environment.VariablePath(path)
}

// HistoryTable renders the operation log.
func (s deleteSession) HistoryTable() string {
	return output.HistoryTable(s.cfg.History)
}

// ChooseHistoryEntry asks which history row to drop.
//
// The note is decoded here: the escaping is a storage detail, and the use case
// shows what the user is about to lose.
func (s deleteSession) ChooseHistoryEntry() (deletecommand.Entry, error) {
	content, err := s.prompt().ChooseHistoryEntry()
	if err != nil {
		return deletecommand.Entry{}, err
	}

	return deletecommand.Entry{
		Version: content.Version,
		File:    content.File,
		Note:    notes.Display(content.Note),
	}, nil
}

// RemoveHistoryEntry drops one entry and the snapshot it points at, and
// rewrites only the index: the command set did not change.
func (s deleteSession) RemoveHistoryEntry(entry deletecommand.Entry) (string, error) {
	content, found := s.historyEntry(entry.Version)
	if !found {
		return "", fmt.Errorf("history entry %s not found", entry.File)
	}

	if !s.cfg.History.RemoveContent(content) {
		return "", fmt.Errorf("history entry %s not found", entry.File)
	}

	if err := s.cfg.DeleteHistoryBackup(content.File); err != nil {
		return "", fmt.Errorf("delete history backup %s: %w", content.File, err)
	}

	if err := s.cfg.SaveOnlyHistory(); err != nil {
		return "", err
	}

	return "history entry deleted", nil
}

// historyEntry finds the stored row the selected entry came from. The version
// is the identity RemoveContent matches on; the file is what the snapshot is
// called.
func (s deleteSession) historyEntry(version int) (timeline.Content, bool) {
	for _, content := range s.cfg.History.AllData {
		if content.Version == version {
			return content, true
		}
	}

	return timeline.Content{}, false
}

// DeleteCommand is the history note for a deleted command.
func (s deleteSession) DeleteCommand(command catalog.Command, args []string) string {
	return notes.DeleteCommandNote(command, args)
}

// DeleteLast is the history note for a `--last` delete.
func (s deleteSession) DeleteLast(op operation.Operation, command catalog.Command) string {
	return notes.DeleteLastNote(op, command)
}

// DeleteFile is the history note for a deleted file.
func (s deleteSession) DeleteFile(args []string) string {
	return notes.DeleteFileNote(args)
}

// DeleteFolder is the history note for a deleted folder.
func (s deleteSession) DeleteFolder(args []string) string {
	return notes.DeleteFolderNote(args)
}

// DeleteSystemVariable is the history note for `delete --system`.
func (s deleteSession) DeleteSystemVariable(args []string) string {
	return notes.DeleteSystemVariableNote(args)
}

// DeleteSystemPath is the history note for `delete --syspath`.
func (s deleteSession) DeleteSystemPath(args []string) string {
	return notes.DeleteSystemPathNote(args)
}

func init() {
	deleteCmd.Flags().BoolVarP(&createFlag, "create", "c", false, "delete recent created command")
	deleteCmd.Flags().BoolVarP(&updateFlag, "update", "u", false, "delete recent updated command")
	deleteCmd.Flags().BoolVarP(&executeFlag, "execute", "x", false, "delete recent executed command")
	deleteCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "ensure arguments are file paths")
	deleteCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "ensure arguments are folder paths")
	deleteCmd.Flags().BoolVar(&systemFlag, "system", false, "delete system variable")
	deleteCmd.Flags().BoolVar(&sysPathFlag, "syspath", false, "delete path variable")
	deleteCmd.Flags().BoolVar(&historyFlagDelete, "history", false, "delete a history content and linked backup cfg")
}
