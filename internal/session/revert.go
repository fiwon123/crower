package session

import (
	"slices"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// RevertTo restores the command map an entry points at, then drops that entry
// and every entry after it.
func (s *Session) RevertTo(entry *timeline.Content) error {
	commands, err := s.historyFiles.ReadSnapshot(entry.File)
	if err != nil {
		return err
	}

	// Revert restores a state that already existed; the load warning already
	// told the user about a literal it holds. Refusing the undo would trap them
	// in the state they are trying to leave, so this path writes it back.
	if err := s.commands.Save(commands); err != nil {
		return err
	}

	// The restored map has to replace the in-memory one as well. The CLI exits
	// after a revert and never noticed, but the Fyne window keeps running: a
	// stale catalog would show the commands the revert just removed, and a
	// later Save would write them back over the restored file.
	s.ReplaceCatalog(catalogOf(commands, sortedNames(commands)))

	s.removeUntilHistory(entry)

	return nil
}

// sortedNames is the order a fresh load sees: the TOML encoder writes map keys
// in sorted order, so the in-memory catalog must use the same order or the
// list reorders on the next start.
func sortedNames(commands command.MapData) []string {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}

	slices.Sort(names)

	return names
}

// removeUntilHistory drops entries from the newest down to the one that was
// restored, and the snapshot file that belongs to each.
func (s *Session) removeUntilHistory(entry *timeline.Content) {
	for {
		last := s.History.GetLast()
		if last == nil || last.Version == entry.Version {
			return
		}

		if s.historyFiles.HasSnapshot(last.File) {
			_ = s.historyFiles.DeleteSnapshot(last.File)
		}

		s.History.RemoveLast()
	}
}
