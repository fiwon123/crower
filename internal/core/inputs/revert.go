package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
)

// Verify parameters to process revert operation
func CheckRevertInput(steps int, app *app.Data) (bool, error) {
	if steps < 0 {
		return false, fmt.Errorf("steps is negative number")
	}

	stopIndex := app.History.GetIndexFromLastTo(steps)
	if stopIndex < 0 {
		return false, fmt.Errorf("can't revert steps is greater than quantity of history registry")
	}

	fmt.Println()
	fmt.Println("Deleted History")
	app.History.ListLastHistory(steps)

	fmt.Println()
	fmt.Println("New History")
	app.History.ListFirstHistory(stopIndex)

	fmt.Println()
	fmt.Printf("History will revert %d registries \n", steps)
	ok := getUserConfirmation("Continue to revert")

	if !ok {
		return false, fmt.Errorf("cancelling revert...")
	}

	return ok, nil
}
