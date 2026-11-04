package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/application/opencontent"
	"github.com/fiwon123/crower/internal/session"

	"github.com/spf13/cobra"
)

// Cmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open [path...]",
	Short: "Open the configuration, a file, a folder or a system dialog.",
	Long: `Open the configuration, a file, a folder or a system dialog.

	crower open
	crower open --file notes.txt
	crower open --folder ~/projects
	crower open --system`,
	Example: `	crower open
	crower open --file notes.txt
	crower open --folder ~/projects`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newSession(cmd)
		if err != nil {
			return err
		}

		request := opencontent.Request{Args: args}

		switch {
		case folderFlag:
			request.Mode = opencontent.ModeFolder
		case fileFlag:
			request.Mode = opencontent.ModeFile
		case systemFlag:
			request.Mode = opencontent.ModeSystem
		case len(args) > 0:
			request.Mode = opencontent.ModePaths
		default:
			return missingArgs(cmd, "a path, or one of --file, --folder, --system")
		}

		session := newOpenSession(cfg)

		return opencontent.New(session, session, session, session, session).Run(request)
	},
}

// openSession adapts the configuration to the ports the open use case declares.
type openSession struct {
	sessionPorts

	opener *system.Opener
	editor *system.SystemEditor
}

// newOpenSession wires the configuration, the platform opener and the system
// editor.
func newOpenSession(cfg *session.Session) openSession {
	return openSession{
		sessionPorts: newSessionPorts(cfg),
		opener:       system.NewOpener(),
		editor:       system.NewSystemEditor(),
	}
}

// Open opens one path with the platform's default application.
func (s openSession) Open(path string) (string, string, error) {
	return s.opener.Open(path)
}

// OpenEditor opens the platform's environment editor.
func (s openSession) OpenEditor() error {
	return s.editor.Open()
}

// OpenNote is the history note for the positional open.
func (s openSession) OpenNote(args []string) string {
	return notes.OpenNote(args)
}

// OpenFileNote is the history note for `open --file`.
func (s openSession) OpenFileNote(args []string) string {
	return notes.OpenFileNote(args)
}

// OpenFolderNote is the history note for `open --folder`.
func (s openSession) OpenFolderNote(args []string) string {
	return notes.OpenFolderNote(args)
}

// OpenSystemNote is the history note for `open --system`.
func (s openSession) OpenSystemNote() string {
	return notes.OpenSystemNote()
}

func init() {
	openCmd.Flags().BoolVarP(&fileFlag, "file", "f", false, "open cfg file or other file")
	openCmd.Flags().BoolVarP(&folderFlag, "folder", "o", false, "open cfg folder or other folder")
	openCmd.Flags().BoolVarP(&systemFlag, "system", "s", false, "open system variable")
}
