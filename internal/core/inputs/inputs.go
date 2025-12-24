package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

const (
	input_n   = "n"
	input_no  = "no"
	input_y   = "y"
	input_yes = "yes"
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

func getUserInput(ask string, fnBefore func(), fnValid func(string, *app.Data) (any, error), app *app.Data) any {
	ok := false
	input := ""
	var output any
	var err error
	for !ok {
		if fnBefore != nil {
			fnBefore()
		}
		fmt.Print(ask)
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if output, err = fnValid(input, app); err != nil {
			fmt.Println(err)
			continue
		}

		ok = true
	}

	return output
}

func getUserConfirmation(ask string, fnBefore func(), fnValid func(string, *app.Data) (any, error), app *app.Data) any {
	ok := false
	input := ""
	var output any
	var err error
	for !ok {
		if fnBefore != nil {
			fnBefore()
		}

		fmt.Print(ask + " ([Y]es/[N]o): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if output, err = fnValid(input, app); err != nil {
			fmt.Println(err)
			continue
		}

		ok = true
	}

	return output
}

func selectInput(input string, app *app.Data) (any, error) {
	if !checkValidAnswer(input) {
		return "", fmt.Errorf("Invalid Input")
	}

	if checkNoAnswer(input) {
		return "", nil
	}

	return input, nil
}

func isValidConfirmation(input string, app *app.Data) (any, error) {
	if !checkValidAnswer(input) {
		return false, fmt.Errorf("Invalid Input")
	}

	if checkNoAnswer(input) {
		return false, nil
	}

	return true, nil
}

func isEmpty(input string, app *app.Data) (any, error) {
	if input == "" {
		return "", fmt.Errorf("input is empty")
	}

	return input, nil
}

func selectAliases(app *app.Data) []string {
	output := []string{}

	fmt.Println()
	alias := "none"
	for alias != "" {
		fmt.Println("current aliases: ", output)
		alias = getUserInput("Add new alias (type enter to skip): ", nil, selectNewAlias, app).(string)

		if alias != "" {
			output = append(output, alias)
		}
	}

	return output
}
