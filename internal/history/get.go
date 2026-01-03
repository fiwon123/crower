package history

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/data/state"
)

// Get last operation
func GetLast(state state.OperationEnum, app *app.Data) *history.Content {
	return app.History.GetLastOperation(state)
}
