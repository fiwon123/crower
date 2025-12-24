package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckDeleteInput(name *string, allAlias *[]string, app *app.Data) bool {

	if *name == "" && len(*allAlias) == 0 {
		handlers.List(app)
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*name = input
	}

	var command *command.Data
	if *name != "" {
		command = app.AllCommandsByName.Get(*name)
	} else if len(*allAlias) > 0 {
		command = app.AllCommandsByAlias.Get((*allAlias)[0])
	}

	if command == nil {
		handlers.List(app)
		fmt.Println("Command not found, try to select one.")
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*name = input

		command = app.AllCommandsByName.Get(*name)
	}

	fmt.Println("-----------------------------------------")
	fmt.Println("Name:    ", command.Name)
	fmt.Println("Aliases: ", command.AllAlias)
	fmt.Println("Exec:    ", command.Exec)
	fmt.Println()

	ok := getUserConfirmation("Continue to delete")
	return ok

}
