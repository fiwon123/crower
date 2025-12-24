package inputs

import "github.com/fiwon123/crower/internal/data/app"

func CheckAddInput(name *string, alias *[]string, exec *string, app *app.Data) {
	if *name == "" {
		*name = getUserInput("New Name ", isValidInput, app).(string)
	}

	if len(*alias) == 0 {
		ok := getUserConfirmation("Do you want to add alias")

		if ok {
			*alias = inputAlias(app)
		}
	}

	if *exec == "" {
		*exec = getUserInput("New Exec ", isValidInput, app).(string)
	}

}
