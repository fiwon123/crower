package cmd

import (
	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/application/resetcommands"
	"github.com/fiwon123/crower/internal/data/note"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/helper"

	"github.com/spf13/cobra"
)

// Cmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset all commands",
	Long:  `reset all commands`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := newApp(cmd)
		if err != nil {
			return err
		}

		return resetcommands.New(
			resetState{cfg: cfg},
			resetOutput{cfg: cfg, table: listdom.NewHandler(cfg)},
			resetPrompt{cfg: cfg},
		).Run()
	},
}

// resetState adapts the legacy app.Config to the use case's State port.
//
// It is the only place that knows the commit is still two writes: the command
// file first, then the history. Making that one unit of work belongs to the
// storage adapter in Phase 2, not to the use case.
type resetState struct {
	cfg *app.Config
}

func (s resetState) Commands() catalog.Catalog {
	return s.cfg.Catalog()
}

func (s resetState) Commit(commands catalog.Catalog, op operation.Operation, replay string) error {
	s.cfg.ReplaceCatalog(commands)

	if err := s.cfg.SaveCommands(); err != nil {
		return err
	}

	s.cfg.History.Add(op, "", note.New().AddMainOperation(op).AddCrowerExec(replay, nil).Build())

	return s.cfg.Save()
}

// resetOutput adapts the diagnostics stream and the table renderer.
type resetOutput struct {
	cfg   *app.Config
	table *listdom.Handler
}

func (o resetOutput) Info(message string) {
	o.cfg.Logger.Info(message)
}

// Table renders the command set the caller asked about.
//
// The renderer still reads the loaded command map until `list` is ported; the
// argument and the map describe the same commands today, so the result is
// identical. Porting `list` makes this adapter honest.
func (o resetOutput) Table(catalog.Catalog) string {
	return o.table.CommandsTable()
}

// resetPrompt adapts the terminal question.
type resetPrompt struct {
	cfg *app.Config
}

func (p resetPrompt) Confirm(question string) (bool, error) {
	return helper.GetUserConfirmation(question, p.cfg)
}
