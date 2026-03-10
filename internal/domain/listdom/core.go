package listdom

import (
	"os"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
)

type Core struct {
	app     *app.Config
	handler *Handler
}

func NewCore(app *app.Config) *Core {

	handler := NewHandler(app)

	return &Core{
		handler: handler,
		app:     app,
	}
}

func (c *Core) ListCommands() {
	c.handler.ListCommands()

	c.app.History.Add(operations.List, newListCommandsNote())
	c.app.Save()
}

func (c *Core) ListProcess(args []string) {
	c.handler.ListProcess(args)

	c.app.History.Add(operations.List, newListProcessNote(args))
	c.app.Save()
}

func (c *Core) ListHistory() {
	c.handler.ListHistory()

	c.app.History.Add(operations.List, newListHistoriesNote())
	c.app.Save()
}

func (c *Core) ListFolder(args []string) {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	out, err := c.handler.ListFolder(currentPath)
	c.assertListResult(out, err)

	c.app.History.Add(operations.List, newListFolderNote(args))
	c.app.Save()
}

func (c *Core) ListSystem() {
	out, err := c.handler.ListSystem()
	c.app.Logger.Info("")
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
	c.assertListResult(out, err)

	c.app.History.Add(operations.List, newListSystemNote())
	c.app.Save()
}

func (c *Core) ListSysPath() {
	out, err := c.handler.ListSysPath()
	c.app.Logger.Info("")
	if err == nil {
		out = formatVariable("PATH", out)
	}

	c.assertListResult(out, err)

	c.app.History.Add(operations.List, newListSystemPathNote())
	c.app.Save()
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

func (c *Core) assertListResult(out string, err error) {
	if err != nil {
		c.app.Logger.Error("failed to list: ", "error", err, "out", out)
		return
	}

	c.app.Logger.Info(out)
}
