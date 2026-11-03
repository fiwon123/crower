package app

import (
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Get last operation
func (app *Config) GetLast(state operation.Operation) *history.Content {
	return app.History.GetLastOperation(state)
}
