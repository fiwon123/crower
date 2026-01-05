package operations

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/handlers"
)

func ListCommands(app *app.Data) {
	handlers.ListCommands(app)
}

func ListProcess(args []string, app *app.Data) {
	handlers.ListProcess(args, app)
}

func ListHistory(app *app.Data) {
	handlers.ListHistory(app)
}

func ListFolder(args []string, app *app.Data) {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	out, err := handlers.ListFolder(currentPath, app)
	assertListResult(out, err)
}

func ListSystem(app *app.Data) {
	out, err := handlers.ListSystem(app)
	fmt.Println()
	if err == nil {
		allSysVariables := strings.Split(out, "\n")
		out = ""
		for _, sysVar := range allSysVariables {
			before, after, ok := strings.Cut(sysVar, "=")
			if !ok {
				continue
			}
			name := before
			paths := after
			out += formatVariable(name, paths)
			out += "\n"
		}
	}
	assertListResult(out, err)
}

func ListSysPath(app *app.Data) {
	out, err := handlers.ListSysPath(app)
	fmt.Println()
	if err == nil {
		out = formatVariable("PATH", out)
	}

	assertListResult(out, err)
}

func formatVariable(name string, paths string) string {
	outBuilder := strings.Builder{}
	splitted := strings.Split(paths, ";")
	outBuilder.WriteString(name)
	outBuilder.WriteString("\n")
	for i, path := range splitted {
		outBuilder.WriteString(strconv.Itoa(i))
		outBuilder.WriteString("- ")
		outBuilder.WriteString(path)
		outBuilder.WriteString("\n")
	}

	return outBuilder.String()
}

func assertListResult(out string, err error) {
	if err != nil {
		fmt.Println("failed to list: ", err, out)
		return
	}

	fmt.Print(out)
}
