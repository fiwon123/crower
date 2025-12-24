package inputs

import "github.com/fiwon123/crower/internal/data/app"

func CheckAddInput(name *string, alias *[]string, exec *string, app *app.Data) error {
	if *name == "" {
		*name = getUserInput("New Name:", nil, isEmpty, app).(string)
	}

	if len(*alias) == 0 {
		ans := getUserConfirmation("Do you want to add alias", nil, isValidConfirmation, app).(bool)

		if ans {
			*alias = selectAliases(app)
		}
	}

	if *exec == "" {
		*exec = getUserInput("New Exec:", nil, isEmpty, app).(string)
	}

	return nil
}

func checkEmptyAnswer(input string) bool {
	if input == input_y || input == input_yes {
		return true
	}

	return false
}
