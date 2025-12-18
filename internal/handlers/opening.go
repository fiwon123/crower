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
		PerformExecute(fmt.Sprintf("start %s", cfgFilePath))
	case "linux":
		PerformExecute(fmt.Sprintf("xdg-open %s", cfgFilePath))
	}

}
