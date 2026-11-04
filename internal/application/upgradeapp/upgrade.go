// Package upgradeapp installs a crower release: the newest, a named one, or the
// previous one.
//
// The decision is small; the technology -- the release download, its checksum,
// the archive and the self-update -- is behind the Updater port.
package upgradeapp

import (
	"fmt"

	"github.com/fiwon123/crower/internal/domain/operation"
)

// Request is one update.
type Request struct {
	// Current is the version this run is executed from.
	Current string

	// Version installs that release explicitly, including an older one.
	Version string

	// Rollback installs the version installed before Current.
	Rollback bool

	// AssumeYes skips the confirmation. A script that wants no prompt asks for
	// it instead of relying on --no-input, which would fail the run.
	AssumeYes bool
}

// State is the history the update is recorded in.
type State interface {
	// Record appends a history entry and persists the index.
	Record(op operation.Operation, commandName, note string) error
}

// Output is the diagnostics stream.
type Output interface {
	Info(message string, data ...any)
}

// Releases reports the newest published tag.
type Releases interface {
	LatestTag() (string, error)
}

// Updater installs one release.
type Updater interface {
	// Upgrade downloads, verifies and applies the release with this version.
	Upgrade(version string) error
}

// Versions remembers what the installation has run.
type Versions interface {
	// Previous returns the most recently installed version other than current.
	Previous(current string) (string, bool)

	// RecordInstalled puts version at the front of the list.
	RecordInstalled(version string) error
}

// Prompter asks the confirmation the update needs.
type Prompter interface {
	Confirm(question string) (bool, error)
}

// Notes renders the history note for the update.
type Notes interface {
	UpgradeNote() string
}

// UseCase runs one update.
type UseCase struct {
	state    State
	output   Output
	releases Releases
	updater  Updater
	versions Versions
	prompt   Prompter
	notes    Notes
}

// New builds the use case from its seven dependencies.
func New(state State, output Output, releases Releases, updater Updater, versions Versions, prompt Prompter, notes Notes) *UseCase {
	return &UseCase{
		state:    state,
		output:   output,
		releases: releases,
		updater:  updater,
		versions: versions,
		prompt:   prompt,
		notes:    notes,
	}
}

// Run installs the requested release.
//
// The newest release is only installed on a yes; a named version and a rollback
// are asked about too, because replacing the running binary is never something
// a user should discover afterwards.
func (u *UseCase) Run(request Request) error {
	u.output.Info("Current Version: ", request.Current)

	target, err := u.target(request)
	if err != nil {
		return err
	}

	if target == "" {
		return nil
	}

	if !request.AssumeYes {
		proceed, err := u.prompt.Confirm(fmt.Sprintf("Install %s?", target))
		if err != nil {
			return err
		}

		if !proceed {
			u.output.Info("Cancelling update...")
			return nil
		}
	}

	if err := u.updater.Upgrade(target); err != nil {
		return err
	}

	if err := u.versions.RecordInstalled(target); err != nil {
		return err
	}

	u.output.Info("crower installed ", target, " (was ", request.Current, ")")

	return u.state.Record(operation.Upgrade, "", u.notes.UpgradeNote())
}

// target decides which release to install; "" means there is nothing to do.
func (u *UseCase) target(request Request) (string, error) {
	switch {
	case request.Version != "":
		return request.Version, nil

	case request.Rollback:
		previous, found := u.versions.Previous(request.Current)
		if !found {
			return "", fmt.Errorf("no previous version recorded; use --install-version <tag>")
		}

		return previous, nil

	default:
		latest, err := u.releases.LatestTag()
		if err != nil {
			return "", err
		}

		if latest == request.Current {
			u.output.Info("already up-to-date")
			return "", nil
		}

		return latest, nil
	}
}
