package searchdom

import (
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

func (c *Core) SearchBrowser(args []string) {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	c.handler.SearchBrowser(content)

	c.app.History.Add(operations.Revert, newSearchBrowserNote(args))
	c.app.Save()
}

func (c *Core) SearchFile(args []string) {
	content, currentPath := assertArgs(args)

	out, err := c.handler.SearchFile(currentPath, content)
	c.assertSearchResult(out, err)

	c.app.History.Add(operations.Revert, newSearchFileNote(args))
	c.app.Save()
}

func (c *Core) SearchFolder(args []string) {
	content, currentPath := assertArgs(args)

	out, err := c.handler.SearchFolder(currentPath, content)
	c.assertSearchResult(out, err)

	c.app.History.Add(operations.Revert, newSearchFolderNote(args))
	c.app.Save()
}

func (c *Core) SearchFileAndFolder(args []string) {
	content, currentPath := assertArgs(args)

	out, err := c.handler.SearchFileAndFolder(currentPath, content)
	c.assertSearchResult(out, err)

	c.app.History.Add(operations.Revert, newSearchFileAndFolderNote(args))
	c.app.Save()
}

func assertArgs(args []string) (string, string) {
	content := ""
	currentPath := "."

	if len(args) > 1 {
		content = args[0]
		currentPath = args[1]
	}

	return content, currentPath
}

func (c *Core) assertSearchResult(out string, err error) {
	if err != nil {
		c.app.Logger.Error("Error trying to search: ", out, err)
	}

	c.app.Logger.Info(out)
}
