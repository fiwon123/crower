package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string, app *app.Data) error {

	if *key == "" {
		handlers.List(app)
		input := getUserInput("Select Row: ", isValidInputKey, app).(string)
		*key = input
	}

	fmt.Println("-----------------------------------------")
	command := app.AllCommandsByName.Get(*key)
	fmt.Println("Name:    ", command.Name)
	fmt.Println("Aliases: ", command.AllAlias)
	fmt.Println("Exec:    ", command.Exec)
	fmt.Println()

	if *name == "" {
		ok := getUserConfirmation("Do you want to update name")

		if ok {
			*name = inputName(app)
		}
	}

	if len(*allAlias) == 0 {
		ok := getUserConfirmation("Do you want to update alias")

		if ok {
			*allAlias = inputAlias(app)
		}
	}

	if *exec == "" {
		ok := getUserConfirmation("Do you want to update exec")

		if ok {
			*exec = inputExec(app)
		}
	}

	return nil
}
