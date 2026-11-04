// Package prompt asks the user the questions a command needs.
//
// It is an adapter: the read, the diagnostic stream and --no-input come from
// the composition root, and the two tables a row question can refer to -- the
// command set and the history -- are passed in as values. Nothing here knows
// about the configuration or the terminal beyond those.
package prompt

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/fiwon123/crower/internal/errors"
)

// Prompt asks questions on one terminal.
type Prompt struct {
	read     func() (string, error)
	info     func(string, ...any)
	noInput  bool
	commands catalog.Catalog
	history  timeline.Log
}

// New builds a Prompt over a line reader, a diagnostic writer, the flag that
// forbids asking, and the two tables a row question can address.
func New(read func() (string, error), info func(string, ...any), noInput bool, commands catalog.Catalog, history timeline.Log) *Prompt {
	return &Prompt{read: read, info: info, noInput: noInput, commands: commands, history: history}
}

// ask prints a question and reads until the answer validates.
//
// With --no-input the ask fails before it is printed: a script gets a clear
// error instead of a wait.
func (p *Prompt) ask(question string, validate func(string) error) (string, error) {
	for {
		if p.noInput {
			return "", errors.InputRequired(question)
		}

		p.info(question + ": ")

		input, err := p.read()
		if err != nil {
			return "", err
		}

		if err := validate(input); err != nil {
			p.info(err.Error())
			continue
		}

		return input, nil
	}
}

// Confirm asks a yes/no question.
func (p *Prompt) Confirm(question string) (bool, error) {
	for {
		if p.noInput {
			return false, errors.InputRequired(question)
		}

		p.info(question + " ([Y]es/[N]o): ")

		input, err := p.read()
		if err != nil {
			return false, err
		}

		switch strings.TrimSpace(input) {
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		default:
			p.info(errors.InvalidInput().Error())
		}
	}
}

// Ask asks a question and requires a non-empty answer.
func (p *Prompt) Ask(question string) (string, error) {
	return p.ask(question, func(input string) error {
		if input == "" {
			return errors.EmptyInput()
		}

		return nil
	})
}

// Aliases asks for aliases until the user enters an empty line.
func (p *Prompt) Aliases() ([]string, error) {
	output := []string{}

	alias := "none"
	for alias != "" {
		p.info("current aliases: ", output)

		value, err := p.ask("Add new alias (type enter to skip): ", func(input string) error {
			for _, r := range input {
				if !unicode.IsNumber(r) && !unicode.IsLetter(r) {
					return errors.OnlyLettersAndNumbers()
				}
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		alias = value
		if alias != "" {
			output = append(output, alias)
		}
	}

	return output, nil
}

// ChooseCommand maps a printed row number to the command name on that row.
func (p *Prompt) ChooseCommand() (string, error) {
	value, err := p.ask("Select Row", func(input string) error {
		index, err := rowIndex(input)
		if err != nil {
			return err
		}

		if index < 0 || index >= p.commands.Len() {
			return errors.InvalidRows()
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	index, _ := rowIndex(value)

	return p.commands.Names()[index], nil
}

// ChooseHistoryEntry maps a printed row number to a history entry.
//
// The history is printed newest first, so row 0 is the last entry.
func (p *Prompt) ChooseHistoryEntry() (timeline.Content, error) {
	value, err := p.ask("Select Row", func(input string) error {
		index, err := rowIndex(input)
		if err != nil {
			return err
		}

		if index < 0 || index >= len(p.history.AllData) {
			return errors.InvalidRows()
		}

		return nil
	})
	if err != nil {
		return timeline.Content{}, err
	}

	index, _ := rowIndex(value)

	return p.history.AllData[len(p.history.AllData)-1-index], nil
}

func rowIndex(input string) (int, error) {
	index, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, errors.InvalidRows()
	}

	return index, nil
}

// String renders a value the way a diagnostic line shows it.
func String(value any) string {
	return fmt.Sprint(value)
}
