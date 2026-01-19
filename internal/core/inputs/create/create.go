package createinputs

import (
	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
)

// Verify parameters to process create operation
func CheckCreateInput(name *string, alias *[]string, exec *string, app *app.Data) {
	if *name == "" {
		*name = inputs.GetUserInput("New Name ", inputs.IsValidInput, app).(string)
	}

	if len(*alias) == 0 {
		ok := inputs.GetUserConfirmation("Do you want to add alias", app)

		if ok {
			*alias = inputs.InputAlias(app)
		}
	}

	if *exec == "" {
		*exec = inputs.GetUserInput("New Exec ", inputs.IsValidInput, app).(string)
	}

}
