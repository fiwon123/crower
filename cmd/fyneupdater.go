package cmd

import (
	"os"
	"path/filepath"
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

	// installer installs the window's own member. update.New() extracts the CLI,
	// and selfupdate writes those bytes over crower-gui, so Restart launched a
	// program that prints help and exits instead of opening the window (#489,
	// #513).
	installer *update.Updater

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
		installer:  update.NewWindow(),
		executable: executable,
		args:       os.Args[1:],
		start:      startReplacement,
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

	installer := u.installer
	if installer == nil {
		installer = update.NewWindow()
	}

	return upgradeapp.New(session, quietOutput{}, latestTag{}, installer, session, session, session).Run(
		upgradeapp.Request{Current: u.current, Version: version, AssumeYes: true},
	)
}

// Restart starts the binary on disk again. On Unix it replaces this process
// image; on Windows it spawns a new process. The caller has already closed the
// window.
func (u fyneUpdater) Restart() error {
	start := u.start
	if start == nil {
		start = startReplacement
	}

	err := start(u.executable, u.args)
	if err == nil {
		return nil
	}

	// There is no window left to report to; leave a note next to the
	// configuration so a failed restart can still be diagnosed.
	if u.cfg != nil {
		_ = os.WriteFile(filepath.Join(u.cfg.ConfigFolder(), "restart-error.log"), []byte(err.Error()+"\n"), 0o644)
	}

	return err
}

// cleanExecutable drops the " (deleted)" suffix Linux appends to
// /proc/self/exe once the running binary has been replaced by an update.
func cleanExecutable(executable string) string {
	return strings.TrimSuffix(executable, " (deleted)")
}

// replacementArgv is the argv for the replacement: the executable first, then
// the arguments this process was started with.
func replacementArgv(executable string, args []string) []string {
	return append([]string{cleanExecutable(executable)}, args...)
}

// quietOutput swallows the use case's progress lines: the window shows them in
// its own status label, and a desktop launch has no terminal to print to.
type quietOutput struct{}

func (quietOutput) Info(string, ...any) {}
