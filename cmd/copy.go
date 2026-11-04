package cmd

import (
	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/transfercontent"
	"github.com/fiwon123/crower/internal/domain/operation"

	"github.com/spf13/cobra"
)

// Cmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy <source...> <destination>",
	Short: "Copy files or folders to a destination.",
	Long: `Copy files or folders to a destination.

	crower copy notes.txt ~/backup
	crower copy a.txt b.txt ~/backup
	crower copy ./reports ~/backup`,
	Example: `	crower copy notes.txt ~/backup
	crower copy a.txt b.txt ~/backup
	crower copy ./reports ~/backup`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return missingArgs(cmd, "source paths and a destination")
		}

		session := newTransferSession(cfg)

		return transfercontent.New(session, session, session, session).Run(
			transfercontent.Request{Op: operation.Copy, Args: args},
		)
	},
}

// transferSession adapts the configuration to the ports the copy and move use
// case declares. One session serves both commands.
type transferSession struct {
	legacySession

	*system.Transfer

	paths *system.Paths
}

// newTransferSession wires the configuration, the platform transfer and the
// filesystem query.
func newTransferSession(cfg *app.Config) transferSession {
	return transferSession{
		legacySession: newLegacySession(cfg),
		Transfer:      system.NewTransfer(),
		paths:         system.NewPaths(),
	}
}

// IsFile reports whether a source is a file rather than a folder.
func (s transferSession) IsFile(path string) bool {
	return s.paths.IsFile(path)
}

// CopyNote is the history note for a copy.
func (s transferSession) CopyNote(subOp operation.SubOperation, args []string) string {
	return notes.CopyNote(subOp, args)
}

// MoveNote is the history note for a move.
func (s transferSession) MoveNote(subOp operation.SubOperation, args []string) string {
	return notes.MoveNote(subOp, args)
}
