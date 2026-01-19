package openhandlers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	executehandlers "github.com/fiwon123/crower/internal/handlers/execute"
)

// Open filepath based on user operational system(OS).
func Open(paths []string, app *app.Data) error {

	for _, f := range paths {
		commandString := ""

		var fstring strings.Builder
		fstring.WriteString(f)

		switch runtime.GOOS {
		case "windows":
			commandString = fmt.Sprintf(`start ' ' '%s'`, fstring.String())
		case "linux":
			commandString = fmt.Sprintf(`"xdg-open '%s'"`, fstring.String())
		}

		if commandString == "" {
			continue
		}

		app.Logger.Info("performing execute...: ", "exec", commandString)

		out, err := executehandlers.PerformExecute(commandString, app)
		if err != nil {
			return fmt.Errorf("error %v out %v", err, string(out))
		}

		app.Logger.Info(string(out))
		return nil
	}

	return nil
}

// Try to open system UI based on operational system (OS)
func OpenSystem(app *app.Data) error {
	switch runtime.GOOS {
	case "windows":
		return executehandlers.PerformExecuteStart("sysdm.cpl", app)
	case "linux":
		executehandlers.PerformInteractiveTerminal("nano", "~/.bashrc")
	}

	return nil
}
