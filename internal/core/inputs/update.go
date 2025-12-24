package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string, app *app.Data) bool {

	if *key == "" {
		handlers.List(app)
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*key = input
	}

	fmt.Println("-----------------------------------------")
	updateCommand := app.AllCommandsByName.Get(*key)
	fmt.Println("Name:    ", updateCommand.Name)
	fmt.Println("Aliases: ", updateCommand.AllAlias)
	fmt.Println("Exec:    ", updateCommand.Exec)
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
	fmt.Println("-----------------------------------------")
	fmt.Println("Old Command: ")
	fmt.Println("Name:    ", updateCommand.Name)
	fmt.Println("Aliases: ", updateCommand.AllAlias)
	fmt.Println("Exec:    ", updateCommand.Exec)
	fmt.Println("-----------------------------------------")
	fmt.Println("New Command: ")
	fmt.Println("Name:    ", *name)
	fmt.Println("Aliases: ", *allAlias)
	fmt.Println("Exec:    ", *exec)
	fmt.Println()

	ok := getUserConfirmation("Continue to update")
	return ok
}
