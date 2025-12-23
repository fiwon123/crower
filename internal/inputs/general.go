package inputs

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fiwon123/crower/internal/data"
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

func getUserInput(ask string, fnBefore func(), fnValid func(string, *data.App) (any, error), app *data.App) any {
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
