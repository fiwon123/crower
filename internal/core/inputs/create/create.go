package createinputs

import (
	inputscore "github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
)

// Verify parameters to process create operation
func CheckCreateInput(name *string, alias *[]string, exec *string, app *app.Data) {
	if *name == "" {
		*name = inputscore.GetUserInput("New Name ", inputscore.IsValidInput, app).(string)
	}

	if len(*alias) == 0 {
		ok := inputscore.GetUserConfirmation("Do you want to add alias", app)

		if ok {
			*alias = inputscore.InputAlias(app)
		}
	}

	if *exec == "" {
		*exec = inputscore.GetUserInput("New Exec ", inputscore.IsValidInput, app).(string)
	}

}
