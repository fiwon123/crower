package handlers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
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
			commandString = fmt.Sprintf(`'xdg-open %s'`, fstring.String())
		}

		if commandString == "" {
			continue
		}

		fmt.Printf("performing execute...: %s \n", commandString)

		out, err := PerformExecute(commandString)
		if err != nil {
			return fmt.Errorf("error %v out %v", err, string(out))
		}

		fmt.Println(string(out))
		return nil
	}

	return nil
}

// Try to open system UI based on operational system (OS)
func OpenSystem(app *app.Data) error {
	switch runtime.GOOS {
	case "windows":
		return PerformExecuteStart("sysdm.cpl")
	case "linux":
		PerformInteractiveTerminal("nano", "~/.bashrc")
	}

	return nil
}
