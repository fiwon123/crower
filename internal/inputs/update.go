package inputs

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string, app *data.App) error {

	if *key == "" {
		handlers.List(app)
		input := getUserInput("Select Row: ", nil, selectInputKey, app).(string)
		*key = input
	}

	fmt.Println("-----------------------------------------")
	command := app.AllCommandsByName.Get(*key)
	fmt.Println("Name:    ", command.Name)
	fmt.Println("Aliases: ", command.AllAlias)
	fmt.Println("Exec:    ", command.Exec)
	fmt.Println()

	if *name == "" {
		*name = getUserInput("Do you want to update name ([Y]es/[N]o): ", nil, selectInputName, app).(string)
	}

	if len(*allAlias) == 0 {
		*allAlias = getUserInput("Do you want to update alias ([Y]es/[N]o): ", nil, selectInputAliases, app).([]string)
	}

	if *exec == "" {
		*exec = getUserInput("Do you want to update exec ([Y]es/[N]o): ", nil, selectInputExec, app).(string)
	}

	return nil
}

func selectInputName(input string, app *data.App) (any, error) {
	if !checkValidAnswer(input) {
		return "", fmt.Errorf("Invalid Input")
	}

	if checkNoAnswer(input) {
		return "", nil
	}

	fmt.Println()
	name := ""
	for name == "" {
		name = getUserInput("Add new name: ", nil, selectNewExec, app).(string)
	}

	return name, nil
}

func selectInputExec(input string, app *data.App) (any, error) {
	if !checkValidAnswer(input) {
		return "", fmt.Errorf("Invalid Input")
	}

	if checkNoAnswer(input) {
		return "", nil
	}

	fmt.Println()
	exec := ""
	for exec == "" {
		exec = getUserInput("Add new exec: ", nil, selectNewExec, app).(string)
	}

	return exec, nil
}

func selectInputAliases(input string, app *data.App) (any, error) {
	if !checkValidAnswer(input) {
		return "", fmt.Errorf("Invalid Input")
	}

	output := []string{}
	if checkNoAnswer(input) {
		return output, nil
	}

	fmt.Println()
	alias := "none"
	for alias != "" {
		fmt.Println("current aliases: ", output)
		alias = getUserInput("Add new alias (type enter to skip): ", nil, selectNewAlias, app).(string)

		if alias != "" {
			output = append(output, alias)
		}
	}

	return output, nil
}

func selectNewExec(input string, app *data.App) (any, error) {
	if input == "" {
		return "", fmt.Errorf("input is empty")
	}

	return input, nil
}

func selectNewAlias(input string, app *data.App) (any, error) {
	if !hasOnlyNumbersAndLetters(input) {
		return "", fmt.Errorf("only letters allowed")
	}

	return input, nil
}

func hasOnlyNumbersAndLetters(input string) bool {
	for _, r := range input {
		if !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func selectInputKey(input string, app *data.App) (any, error) {
	index, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("Invalid row")
	}

	if index < 0 || index >= len(app.OrderKeys) {
		return "", fmt.Errorf("Invalid row")
	}

	return app.OrderKeys[index], nil
}
