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
	payload := commandsPayload(c.app.OrderKeys, c.app.AllCommandsByName)

	return c.emit(c.handler.CommandsTable(), payload, newListCommandsNote())
}

// ListProcess prints the running processes and records the operation.
func (c *Core) ListProcess(args []string) error {
	rows, err := c.handler.ListProcess(args)
	if err != nil {
		return err
	}

	return c.emit(processTable(rows), processesPayload(rows), newListProcessNote(args))
}

// ListHistory prints the history table and records the operation.
func (c *Core) ListHistory() error {
	payload := historyPayload(c.app.History.AllData)

	return c.emit(c.handler.ListHistory(), payload, newListHistoriesNote())
}

// ListFolder prints the files and folders under the requested path.
func (c *Core) ListFolder(args []string) error {
	currentPath := "./"
	if len(args) > 0 {
		currentPath = args[0]
	}

	// The JSON result is a directory read, not a parse of `ls -la`, so it does
	// not go through the shell and cannot fail on a platform whose listing
	// command is unsupported.
	if c.app.Format == app.FormatJSON {
		payload, err := c.handler.FolderResult(currentPath)
		if err != nil {
			return err
		}

		return c.emit("", payload, newListFolderNote(args))
	}

	out, err := c.handler.ListFolder(currentPath)
	if err != nil {
		return err
	}

	return c.emit(out, nil, newListFolderNote(args))
}

// ListSystem prints every system variable, one per line.
func (c *Core) ListSystem() error {
	out, err := c.handler.ListSystem()
	if err != nil {
		return err
	}

	return c.emit(formatVariables(out), variablesPayload(out), newListSystemNote())
}

// ListSysPath prints the system path variable, one entry per line.
func (c *Core) ListSysPath() error {
	out, err := c.handler.ListSysPath()
	if err != nil {
		return err
	}

	return c.emit(formatVariable("PATH", out), variablePayload("PATH", out), newListSystemPathNote())
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

// emit writes the result in the configured format and records the operation.
//
// A failure is returned instead of printed, so the command exits non-zero; a
// result is written only once the operation succeeded.
func (c *Core) emit(text string, payload any, note string) error {
	if err := c.app.Result(text, payload); err != nil {
		return err
	}

	c.app.History.Add(operations.List, "", note)
	return c.app.Save()
}
