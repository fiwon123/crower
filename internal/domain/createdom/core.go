package createdom

import (
	stderrors "errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/opendom"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/fiwon123/crower/internal/helper"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/pkg/utils"
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

	c.app.History.Add(operations.Create, newCreateCommandNote(created, args))

	return c.app.Save()
}

func (c *Core) performCreateCommand(name string, allAlias []string, exec string) (*command.Data, error) {
	created, err := c.handler.CreateCommand(name, allAlias, exec)
	if err != nil {
		return nil, fmt.Errorf("add command %q: %w", name, err)
	}

	if err := utils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath); err != nil {
		return nil, err
	}

	c.app.Logger.Info("added new command: ", name)

	return created, nil
}

// CreateProcess adds a managed command built from a running process.
func (c *Core) CreateProcess(name string, args []string) error {
	created, err := c.handler.CreateProcess(name, args)
	if err != nil {
		return fmt.Errorf("add command by process: %w", err)
	}

	if err := utils.WriteToml(c.app.AllCommandsByName, c.app.CfgFilePath); err != nil {
		return err
	}

	c.app.Logger.Info("added new command by process: ", name)

	c.app.History.Add(operations.Create, newCreateProcessNote(created, args))

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

	c.app.History.Add(operations.Create, newCreateSystemVariableNote(args))

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

	c.app.History.Add(operations.Create, newCreateSystemPathVariableNote(args))

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

	c.app.History.Add(operations.Create, generateCreateFile(args))

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

	c.app.History.Add(operations.Create, newCreateFolder(args))

	return c.app.Save()
}

// CreateLastCommand re-creates the command touched by the last operation of the
// given kind.
func (c *Core) CreateLastCommand(op operations.MainOperationEnum, args []string) error {
	if len(args) == 0 {
		return errors.NotEnoughArgs("name")
	}

	name := args[0]

	content := c.app.GetLast(op)
	if content == nil {
		return errors.CommandNotFound(string(op))
	}

	exec := execFromNote(content.Note)
	if exec == "" {
		return fmt.Errorf("the last %q operation carries no exec to reuse", op)
	}

	created, err := c.performCreateCommand(name, []string{}, exec)
	if err != nil {
		return err
	}

	c.app.History.Add(operations.Create, newCreateCommandLastExecuteNote(created))

	return c.app.Save()
}

// execFromNote reads the exec field out of a history note such as
// "exec=echo hi;name=hello".
//
// The field is checked for a separator first: the old code indexed the split
// result without looking, so a note carrying a bare "exec" word panicked.
func execFromNote(note string) string {
	for keyValRaw := range strings.SplitSeq(note, ";") {
		keyVal := strings.SplitN(keyValRaw, "=", 2)
		if len(keyVal) == 2 && keyVal[0] == "exec" {
			return keyVal[1]
		}
	}

	return ""
}

// CreateScriptCommand creates a script file, registers it as a command, and
// opens the folder that holds it.
func (c *Core) CreateScriptCommand(args []string) error {
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
	created, err := c.performCreateCommand(name, []string{}, scriptFilePath)
	if err != nil {
		return err
	}

	if err := c.openHandler.Open([]string{filepath.Dir(scriptFilePath)}); err != nil {
		return err
	}

	c.app.History.Add(operations.Create, newCreateScriptCommandNote(created, args))

	return c.app.Save()
}
