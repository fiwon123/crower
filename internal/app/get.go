package app

import (
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/data/operations"
)

// Get last operation
func (app *Config) GetLast(state operations.MainOperationEnum) *history.Content {
	return app.History.GetLastOperation(state)
}
