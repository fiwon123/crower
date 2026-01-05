package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
)

const (
	input_n   string = "n"
	input_no  string = "no"
	input_y   string = "y"
	input_yes string = "yes"
)

func checkValidAnswer(input string) bool {
	if input == input_y ||
		input == input_n ||
		input == input_yes ||
		input == input_no {
		return true
	}

	return false
}

func checkNoAnswer(input string) bool {
	if input == input_n || input == input_no {
		return true
	}

	return false
}

func checkYesAnswer(input string) bool {
	if input == input_y || input == input_yes {
		return true
	}

	return false
}

func getUserInput(ask string, fnValid func(string, *app.Data) (any, error), app *app.Data) any {
	ok := false
	input := ""
	var output any
	var err error
	for !ok {
		fmt.Print(ask + ": ")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if output, err = fnValid(input, app); err != nil {
			fmt.Println(err)
			continue
		}

		ok = true
	}

	return output
}

func getUserConfirmation(ask string) bool {
	ok := false
	input := ""
	var confirmation bool
	var err error
	for !ok {

		fmt.Print(ask + " ([Y]es/[N]o): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if confirmation, err = isValidConfirmation(input); err != nil {
			fmt.Println(err)
			continue
		}

		ok = true
	}

	return confirmation
}

func isValidConfirmation(input string) (bool, error) {
	if !checkValidAnswer(input) {
		return false, crerrors.InvalidInput()
	}

	if checkNoAnswer(input) {
		return false, nil
	}

	return true, nil
}

func isValidInput(input string, app *app.Data) (any, error) {
	if input == "" {
		return "", crerrors.EmptyInput()
	}

	return input, nil
}

func inputAlias(app *app.Data) []string {
	output := []string{}

	alias := "none"
	for alias != "" {
		fmt.Println("current aliases: ", output)
		alias = getUserInput("Add new alias (type enter to skip): ", isValidAlias, app).(string)

		if alias != "" {
			output = append(output, alias)
		}
	}

	return output
}

func inputName(app *app.Data) string {
	name := ""
	for name == "" {
		name = getUserInput("Add new name: ", isValidInput, app).(string)
	}

	return name
}

func inputExec(app *app.Data) string {
	exec := ""
	for exec == "" {
		exec = getUserInput("Add new exec: ", isValidInput, app).(string)
	}

	return exec
}

func isValidAlias(input string, app *app.Data) (any, error) {
	for _, r := range input {
		if !unicode.IsNumber(r) && !unicode.IsLetter(r) {
			return "", crerrors.OnlyLettersAndNumbers()
		}
	}

	return input, nil
}

func isValidInputKey(input string, app *app.Data) (any, error) {
	index, err := strconv.Atoi(input)
	if err != nil {
		return "", crerrors.InvalidRows()
	}

	if index < 0 || index >= len(app.OrderKeys) {
		return "", crerrors.InvalidRows()
	}

	return app.OrderKeys[index], nil
}
