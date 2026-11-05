//go:build fyne

package fyneui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	// projectURL is where the About tab sends a reader.
	projectURL = "https://github.com/fiwon123/crower"

	// projectSummary is the one-line description the window shows.
	projectSummary = "A terminal file manager and command-alias manager with undoable history and per-item notes."
)

// Updater answers version questions and installs a release.
//
// It is declared here, at the point of use: the window owns the confirmation
// and the restart, and the command layer owns the release download, its
// checksum and the self-update.
type Updater interface {
	// Current is the version this window runs from.
	Current() string

	// Latest reports the newest published release.
	Latest() (string, error)

	// Previous reports the version installed before the current one, if the
	// installation has recorded one.
	Previous() (string, bool)

	// Install downloads, verifies and applies one release.
	Install(version string) error

	// Restart starts the binary on disk again as a new process.
	Restart() error
}

// aboutTab is the window's identity and its update controls.
func (u *ui) aboutTab() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(logoResource())
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(96, 96))

	name := widget.NewLabelWithStyle("crower", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	version := widget.NewLabelWithStyle("Version "+u.updater.Current(), fyne.TextAlignCenter, fyne.TextStyle{})

	summary := widget.NewLabel(projectSummary)
	summary.Alignment = fyne.TextAlignCenter
	summary.Wrapping = fyne.TextWrapWord

	var homepage fyne.CanvasObject

	if parsed, err := url.Parse(projectURL); err == nil {
		link := widget.NewHyperlink(projectURL, parsed)
		link.Alignment = fyne.TextAlignCenter
		homepage = link
	} else {
		// The constant cannot be malformed; if it ever is, a plain label beats
		// a link that goes nowhere.
		label := widget.NewLabel(projectURL)
		label.Alignment = fyne.TextAlignCenter
		homepage = label
	}

	u.updateStatus = widget.NewLabel("Updates are checked when the window opens.")
	u.updateStatus.Wrapping = fyne.TextWrapWord

	u.checkButton = widget.NewButtonWithIcon("Check for updates", theme.ViewRefreshIcon(), u.checkForUpdates)
	u.installButton = widget.NewButtonWithIcon("Install latest", theme.DownloadIcon(), u.installLatest)

	actions := container.NewHBox(u.checkButton, u.installButton)

	if previous, found := u.updater.Previous(); found {
		u.restoreButton = widget.NewButton("Restore previous ("+previous+")", u.restorePrevious)
		actions.Add(u.restoreButton)
	}

	restart := widget.NewButtonWithIcon("Restart", theme.ViewRefreshIcon(), u.restart)

	updates := widget.NewCard(
		"Updates",
		"Checks GitHub for the newest release. Nothing is replaced without confirmation.",
		container.NewVBox(
			u.updateStatus,
			actions,
			container.NewHBox(restart),
		),
	)

	body := container.NewVBox(
		container.NewCenter(logo),
		name,
		version,
		summary,
		container.NewCenter(homepage),
		widget.NewSeparator(),
		updates,
	)

	return container.NewVScroll(body)
}

// checkForUpdates is the startup courtesy: one lookup, in the background, with
// the window already open.
func (u *ui) checkForUpdates() {
	if u.updateStatus == nil || u.checkButton == nil {
		return
	}

	u.setBusy("Checking for updates…")

	go func() {
		latest, err := u.updater.Latest()

		fyne.Do(func() {
			u.setIdle()
			u.applyLatest(latest, err)
		})
	}()
}

// applyLatest is the state change without the goroutine, so a test can call it
// with a canned answer.
func (u *ui) applyLatest(latest string, err error) {
	if err != nil {
		u.latest = ""
		u.updateStatus.SetText("Could not check for updates: " + err.Error())
		return
	}

	u.latest = latest

	switch {
	case latest == "":
		u.updateStatus.SetText("No release information was returned.")
	case latest == u.updater.Current():
		u.updateStatus.SetText("crower " + u.updater.Current() + " is the newest release.")
	default:
		u.updateStatus.SetText("crower " + latest + " is available (you have " + u.updater.Current() + ").")
	}
}

func (u *ui) setBusy(status string) {
	u.updateStatus.SetText(status)
	u.checkButton.Disable()
	u.installButton.Disable()

	if u.restoreButton != nil {
		u.restoreButton.Disable()
	}
}

func (u *ui) setIdle() {
	u.checkButton.Enable()
	u.installButton.Enable()

	if u.restoreButton != nil {
		u.restoreButton.Enable()
	}
}

// installLatest looks the newest release up again before asking, so a stale
// startup answer cannot install the wrong tag.
func (u *ui) installLatest() {
	u.setBusy("Looking up the newest release…")

	go func() {
		latest, err := u.updater.Latest()

		fyne.Do(func() {
			u.setIdle()

			if err != nil {
				u.applyLatest("", err)
				return
			}

			u.applyLatest(latest, nil)

			if latest == "" || latest == u.updater.Current() {
				return
			}

			u.confirmInstall(latest)
		})
	}()
}

func (u *ui) confirmInstall(version string) {
	question := "Install crower " + version + "?\n\nIt replaces the binary on disk; a restart finishes the update."

	dialog.ShowConfirm("Install "+version, question, func(ok bool) {
		if !ok {
			return
		}

		u.installInBackground(version)
	}, u.window)
}

// restorePrevious installs the version the installation ran before this one.
func (u *ui) restorePrevious() {
	previous, found := u.updater.Previous()
	if !found {
		u.updateStatus.SetText("No previous version is recorded yet; install one with crower --install-version.")
		return
	}

	question := "Install crower " + previous + " again?\n\nIt replaces the binary on disk; a restart finishes the change."

	dialog.ShowConfirm("Restore "+previous, question, func(ok bool) {
		if !ok {
			return
		}

		u.installInBackground(previous)
	}, u.window)
}

func (u *ui) installInBackground(version string) {
	u.setBusy("Installing crower " + version + "…")

	go func() {
		err := u.installNow(version)

		fyne.Do(func() {
			u.setIdle()

			if err != nil {
				return
			}

			u.offerRestart(version)
		})
	}()
}

// installNow is the state change without the goroutine: the wrapper adds the
// busy state and the restart question.
func (u *ui) installNow(version string) error {
	if err := u.updater.Install(version); err != nil {
		u.updateStatus.SetText("Install failed: " + err.Error())
		return err
	}

	u.updateStatus.SetText("Installed crower " + version + ". Restart to use it.")

	return nil
}

func (u *ui) offerRestart(version string) {
	question := "crower " + version + " is installed. Restart the window to use it?"

	dialog.ShowConfirm("Restart now?", question, func(ok bool) {
		if !ok {
			return
		}

		u.restart()
	}, u.window)
}

func (u *ui) restart() {
	if err := u.updater.Restart(); err != nil {
		showError("Restart failed: "+err.Error()+". Close the window and start crower-gui again.", u.window)
		return
	}

	u.window.Close()
}
