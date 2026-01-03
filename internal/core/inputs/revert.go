package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
)

// Verify parameters to process revert operation
func CheckRevertInput(app *app.Data) bool {
	fmt.Println()
	fmt.Println("New History")
	app.History.ListGoBack(1)

	fmt.Println()
	fmt.Println("Deleted History")
	app.History.ListLast(1)

	fmt.Println()
	fmt.Println("History will revert 1 commit.")
	ok := getUserConfirmation("Continue to revert")
	return ok
}
