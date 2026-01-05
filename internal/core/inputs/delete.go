package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/handlers"
)

// Verify parameters to process delete operation
func CheckDeleteInput(key *string, app *app.Data) bool {

	if *key == "" {
		handlers.ListCommands(app)
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*key = input
	}

	var command *command.Data
	if *key != "" {
		command = app.AllCommandsByName.Get(*key)
		if command == nil {
			command = app.AllCommandsByAlias.Get(*key)
		}
	}

	if command == nil {
		handlers.ListCommands(app)
		fmt.Println("Command not found, try to select one.")
		input := getUserInput("Select Row", isValidInputKey, app).(string)
		*key = input

		command = app.AllCommandsByName.Get(*key)
	}

	fmt.Println("-----------------------------------------")
	fmt.Println("Name:    ", command.Name)
	fmt.Println("Aliases: ", command.AllAlias)
	fmt.Println("Exec:    ", command.Exec)
	fmt.Println()

	ok := getUserConfirmation("Continue to delete")
	return ok

}
