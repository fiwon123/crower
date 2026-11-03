package helper

import (
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

// GetUserInput asks a question until fnValid accepts the answer.
//
// It returns an error instead of looping forever. A closed input stream used to
// be an empty answer, so `crower create </dev/null` re-prompted as fast as it
// could write: 876,468 lines in five seconds and no exit. With --no-input the
// prompt fails before it is printed, so a script gets a clear error rather than
// a wait.
func GetUserInput(ask string, fnValid func(string, *app.Config) (any, error), a *app.Config) (any, error) {
	for {
		if a.NoInput {
			return nil, errors.InputRequired(ask)
		}

		a.Logger.Info(ask + ": ")

		input, err := a.ReadLine()
		if err != nil {
			return nil, err
		}

		output, err := fnValid(input, a)
		if err == nil {
			return output, nil
		}

		a.Logger.Info(err.Error())
	}
}

// GetUserConfirmation asks a yes/no question.
//
// A closed stream is an error, not a "no": a script must not silently skip an
// operation because nobody was there to answer.
func GetUserConfirmation(ask string, a *app.Config) (bool, error) {
	for {
		if a.NoInput {
			return false, errors.InputRequired(ask)
		}

		a.Logger.Info(ask + " ([Y]es/[N]o): ")

		input, err := a.ReadLine()
		if err != nil {
			return false, err
		}

		confirmation, err := IsValidConfirmation(input)
		if err == nil {
			return confirmation, nil
		}

		a.Logger.Info(err.Error())
	}
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

func IsValidInput(input string, a *app.Config) (any, error) {
	if input == "" {
		return "", errors.EmptyInput()
	}

	return input, nil
}

// InputAlias asks for aliases until the user enters an empty line.
func InputAlias(a *app.Config) ([]string, error) {
	output := []string{}

	alias := "none"
	for alias != "" {
		a.Logger.Info("current aliases: ", output)

		value, err := GetUserInput("Add new alias (type enter to skip): ", isValidAlias, a)
		if err != nil {
			return nil, err
		}

		alias = value.(string)
		if alias != "" {
			output = append(output, alias)
		}
	}

	return output, nil
}

// InputName asks for a required name.
func InputName(a *app.Config) (string, error) {
	return requiredInput("Add new name: ", a)
}

// InputExec asks for a required command line.
func InputExec(a *app.Config) (string, error) {
	return requiredInput("Add new exec: ", a)
}

// requiredInput asks until the answer is not empty.
func requiredInput(ask string, a *app.Config) (string, error) {
	for {
		value, err := GetUserInput(ask, IsValidInput, a)
		if err != nil {
			return "", err
		}

		input, ok := value.(string)
		if ok && input != "" {
			return input, nil
		}
	}
}

func isValidAlias(input string, a *app.Config) (any, error) {
	for _, r := range input {
		if !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			return "", errors.OnlyLettersAndNumbers()
		}
	}

	return input, nil
}

// IsValidInputKey maps a printed row number to the command name on that row.
func IsValidInputKey(input string, a *app.Config) (any, error) {
	index, err := rowIndex(input)
	if err != nil {
		return "", err
	}

	if index < 0 || index >= len(a.OrderKeys) {
		return "", errors.InvalidRows()
	}

	return a.OrderKeys[index], nil
}

// IsValidContentKey maps a printed row number to a history entry. The history
// is printed newest first, so row 0 is the last entry.
func IsValidContentKey(input string, a *app.Config) (any, error) {
	index, err := rowIndex(input)
	if err != nil {
		return "", err
	}

	contents := a.History.AllData
	correctIndex := len(contents) - 1 - index
	if correctIndex < 0 || correctIndex >= len(contents) {
		return "", errors.InvalidRows()
	}

	return contents[correctIndex], nil
}

func rowIndex(input string) (int, error) {
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, errors.InvalidRows()
	}

	return index, nil
}
