package handlers

import (
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/pkg/utils"
)

// List all Process running on user operational system (OS).
func Process(app *data.App) {
	err := utils.ListAllProcess("", true)
	if err != nil {
		app.LoggerInfo.Error("Error getting processes:", err)
	}
}
