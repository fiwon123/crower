package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

// Open cfg file based on user operational system(OS).
func Open(cfgFilePath string, app *app.Data) {
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf("'start %s'", cfgFilePath))
	case "linux":
		PerformExecute(fmt.Sprintf("'xdg-open %s'", cfgFilePath))
	}

}

func OpenSystem(app *app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute("'sysdm.cpl'")
	case "linux":
		PerformInteractiveTerminal("nano", "~/.bashrc")
	}

	return nil, nil
}
