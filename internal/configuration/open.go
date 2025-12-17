package configuration

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

// Open cfg file based on user operational system(OS).
func Open(cfgFilePath string, app *data.App) {
	switch runtime.GOOS {
	case "windows":
		handlers.Execute(*data.NewCommand(
			"open",
			[]string{},
			fmt.Sprintf("start %s", cfgFilePath)),
			app)
	case "linux":
		handlers.Execute(*data.NewCommand(
			"open",
			[]string{},
			fmt.Sprintf("xdg-open %s", cfgFilePath)),
			app)
	}

}
