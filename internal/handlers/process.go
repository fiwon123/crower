package handlers

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/utils"
)

// List all Process running on user operational system (OS).
func Process(args []string, app *app.Data) {

	partName := ""
	if len(args) > 0 {
		partName = args[0]
	}

	err := utils.ListAllProcess(partName, true)
	if err != nil {
		app.LoggerInfo.Error("Error getting processes:", err)
	}
}
