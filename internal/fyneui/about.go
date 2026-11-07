//go:build fyne

package fyneui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
	version.Selectable = true

	summary := widget.NewLabel(projectSummary)
	summary.Alignment = fyne.TextAlignCenter
	summary.Wrapping = fyne.TextWrapWord
	summary.Selectable = true

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
	u.updateStatus.Alignment = fyne.TextAlignCenter
	u.updateStatus.Selectable = true

	u.checkButton = newButton("Check for updates", theme.ViewRefreshIcon(), widget.MediumImportance, u.checkForUpdates)

	// Install and restart are state, not furniture: install appears when the
	// check found a release, restart appears once a release was put on disk.
	u.installButton = newButton("Install", theme.DownloadIcon(), widget.HighImportance, u.installLatest)
	u.installButton.Hide()

	u.restartButton = newButton("Restart now", theme.ViewRefreshIcon(), widget.SuccessImportance, u.restart)
	u.restartButton.Hide()

	// The update actions are vertical: one centred button per row follows the
	// window width and reads as a list of steps. It is kept on ui because hiding
	// a button alone does not re-run this layout.
	actions := container.NewVBox(
		container.NewCenter(u.checkButton),
		container.NewCenter(u.installButton),
		container.NewCenter(u.restartButton),
	)
	u.aboutActions = actions

	if previous, found := u.updater.Previous(); found {
		u.restoreButton = newButton("Restore "+previous, nil, widget.WarningImportance, u.restorePrevious)
		actions.Add(container.NewCenter(u.restoreButton))
	}

	updates := newCenteredCard(
		"Updates",
		"Checks GitHub for the newest release. Nothing is replaced without confirmation.",
		container.NewVBox(u.updateStatus, actions),
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

	u.aboutBody = body
	u.aboutUpdates = updates

	// The body fills the width so the summary and the card wrap with the window;
	// the page adds the margins and clears the scrollbar. Centring the block
	// pinned it to its minimum width and defeated both.
	return container.NewVScroll(page(body))
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
		u.showInstall("")
		return
	}

	u.latest = latest

	switch {
	case latest == "":
		u.updateStatus.SetText("No release information was returned.")
		u.showInstall("")
	case latest == u.updater.Current():
		u.updateStatus.SetText("crower " + u.updater.Current() + " is the newest release.")
		u.showInstall("")
	default:
		u.updateStatus.SetText("crower " + latest + " is available (you have " + u.updater.Current() + ").")
		u.showInstall(latest)
	}
}

// showInstall offers the install action only when there is a release to install,
// so the button appears because the check found one instead of sitting there
// permanently.
func (u *ui) showInstall(version string) {
	if u.installButton == nil {
		return
	}

	if version == "" {
		u.installButton.Hide()
		u.refreshAboutActions()
		return
	}

	u.installButton.SetText("Install " + version)
	u.installButton.Show()
	u.refreshAboutActions()
}

// refreshAboutActions re-runs the action flow and then the body above it.
//
// Hiding a widget only repaints it, and the flow layout reports the minimum
// height of its last layout, so the actions container alone would keep the old
// height. Refreshing the body re-reads the card's minimum and lets it shrink.
func (u *ui) refreshAboutActions() {
	if u.aboutActions != nil {
		u.aboutActions.Refresh()
	}

	if u.aboutBody != nil {
		u.aboutBody.Refresh()
	}
}

// afterInstall is what a successful install or restore leaves: no install
// action (there is nothing newer to fetch), no stale restore, and a restart,
// because the binary on disk is now a different one.
func (u *ui) afterInstall() {
	u.showInstall("")

	if u.restoreButton != nil {
		u.restoreButton.Hide()
	}

	if u.restartButton != nil {
		u.restartButton.Show()
	}

	u.refreshAboutActions()
}

func (u *ui) setBusy(status string) {
	u.updateStatus.SetText(status)
	u.checkButton.Disable()
	u.installButton.Disable()

	if u.restoreButton != nil {
		u.restoreButton.Disable()
	}

	if u.restartButton != nil {
		u.restartButton.Disable()
	}
}

func (u *ui) setIdle() {
	u.checkButton.Enable()
	u.installButton.Enable()

	if u.restoreButton != nil {
		u.restoreButton.Enable()
	}

	if u.restartButton != nil {
		u.restartButton.Enable()
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

	showConfirm("Install "+version, question, u.window, func(ok bool) {
		if !ok {
			return
		}

		u.installInBackground(version)
	})
}

// restorePrevious installs the version the installation ran before this one.
func (u *ui) restorePrevious() {
	previous, found := u.updater.Previous()
	if !found {
		u.updateStatus.SetText("No previous version is recorded yet; install one with crower --install-version.")
		return
	}

	question := "Install crower " + previous + " again?\n\nIt replaces the binary on disk; a restart finishes the change."

	showConfirm("Restore "+previous, question, u.window, func(ok bool) {
		if !ok {
			return
		}

		u.installInBackground(previous)
	})
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

			u.afterInstall()
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

// restart closes the window; Run performs the restart once the window loop has
// returned, so the new process starts with the old one shut down.
func (u *ui) restart() {
	if u.updater == nil {
		return
	}

	u.restartRequested = true
	u.window.Close()
}
