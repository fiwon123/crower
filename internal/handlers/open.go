package handlers

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
)

// Open filepath based on user operational system(OS).
func Open(paths []string, app *app.Data) {

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
			fmt.Println("command")
			return
		}

		fmt.Printf("performing execute...: %s \n", commandString)

		out, err := PerformExecute(commandString)
		if err != nil {
			fmt.Println(err)
			fmt.Println(string(out))
			return
		}

		fmt.Println(string(out))
	}

}

// Try to open system UI based on operational system (OS)
func OpenSystem(app *app.Data) ([]byte, error) {
	switch runtime.GOOS {
	case "windows":
		return PerformExecute("'sysdm.cpl'")
	case "linux":
		PerformInteractiveTerminal("nano", "~/.bashrc")
	}

	return nil, nil
}
