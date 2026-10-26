package updatedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/data/operations"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/fiwon123/crower/internal/interfaces"
)

type Input struct {
	app         *app.Config
	listHandler interfaces.ListHandler
}

func newInput(app *app.Config, listHandler interfaces.ListHandler) *Input {
	return &Input{
		app:         app,
		listHandler: listHandler,
	}
}

// CheckUpdateInput asks for the row to update and for any field left unset.
//
// It returns false with a nil error when the user declines the final
// confirmation: declining is a choice, not a failure.
func (i *Input) CheckUpdateInput(key *string, name *string, allAlias *[]string, exec *string) (bool, error) {

	if *key == "" {
		i.app.Logger.Info(i.listHandler.CommandsTable())

		value, err := helper.GetUserInput("Select Row", helper.IsValidInputKey, i.app)
		if err != nil {
			return false, err
		}

		*key = value.(string)
	}

	updateCommand := i.app.AllCommandsByName.Get(*key)
	if updateCommand == nil {
		return false, errors.CommandNotFound(*key)
	}

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Name:    ", updateCommand.Name)
	i.app.Logger.Info("Aliases: ", updateCommand.AllAlias)
	i.app.Logger.Info("Exec:    ", updateCommand.Exec)
	i.app.Logger.Info("")

	if *name == "" {
		ok, err := helper.GetUserConfirmation("Do you want to update name", i.app)
		if err != nil {
			return false, err
		}

		if ok {
			*name, err = helper.InputName(i.app)
			if err != nil {
				return false, err
			}
		}
	}

	if len(*allAlias) == 0 {
		ok, err := helper.GetUserConfirmation("Do you want to update alias", i.app)
		if err != nil {
			return false, err
		}

		if ok {
			*allAlias, err = helper.InputAlias(i.app)
			if err != nil {
				return false, err
			}
		}
	}

	if *exec == "" {
		ok, err := helper.GetUserConfirmation("Do you want to update exec", i.app)
		if err != nil {
			return false, err
		}

		if ok {
			*exec, err = helper.InputExec(i.app)
			if err != nil {
				return false, err
			}
		}
	}

	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("Old Command: ")
	i.app.Logger.Info("Name:    ", updateCommand.Name)
	i.app.Logger.Info("Aliases: ", updateCommand.AllAlias)
	i.app.Logger.Info("Exec:    ", updateCommand.Exec)
	i.app.Logger.Info("-----------------------------------------")
	i.app.Logger.Info("New Command: ")
	i.app.Logger.Info("Name:    ", *name)
	i.app.Logger.Info("Aliases: ", *allAlias)
	i.app.Logger.Info("Exec:    ", *exec)
	i.app.Logger.Info("")

	return helper.GetUserConfirmation("Continue to update", i.app)
}

// commandNotFound reports a last-operation lookup that found nothing.
func commandNotFound(op operations.MainOperationEnum) error {
	return fmt.Errorf("no %q command in history", op)
}
