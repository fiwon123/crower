package history

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/internal/data/operation"
)

func GetLast(state operation.State, app *app.Data) *history.Content {
	return app.History.GetLastOperation(state)
}
