package session

import (
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// RevertTo restores the command map an entry points at, then drops that entry
// and every entry after it.
func (s *Session) RevertTo(entry *timeline.Content) error {
	commands, err := s.historyFiles.ReadSnapshot(entry.File)
	if err != nil {
		return err
	}

	if err := s.commands.Save(commands); err != nil {
		return err
	}

	s.removeUntilHistory(entry)

	return nil
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
