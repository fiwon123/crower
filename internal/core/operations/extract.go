package operations

import (
	"fmt"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func Extract(args []string, outDir string, app *app.Data) {

	if outDir == "" {
		outDir = "./"
	}

	if len(args) == 0 {
		fmt.Println("empty paths")
		return
	}

	handlers.Extract(args, outDir, app)
}
