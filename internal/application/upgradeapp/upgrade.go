// Package upgradeapp replaces the running binary with the newest release.
//
// The decision is four lines of flow; the technology -- the release download,
// its checksum, the archive and the self-update -- is behind the Updater port.
package upgradeapp

import (
	"github.com/fiwon123/crower/internal/domain/operation"
)

// State is the history the upgrade is recorded in.
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

// Notes renders the history note for the upgrade.
type Notes interface {
	UpgradeNote() string
}

// UseCase runs one upgrade.
type UseCase struct {
	state    State
	output   Output
	releases Releases
	updater  Updater
	notes    Notes
}

// New builds the use case from its five dependencies.
func New(state State, output Output, releases Releases, updater Updater, notes Notes) *UseCase {
	return &UseCase{state: state, output: output, releases: releases, updater: updater, notes: notes}
}

// Run upgrades from currentVersion when a newer release exists.
func (u *UseCase) Run(currentVersion string) error {
	latest, err := u.releases.LatestTag()
	if err != nil {
		return err
	}

	u.output.Info("Current Version: ", currentVersion)

	if latest == currentVersion {
		u.output.Info("already up-to-date")
		return nil
	}

	if err := u.updater.Upgrade(latest); err != nil {
		return err
	}

	u.output.Info("crower upgraded from ", currentVersion, " to ", latest)

	return u.state.Record(operation.Upgrade, "", u.notes.UpgradeNote())
}
