package listoperations

import (
	"os"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/internal/data/state"
	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
	"github.com/fiwon123/crower/internal/history"
	"github.com/fiwon123/crower/internal/history/notes"
)

func ListCommands(app *app.Data) {
	listhandlers.ListCommands(app)

	app.History.Add(state.List, notes.GenerateListCommandsNote())
	history.Save(app)
}

func ListProcess(args []string, app *app.Data) {
	listhandlers.ListProcess(args, app)

	app.History.Add(state.List, notes.GenerateListProcessNote(args))
	history.Save(app)
}

func ListHistory(app *app.Data) {
	listhandlers.ListHistory(app)

	app.History.Add(state.List, notes.GenerateListHistoriesNote())
	history.Save(app)
}

func ListFolder(args []string, app *app.Data) {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	out, err := listhandlers.ListFolder(currentPath, app)
	assertListResult(out, err, app)

	app.History.Add(state.List, notes.GenerateListFolderNote(args))
	history.Save(app)
}

func ListSystem(app *app.Data) {
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

	app.History.Add(state.List, notes.GenerateListSystemNote())
	history.Save(app)
}

func ListSysPath(app *app.Data) {
	out, err := listhandlers.ListSysPath(app)
	app.Logger.Info("")
	if err == nil {
		out = formatVariable("PATH", out)
	}

	assertListResult(out, err, app)

	app.History.Add(state.List, notes.GenerateListSystemPathNote())
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

func assertListResult(out string, err error, app *app.Data) {
	if err != nil {
		app.Logger.Error("failed to list: ", "error", err, "out", out)
		return
	}

	app.Logger.Info(out)
}
