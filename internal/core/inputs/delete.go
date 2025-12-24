package inputs

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckDeleteInput(name *string, allAlias *[]string, app *app.Data) {

	if *name == "" && len(*allAlias) == 0 {
		handlers.List(app)
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*name = input
	}

}
