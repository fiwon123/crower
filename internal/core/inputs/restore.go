package inputs

import (
	"fmt"
	"strconv"

	"github.com/fiwon123/crower/internal/crerrors"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/history"
)

// Verify parameters to process restore operation
func CheckRestoreInput(app *app.Data) (history.Content, bool) {

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

func isValidContentKey(input string, app *app.Data) (any, error) {
	index, err := strconv.Atoi(input)
	if err != nil {
		return "", crerrors.InvalidRows()
	}

	contents := app.History.AllData
	correctIndex := len(contents) - 1 - index
	if correctIndex < 0 || correctIndex >= len(contents) {
		return "", crerrors.InvalidRows()
	}

	return contents[correctIndex], nil
}
