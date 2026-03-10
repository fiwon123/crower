package historyhelper

import (
	"github.com/fiwon123/crower/internal/app"
	historydata "github.com/fiwon123/crower/internal/data/history"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Get last operation
func GetLast(state operationsdata.MainOperationEnum, app *app.Config) *historydata.Content {
	return app.History.GetLastOperation(state)
}
