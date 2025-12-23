package handlers

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/pkg/utils"
)

func History(app *data.App) error {

	var history data.History
	err := utils.ReadJson(app.HistoryFilePath, &history)
	if err != nil {
		app.LoggerInfo.Error("history error: ", err)
	}
	
	history.List()

	return nil
}
