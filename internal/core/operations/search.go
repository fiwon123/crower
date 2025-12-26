package operations

import (
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Search(args []string, app *app.Data) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	handlers.Search(content, app)
}
