package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

// Open filepath based on user operational system(OS).
func Open(paths []string, app *app.Data) {

	for _, f := range paths {
		switch runtime.GOOS {
		case "windows":
			PerformExecute(fmt.Sprintf("'start %s'", f))
		case "linux":
			PerformExecute(fmt.Sprintf("'xdg-open %s'", f))
		}

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
