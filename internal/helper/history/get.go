package historyhelper

import (
	appdata "github.com/fiwon123/crower/internal/data/app"
	historydata "github.com/fiwon123/crower/internal/data/history"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
)

// Get last operation
func GetLast(state operationsdata.MainOperationEnum, app *appdata.Data) *historydata.Content {
	return app.History.GetLastOperation(state)
}
