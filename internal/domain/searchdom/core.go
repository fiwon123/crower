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

// SearchBrowser opens a search in the default browser.
func (c *Core) SearchBrowser(args []string) error {
	content := ""
	if len(args) > 0 {
		content = args[0]
	}

	out, err := c.handler.SearchBrowser(content)

	return c.finishSearch(out, err, newSearchBrowserNote(args))
}

// SearchFile searches files under a path.
func (c *Core) SearchFile(args []string) error {
	content, currentPath := assertArgs(args)

	out, err := c.handler.SearchFile(currentPath, content)

	return c.finishSearch(out, err, newSearchFileNote(args))
}

// SearchFolder searches folders under a path.
func (c *Core) SearchFolder(args []string) error {
	content, currentPath := assertArgs(args)

	out, err := c.handler.SearchFolder(currentPath, content)

	return c.finishSearch(out, err, newSearchFolderNote(args))
}

// SearchFileAndFolder searches files and folders under a path.
func (c *Core) SearchFileAndFolder(args []string) error {
	content, currentPath := assertArgs(args)

	out, err := c.handler.SearchFileAndFolder(currentPath, content)

	return c.finishSearch(out, err, newSearchFileAndFolderNote(args))
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

// finishSearch reports a search and records it.
//
// The operation recorded is Search. It was Revert for every search, which
// mislabelled the history and made `revert`/`restore` address the wrong entry.
func (c *Core) finishSearch(out string, err error, note string) error {
	if err != nil {
		return err
	}

	c.app.Logger.Out(out)

	c.app.History.Add(operations.Search, note)

	return c.app.Save()
}
