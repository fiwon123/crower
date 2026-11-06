package cmd

import (
	"os"
	osexec "os/exec"
	"strings"
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

	// executable and args are captured before any install: selfupdate moves the
	// running binary aside, so os.Executable() after the update points at the old
	// file (or a deleted one) instead of the new binary at the target path.
	executable string
	args       []string

	// start is the process spawn. A test replaces it so a restart assertion
	// does not launch the test binary.
	start func(executable string, args []string) error
}

func newFyneUpdater(cfg *session.Session, current string) fyneUpdater {
	executable, _ := os.Executable()

	return fyneUpdater{
		cfg:        cfg,
		current:    current,
		executable: executable,
		args:       os.Args[1:],
		start:      startReplacementProcess,
	}
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
	start := u.start
	if start == nil {
		start = startReplacementProcess
	}

	return start(u.executable, u.args)
}

// quietOutput swallows the use case's progress lines: the window shows them in
// its own status label, and a desktop launch has no terminal to print to.
type quietOutput struct{}

func (quietOutput) Info(string, ...any) {}

// replacementCommand builds the process that runs the binary at the path the
// window was started from. The path is captured before an update because
// selfupdate renames the running binary aside: on Linux /proc/self/exe then
// resolves to the old file, or to "<path> (deleted)" once it is removed.
func replacementCommand(executable string, args []string) (*osexec.Cmd, error) {
	if executable == "" {
		var err error

		executable, err = os.Executable()
		if err != nil {
			return nil, err
		}
	}

	executable = strings.TrimSuffix(executable, " (deleted)")

	command := osexec.Command(executable, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command, nil
}

// startReplacementProcess runs the executable at the current path with the
// same arguments. After a self-update that path holds the new binary.
func startReplacementProcess(executable string, args []string) error {
	command, err := replacementCommand(executable, args)
	if err != nil {
		return err
	}

	return command.Start()
}
