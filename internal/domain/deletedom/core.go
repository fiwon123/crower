package deletedom

import (
	"github.com/fiwon123/crower/internal/app"

	"github.com/fiwon123/crower/internal/errors"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/pkg/utils"
)

type Core struct {
	app     *app.Config
	handler *Handler
	input   *Input
}

func NewCore(app *app.Config, list listHandler) *Core {

	handler := NewHandler(app)
	input := newInput(app, list)

	return &Core{
		app:     app,
		handler: handler,
		input:   input,
	}
}

// Delete removes a managed command by name or alias.
func (c *Core) Delete(args []string) error {

	key := ""
	if len(args) > 0 {
		key = args[0]
	}

	ok, err := c.input.CheckDeleteInput(&key)
	if err != nil {
		return err
	}

	if !ok {
		c.app.Logger.Info("Cancelling delete...")
		return nil
	}

	deleted, err := c.performDeleteCommand(key)
	if err != nil {
		return err
	}

	c.app.History.Add(operations.Delete, newDeleteCommandNote(deleted, args))

	return c.app.Save()
}

// performDeleteCommand removes key and saves the configuration.
//
// It returns the deleted command on success and an error on failure. The old
// version returned a nil pointer for "failed", and both callers treated a
// non-nil result as "stop" -- so a successful delete was never recorded in
// history (and could not be reverted), while a failed one recorded a nil
// command and panicked when the note was built.
func (c *Core) performDeleteCommand(key string) (*command.Data, error) {

	deleted, ok := c.handler.DeleteCommand(key)
	if !ok {
		return nil, errors.CommandNotFound(key)
	}

	if err := utils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath); err != nil {
		return nil, err
	}

	c.app.Logger.Info("deleted command: ", key)

	return deleted, nil
}

// DeleteLast deletes the command touched by the last operation of the given
// kind.
func (c *Core) DeleteLast(op operations.MainOperationEnum) error {
	content := c.app.GetLast(op)

	if content == nil {
		return errors.CommandNotFound(string(op))
	}

	deleted, err := c.performDeleteCommand(content.CommandName)
	if err != nil {
		return err
	}

	c.app.History.Add(operations.Delete, newDeleteLastNote(op, deleted))

	return c.app.Save()
}

// DeleteSystemVariable removes a system variable.
func (c *Core) DeleteSystemVariable(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("var name")
	}

	out, err := c.handler.DeleteSystemVariable(args[0])
	if err != nil {
		return err
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operations.Delete, newDeleteSystemVariable(args))

	return c.app.Save()
}

// DeleteSystemPathVariable removes a folder from the system path variable.
func (c *Core) DeleteSystemPathVariable(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("path")
	}

	out, err := c.handler.DeleteSystemPathVariable(args[0])
	if err != nil {
		return err
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operations.Delete, newDeleteSystemPathVariable(args))

	return c.app.Save()
}

// DeleteFile deletes a file.
func (c *Core) DeleteFile(args []string) error {
	if len(args) == 0 {
		return errors.NotFileAndOutputPath()
	}

	if err := c.handler.DeleteFile(args[0]); err != nil {
		return err
	}

	c.app.History.Add(operations.Delete, newDeleteFileNote(args))

	return c.app.Save()
}

// DeleteFolder deletes a folder.
func (c *Core) DeleteFolder(args []string) error {
	if len(args) == 0 {
		return errors.NotFileAndOutputPath()
	}

	if err := c.handler.DeleteFolder(args[0]); err != nil {
		return err
	}

	c.app.History.Add(operations.Delete, generateDeleteFolderNote(args))

	return c.app.Save()
}

// DeleteHistoryContent drops one history entry, choosing interactively.
func (c *Core) DeleteHistoryContent(args []string) error {
	content, ok, err := c.input.CheckDeleteHistoryContentInput()
	if err != nil {
		return err
	}

	if !ok {
		c.app.Logger.Info("Cancelling Delete History Content...")
		return nil
	}

	out, err := c.handler.DeleteHistoryContent(content)
	if err != nil {
		return err
	}

	c.app.Logger.Info(out)

	return c.app.SaveOnlyHistory()
}
