// Package checkrelease reports whether a newer crower has been published.
package checkrelease

import (
	"github.com/fiwon123/crower/internal/domain/operation"
)

// Releases names the newest published release.
type Releases interface {
	LatestTag() (string, error)
}

// Output is the diagnostics stream. The check prints no result: it is a status
// line for a person, and a script reads the exit status.
type Output interface {
	Info(message string, data ...any)
}

// History records the check.
//
// It is separate from the command-set state on purpose: a check changes no
// command, so it must not rewrite the command file.
type History interface {
	Record(op operation.Operation, commandName, note string) error
}

// Notes renders the history note for the check.
type Notes interface {
	Check() string
}

// UseCase compares the running version with the published one.
type UseCase struct {
	releases Releases
	output   Output
	history  History
	notes    Notes
}

// New builds the use case from its four dependencies.
func New(releases Releases, output Output, history History, notes Notes) *UseCase {
	return &UseCase{releases: releases, output: output, history: history, notes: notes}
}

// Run reports the current version and, when they differ, the published one.
//
// Being up to date is not recorded: it changed nothing, and a history entry
// per `--check` would fill the timeline with non-events.
func (useCase *UseCase) Run(currentVersion string) error {
	newVersion, err := useCase.releases.LatestTag()
	if err != nil {
		return err
	}

	useCase.output.Info("Current Version: ", currentVersion)

	if newVersion == currentVersion {
		useCase.output.Info("already up-to-date")
		return nil
	}

	useCase.output.Info("New Version Found: ", newVersion)
	useCase.output.Info("Check: https://github.com/fiwon123/crower/releases/latest")

	return useCase.history.Record(operation.Check, "", useCase.notes.Check())
}
