package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/operation"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
	"github.com/fiwon123/crower/pkg/utils"
)

func Reset(app *app.Data) {
	ok := inputs.CheckResetInput(app)

	if !ok {
		fmt.Println("Cancelling reset...")
		return
	}

	app.LoggerInfo.Info("reset all commands: ", app.AllCommandsByName)
	handlers.Reset(app)
	utils.WriteToml(app.AllCommandsByName, app.CfgFilePath)

	app.History.Add(operation.Reset, "", notes.GenerateResetNote())
	history.Save(app)
}
