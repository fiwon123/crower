package handlers

import (
	"fmt"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
)

// Search based on user operational system(OS).
func Search(content string, app *app.Data) {
	switch runtime.GOOS {
	case "windows":
		PerformExecute(fmt.Sprintf(`'start "https://duckduckgo.com/?q=%s" '`, content))
	case "linux":
		PerformExecute(fmt.Sprintf(`'xdg-open "https://duckduckgo.com/?q=%s" '`, content))
	}

}
