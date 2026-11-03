package createdom

import (
	stderrors "errors"
	"fmt"
	"path/filepath"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/legacy/opendom"

	"github.com/fiwon123/crower/internal/helper"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/operation"
)

type Core struct {
	app         *app.Config
	handler     *Handler
	openHandler opendom.Handler
}

func NewCore(app *app.Config, openHandler opendom.Handler) *Core {

	handler := NewHandler(app)

	return &Core{
		handler:     handler,
		app:         app,
		openHandler: openHandler,
	}
}

// CheckCreateInput asks for the fields a create command was not given.
func (c *Core) CheckCreateInput(name *string, alias *[]string, exec *string) error {
	if *name == "" {
		value, err := helper.GetUserInput("New Name ", helper.IsValidInput, c.app)
		if err != nil {
			return err
		}

		*name = value.(string)
	}

	if len(*alias) == 0 {
		ok, err := helper.GetUserConfirmation("Do you want to add alias", c.app)
		if err != nil {
			return err
		}

		if ok {
			*alias, err = helper.InputAlias(c.app)
			if err != nil {
				return err
			}
		}
	}

	if *exec == "" {
		value, err := helper.GetUserInput("New Exec ", helper.IsValidInput, c.app)
		if err != nil {
			return err
		}

		*exec = value.(string)
	}

	return nil
}

// CreateCommand adds a managed command.
func (c *Core) CreateCommand(allAlias []string, args []string) error {
	name := ""
	exec := ""
	if len(args) == 2 {
		name = args[0]
		exec = args[1]
	} else {
		if err := c.CheckCreateInput(&name, &allAlias, &exec); err != nil {
			return err
		}
	}

	created, err := c.performCreateCommand(name, allAlias, exec)
	if err != nil {
		return err
	}

	c.app.History.Add(operation.Create, created.Name, newCreateCommandNote(created))

	return c.app.Save()
}

func (c *Core) performCreateCommand(name string, allAlias []string, exec string) (*command.Data, error) {
	created, err := c.handler.CreateCommand(name, allAlias, exec)
	if err != nil {
		return nil, fmt.Errorf("add command %q: %w", name, err)
	}

	if err := c.app.SaveCommands(); err != nil {
		return nil, err
	}

	c.app.Logger.Info("added new command: ", name)

	return created, nil
}

// CreateProcess adds a managed command built from a running process.
func (c *Core) CreateProcess(name string, allAlias []string, args []string) error {
	created, err := c.handler.CreateProcess(name, allAlias, args)
	if err != nil {
		return fmt.Errorf("add command by process: %w", err)
	}

	if err := c.app.SaveCommands(); err != nil {
		return err
	}

	c.app.Logger.Info("added new command by process: ", name)

	c.app.History.Add(operation.Create, created.Name, newCreateProcessNote(created, args))

	return c.app.Save()
}

// CreateSystemVariable adds a system variable.
func (c *Core) CreateSystemVariable(args []string) error {
	if len(args) < 2 {
		return errors.NotEnoughArgs("var name and var value")
	}

	newVar := args[0]
	value := args[1]

	out, err := c.handler.CreateSystemVariable(newVar, value)
	if err != nil {
		return err
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operation.Create, "", newCreateSystemVariableNote(args))

	return c.app.Save()
}

// CreateSystemPathVariable adds a folder to the system path variable.
func (c *Core) CreateSystemPathVariable(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("path")
	}

	out, err := c.handler.CreateSystemPathVariable(args[0])
	if err != nil {
		return err
	}

	c.app.Logger.Info(out)

	c.app.History.Add(operation.Create, "", newCreateSystemPathVariableNote(args))

	return c.app.Save()
}

// CreateFile creates every path, reporting all failures together.
func (c *Core) CreateFile(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("file path")
	}

	var failures []error

	for _, path := range args {
		if err := c.handler.CreateFile(path); err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	c.app.History.Add(operation.Create, "", generateCreateFile(args))

	return c.app.Save()
}

// CreateFolder creates every path, reporting all failures together.
func (c *Core) CreateFolder(args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("folder path")
	}

	var failures []error

	for _, path := range args {
		if err := c.handler.CreateFolder(path); err != nil {
			failures = append(failures, err)
		}
	}

	if len(failures) > 0 {
		return stderrors.Join(failures...)
	}

	c.app.History.Add(operation.Create, "", newCreateFolder(args))

	return c.app.Save()
}

// CreateLastCommand re-creates the command touched by the last operation of the
// given kind.
//
// It takes exactly one argument, the new name. Extra words used to be ignored,
// so `create --execute name typo` succeeded and dropped the typo; and the
// --alias flag was accepted but never reached the command, because the aliases
// were hard-coded to an empty slice.
func (c *Core) CreateLastCommand(op operation.Operation, allAlias []string, args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("name")
	}

	if len(args) > 1 {
		return errors.TooManyArgs("one name")
	}

	name := args[0]

	content := c.app.GetLast(op)
	if content == nil {
		return errors.CommandNotFound(string(op))
	}

	// A command with steps never runs its exec, so reusing that value would
	// rebuild a command that does something else. Say so instead.
	steps, hasSteps := note.CommandSteps(content.Note)

	// Notes written before command_steps existed do not carry it. The command
	// may still be in the map, where its steps are visible.
	if !hasSteps && content.CommandName != "" {
		if cmd := c.app.AllCommandsByName.Get(content.CommandName); cmd != nil && len(cmd.Steps) > 0 {
			steps, hasSteps = cmd.Steps, true
		}
	}

	if hasSteps {
		if len(steps) == 0 {
			return fmt.Errorf("the last %q operation (%q) has steps that cannot be read; create --execute cannot rebuild them", op, content.CommandName)
		}

		return fmt.Errorf("the last %q operation (%q) ran %d step(s); create --execute cannot rebuild them", op, content.CommandName, len(steps))
	}

	exec := note.CommandExec(content.Note)
	if exec == "" {
		return fmt.Errorf("the last %q operation carries no exec to reuse", op)
	}

	created, err := c.performCreateCommand(name, allAlias, exec)
	if err != nil {
		return err
	}

	c.app.History.Add(operation.Create, created.Name, newCreateCommandLastExecuteNote(created))

	return c.app.Save()
}

// CreateScriptCommand creates a script file, registers it as a command, and
// opens the folder that holds it.
func (c *Core) CreateScriptCommand(allAlias []string, args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("name")
	}

	name := args[0]

	scriptFilePath, err := c.handler.CreateScriptCommand(name)
	if err != nil {
		return err
	}

	// The handler only returns a path on a platform it can build one for, so
	// there is no per-platform choice left to make here. The old switch
	// repeated the same call for two systems and returned success on any
	// other, which is how `create --script` could report success and create
	// nothing.
	created, err := c.performCreateCommand(name, allAlias, scriptFilePath)
	if err != nil {
		return err
	}

	if err := c.openHandler.Open([]string{filepath.Dir(scriptFilePath)}); err != nil {
		return err
	}

	c.app.History.Add(operation.Create, created.Name, newCreateScriptCommandNote(created))

	return c.app.Save()
}
