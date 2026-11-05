package cmd

import (
	"os"
	osexec "os/exec"
	"time"

	"github.com/fiwon123/crower/internal/adapters/github"
	"github.com/fiwon123/crower/internal/adapters/storage"
	"github.com/fiwon123/crower/internal/adapters/update"
	"github.com/fiwon123/crower/internal/application/upgradeapp"
	"github.com/fiwon123/crower/internal/session"
)

// updateCheckTimeout bounds the window's release lookup. Opening crower must
// not depend on GitHub answering, so a slow network is a status line, not a
// frozen window.
const updateCheckTimeout = 5 * time.Second

// fyneUpdater answers the About tab's version questions and installs releases
// through the same use case `crower --upgrade` runs.
//
// It has no Fyne import, so the adapter is covered by the normal cmd tests.
type fyneUpdater struct {
	cfg     *session.Session
	current string

	// start is the process spawn. A test replaces it so a restart assertion
	// does not launch the test binary.
	start func() error
}

func newFyneUpdater(cfg *session.Session, current string) fyneUpdater {
	return fyneUpdater{cfg: cfg, current: current, start: startReplacementProcess}
}

// Current is the version this window was built from.
func (u fyneUpdater) Current() string {
	return u.current
}

// Latest is the newest published tag.
func (u fyneUpdater) Latest() (string, error) {
	return github.LatestTagTimeout(updateCheckTimeout)
}

// Previous is the version installed before this one, when one was recorded.
func (u fyneUpdater) Previous() (string, bool) {
	return storage.NewVersions(u.cfg.ConfigFolder()).Previous(u.current)
}

// Install downloads, verifies and applies one release.
//
// The window has already asked, so the use case runs with AssumeYes; its
// progress lines go to a quiet output because the window reports them itself.
func (u fyneUpdater) Install(version string) error {
	session := newUpdateSession(u.cfg)

	return upgradeapp.New(session, quietOutput{}, latestTag{}, update.New(), session, session, session).Run(
		upgradeapp.Request{Current: u.current, Version: version, AssumeYes: true},
	)
}

// Restart starts the binary on disk again. The caller closes the window so the
// old process exits.
func (u fyneUpdater) Restart() error {
	if u.start == nil {
		return startReplacementProcess()
	}

	return u.start()
}

// quietOutput swallows the use case's progress lines: the window shows them in
// its own status label, and a desktop launch has no terminal to print to.
type quietOutput struct{}

func (quietOutput) Info(string, ...any) {}

// startReplacementProcess runs the executable at the current path with the
// same arguments. After a self-update that path holds the new binary.
func startReplacementProcess() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	command := osexec.Command(executable, os.Args[1:]...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Start()
}
