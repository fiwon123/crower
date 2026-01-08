package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
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

// Verify parameters to process delete history content operation
func CheckDeleteHistoryContentInput(app *app.Data) (history.Content, bool) {

	var content history.Content
	app.History.List()
	content = getUserInput("Select Row", isValidContentKey, app).(history.Content)

	fmt.Println("-----------------------------------------")
	fmt.Println("Content")
	fmt.Println("Version:    ", content.Version)
	fmt.Println("File:    ", content.File)
	fmt.Println("Note:    ", content.Note)
	fmt.Println()

	ok := getUserConfirmation("Continue to restore")
	return content, ok
}
