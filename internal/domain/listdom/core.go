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
		app:     app,
		handler: handler,
	}
}

// ListCommands prints the command table and records the operation.
func (c *Core) ListCommands() error {
	c.handler.ListCommands()

	c.app.History.Add(operations.List, newListCommandsNote())
	return c.app.Save()
}

// ListProcess prints the running processes and records the operation.
func (c *Core) ListProcess(args []string) error {
	out, err := c.handler.ListProcess(args)

	return c.finishList(out, err, newListProcessNote(args))
}

// ListHistory prints the history table and records the operation.
func (c *Core) ListHistory() error {
	return c.finishList(c.handler.ListHistory(), nil, newListHistoriesNote())
}

// ListFolder prints the files and folders under the requested path.
func (c *Core) ListFolder(args []string) error {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	out, err := c.handler.ListFolder(currentPath)

	return c.finishList(out, err, newListFolderNote(args))
}

// ListSystem prints every system variable, one per line.
func (c *Core) ListSystem() error {
	out, err := c.handler.ListSystem()
	if err == nil {
		out = formatVariables(out)
	}

	return c.finishList(out, err, newListSystemNote())
}

// ListSysPath prints the system path variable, one entry per line.
func (c *Core) ListSysPath() error {
	out, err := c.handler.ListSysPath()
	if err == nil {
		out = formatVariable("PATH", out)
	}

	return c.finishList(out, err, newListSystemPathNote())
}

// formatVariables renders every KEY=VALUE line as a named block.
func formatVariables(out string) string {
	var builder strings.Builder

	for _, sysVar := range strings.Split(out, "\n") {
		before, after, ok := strings.Cut(sysVar, "=")
		if !ok {
			continue
		}

		builder.WriteString(formatVariable(before, after))
		builder.WriteString("\n")
	}

	return builder.String()
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

// finishList reports a listing: the result goes to the result stream on
// success, and a failure is returned instead of printed so the command exits
// non-zero.
func (c *Core) finishList(out string, err error, note string) error {
	if err != nil {
		return err
	}

	c.app.Logger.Out(out)

	c.app.History.Add(operations.List, note)
	return c.app.Save()
}
