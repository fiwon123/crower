package listoperations

import (
	"os"
	"strconv"
	"strings"

	appdata "github.com/fiwon123/crower/internal/data/app"
	operationsdata "github.com/fiwon123/crower/internal/data/operations"
	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func ListCommands(app *appdata.Data) {
	listhandlers.ListCommands(app)

	app.History.Add(operationsdata.List, notes.GenerateListCommandsNote())
	history.Save(app)
}

func ListProcess(args []string, app *appdata.Data) {
	listhandlers.ListProcess(args, app)

	app.History.Add(operationsdata.List, notes.GenerateListProcessNote(args))
	history.Save(app)
}

func ListHistory(app *appdata.Data) {
	listhandlers.ListHistory(app)

	app.History.Add(operationsdata.List, notes.GenerateListHistoriesNote())
	history.Save(app)
}

func ListFolder(args []string, app *appdata.Data) {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	out, err := listhandlers.ListFolder(currentPath, app)
	assertListResult(out, err, app)

	app.History.Add(operationsdata.List, notes.GenerateListFolderNote(args))
	history.Save(app)
}

func ListSystem(app *appdata.Data) {
	out, err := listhandlers.ListSystem(app)
	app.Logger.Info("")
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
	assertListResult(out, err, app)

	app.History.Add(operationsdata.List, notes.GenerateListSystemNote())
	history.Save(app)
}

func ListSysPath(app *appdata.Data) {
	out, err := listhandlers.ListSysPath(app)
	app.Logger.Info("")
	if err == nil {
		out = formatVariable("PATH", out)
	}

	assertListResult(out, err, app)

	app.History.Add(operationsdata.List, notes.GenerateListSystemPathNote())
	history.Save(app)
}

func formatVariable(name string, paths string) string {
	outBuilder := strings.Builder{}
	splitted := strings.Split(paths, string(os.PathListSeparator))
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

func assertListResult(out string, err error, app *appdata.Data) {
	if err != nil {
		app.Logger.Error("failed to list: ", "error", err, "out", out)
		return
	}

	app.Logger.Info(out)
}
