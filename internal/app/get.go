package app

import (
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Get last operation
func (app *Config) GetLast(state operation.Operation) *timeline.Content {
	return app.History.GetLastOperation(state)
}
