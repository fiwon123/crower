package helper

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/errors"
)

const (
	input_n   string = "n"
	input_no  string = "no"
	input_y   string = "y"
	input_yes string = "yes"
)

func CheckValidAnswer(input string) bool {
	if input == input_y ||
		input == input_n ||
		input == input_yes ||
		input == input_no {
		return true
	}

	return false
}

func CheckNoAnswer(input string) bool {
	if input == input_n || input == input_no {
		return true
	}

	return false
}

func CheckYesAnswer(input string) bool {
	if input == input_y || input == input_yes {
		return true
	}

	return false
}

func GetUserInput(ask string, fnValid func(string, *app.Config) (any, error), app *app.Config) any {
	ok := false
	input := ""
	var output any
	var err error
	for !ok {
		app.Logger.Info(ask + ": ")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if output, err = fnValid(input, app); err != nil {
			app.Logger.Info(err.Error())
			continue
		}

		ok = true
	}

	return output
}

func GetUserConfirmation(ask string, app *app.Config) bool {
	ok := false
	input := ""
	var confirmation bool
	var err error
	for !ok {

		app.Logger.Info(ask + " ([Y]es/[N]o): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if confirmation, err = IsValidConfirmation(input); err != nil {
			app.Logger.Info(err.Error())
			continue
		}

		ok = true
	}

	return confirmation
}

func IsValidConfirmation(input string) (bool, error) {
	if !CheckValidAnswer(input) {
		return false, errors.InvalidInput()
	}

	if CheckNoAnswer(input) {
		return false, nil
	}

	return true, nil
}

func IsValidInput(input string, app *app.Config) (any, error) {
	if input == "" {
		return "", errors.EmptyInput()
	}

	return input, nil
}

func InputAlias(app *app.Config) []string {
	output := []string{}

	alias := "none"
	for alias != "" {
		app.Logger.Info("current aliases: ", output)
		alias = GetUserInput("Add new alias (type enter to skip): ", isValidAlias, app).(string)

		if alias != "" {
			output = append(output, alias)
		}
	}

	return output
}

func InputName(app *app.Config) string {
	name := ""
	for name == "" {
		name = GetUserInput("Add new name: ", IsValidInput, app).(string)
	}

	return name
}

func InputExec(app *app.Config) string {
	exec := ""
	for exec == "" {
		exec = GetUserInput("Add new exec: ", IsValidInput, app).(string)
	}

	return exec
}

func isValidAlias(input string, app *app.Config) (any, error) {
	for _, r := range input {
		if !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			return "", errors.OnlyLettersAndNumbers()
		}
	}

	return input, nil
}

func IsValidInputKey(input string, app *app.Config) (any, error) {
	index, err := strconv.Atoi(input)
	if err != nil {
		return "", errors.InvalidRows()
	}

	if index < 0 || index >= len(app.OrderKeys) {
		return "", errors.InvalidRows()
	}

	return app.OrderKeys[index], nil
}

func IsValidContentKey(input string, app *app.Config) (any, error) {
	index, err := strconv.Atoi(input)
	if err != nil {
		return "", errors.InvalidRows()
	}

	contents := app.History.AllData
	correctIndex := len(contents) - 1 - index
	if correctIndex < 0 || correctIndex >= len(contents) {
		return "", errors.InvalidRows()
	}

	return contents[correctIndex], nil
}
