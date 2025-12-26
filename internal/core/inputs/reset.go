package inputs

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func CheckResetInput(app *app.Data) bool {
	fmt.Println("-----------------------------------------")
	handlers.ListCommands(app)

	fmt.Println()
	fmt.Println("All commands will be erased...")
	ok := getUserConfirmation("Continue to reset")
	return ok
}
