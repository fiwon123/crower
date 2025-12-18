package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data"
)

// Open cfg file based on user operational system(OS).
func Open(cfgFilePath string, app *data.App) {
	switch runtime.GOOS {
	case "windows":
		Execute(*data.NewCommand(
			"open",
			[]string{},
			fmt.Sprintf("start %s", cfgFilePath)),
			app)
	case "linux":
		Execute(*data.NewCommand(
			"open",
			[]string{},
			fmt.Sprintf("xdg-open %s", cfgFilePath)),
			app)
	}

}
