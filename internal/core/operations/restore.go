package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/core/inputs"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func Restore(args []string, app *app.Data) {
	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	content, ok := inputs.CheckRestoreInput(app)
	if !ok {
		fmt.Println("Cancelling Restore...")
	}

	out, err := handlers.RestoreHistory(key, content, app)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	fmt.Println("restored command: ", out)

	app.History.Add(state.Restore, "", notes.GenerateRestoreNote(out))
	history.Save(app)
}
