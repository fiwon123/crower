package session

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Get last operation
func (s *Session) GetLast(state operation.Operation) *timeline.Content {
	return s.History.GetLastOperation(state)
}
