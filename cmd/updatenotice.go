package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fiwon123/crower/internal/adapters/github"
	"github.com/fiwon123/crower/internal/adapters/storage"
)

const (
	// noticeDaily and noticeWeekly are how often the periodic check may look.
	noticeDaily  = 24 * time.Hour
	noticeWeekly = 7 * 24 * time.Hour

	// noticeTimeout bounds the check: a notice that delays every command must
	// fail fast, not wait for the network.
	noticeTimeout = 2 * time.Second
)

// noticeDue reports whether the periodic notice may look for a newer release.
//
// mode is the setting: "never" turns it off, "weekly" slows it down, anything
// else (including empty) is the daily default. A last check that cannot be read
// is treated as "never checked", which errs toward telling the user once.
func noticeDue(mode, lastCheck string, now time.Time) bool {
	period := noticeDaily

	switch mode {
	case "never":
		return false
	case "weekly":
		period = noticeWeekly
	}

	if lastCheck == "" {
		return true
	}

	last, err := time.Parse(time.RFC3339, lastCheck)
	if err != nil {
		return true
	}

	return now.Sub(last) >= period
}

// noticeLine is the one-line notice, or "" when there is nothing to say.
func noticeLine(current, latest string) string {
	if latest == "" || latest == current {
		return ""
	}

	return fmt.Sprintf("crower %s is available (you have %s) — run crower --upgrade", latest, current)
}

// maybeNotifyUpdate looks for a newer release at most once per period and says
// so on stderr.
//
// It is best effort: no terminal, --no-input, an explicit update flag or a
// failed lookup all end it silently. The lookup never blocks a command for long
// and never changes the exit status.
func maybeNotifyUpdate() {
	// The update flags already do the real work; do not add a second lookup.
	if versionFlag || checkNewVersion || upgradeFlag || rollbackFlag || installVersion != "" {
		return
	}

	noInput, err := rootCmd.PersistentFlags().GetBool(flagNoInput)
	if err != nil || noInput || !stderrIsTerminal() {
		return
	}

	cfgPath, err := rootCmd.PersistentFlags().GetString(flagConfig)
	if err != nil {
		return
	}

	if cfgPath == "" {
		if cfgPath, err = defaultConfigPath(); err != nil {
			return
		}
	}

	settings := storage.NewSettings(filepath.Dir(cfgPath))
	data := settings.Load()

	now := time.Now()
	if !noticeDue(data.UpdateCheck, data.LastCheck, now) {
		return
	}

	latest, lookupErr := github.LatestTagTimeout(noticeTimeout)

	// Record the look either way: a failing network must not make every run
	// retry it.
	data.LastCheck = now.Format(time.RFC3339)
	if lookupErr == nil && latest != "" {
		data.Latest = latest
	}

	_ = settings.Save(data)

	if lookupErr != nil {
		return
	}

	if line := noticeLine(Version, latest); line != "" {
		fmt.Fprintln(streams.Err)
		fmt.Fprintln(streams.Err, line)
	}
}

// stderrIsTerminal reports whether the diagnostics stream is a terminal. A
// notice belongs on an interactive terminal, not in a pipe or CI log.
//
// /dev/null is a character device too, so the char-device bit alone would let a
// redirected run write a settings file and spend the lookup. On Linux the
// descriptor's target settles it; elsewhere the char-device answer stands.
func stderrIsTerminal() bool {
	info, err := os.Stderr.Stat()
	if err != nil {
		return false
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		return false
	}

	if target, err := os.Readlink("/proc/self/fd/2"); err == nil && target == os.DevNull {
		return false
	}

	return true
}
