//go:build fyne

package fyneui

import (
	"errors"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// stubUpdater answers the About tab without touching GitHub.
type stubUpdater struct {
	current     string
	latest      string
	latestErr   error
	previous    string
	hasPrevious bool
	installErr  error
	installed   []string
	restarts    int
	restartErr  error
}

func (u *stubUpdater) Current() string { return u.current }

func (u *stubUpdater) Latest() (string, error) { return u.latest, u.latestErr }

func (u *stubUpdater) Previous() (string, bool) { return u.previous, u.hasPrevious }

func (u *stubUpdater) Install(version string) error {
	if u.installErr != nil {
		return u.installErr
	}

	u.installed = append(u.installed, version)

	return nil
}

func (u *stubUpdater) Restart() error {
	u.restarts++

	return u.restartErr
}

// findLabelText walks the widget tree for a label whose text contains want.
func findLabelText(object fyne.CanvasObject, want string) bool {
	switch typed := object.(type) {
	case *widget.Label:
		return strings.Contains(typed.Text, want)
	case *fyne.Container:
		for _, child := range typed.Objects {
			if findLabelText(child, want) {
				return true
			}
		}
	case *container.Scroll:
		if typed.Content != nil {
			return findLabelText(typed.Content, want)
		}
	case *widget.Card:
		if typed.Content != nil {
			return findLabelText(typed.Content, want)
		}
	}

	return false
}

func TestAboutTabShowsTheVersionAndTheDescription(t *testing.T) {
	updater := &stubUpdater{current: "v9.9.9"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)

	content := front.aboutTab()
	if content == nil {
		t.Fatal("aboutTab returned nil")
	}

	if front.updateStatus == nil || front.checkButton == nil || front.installButton == nil {
		t.Fatal("the update controls were not built")
	}

	if front.installButton.Visible() {
		t.Error("the install button is visible before a check found a release")
	}

	if front.restartButton == nil || front.restartButton.Visible() {
		t.Error("the restart button must exist and start hidden")
	}

	if !findLabelText(content, "v9.9.9") {
		t.Error("the About tab does not show the running version")
	}

	if !findLabelText(content, "undoable history") {
		t.Error("the About tab does not show the project description")
	}
}

// A recorded previous version offers the restore button; a fresh installation
// does not.
func TestAboutOffersRestoreOnlyWhenPreviousIsRecorded(t *testing.T) {
	withPrevious := &stubUpdater{current: "v1.0.0", previous: "v0.9.0", hasPrevious: true}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, withPrevious)
	front.aboutTab()

	if front.restoreButton == nil {
		t.Fatal("no restore button with a recorded previous version")
	}

	if !strings.Contains(front.restoreButton.Text, "v0.9.0") {
		t.Errorf("restore button = %q, want the previous version", front.restoreButton.Text)
	}

	fresh := &stubUpdater{current: "v1.0.0"}
	front, _ = newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, fresh)
	front.aboutTab()

	if front.restoreButton != nil {
		t.Error("a fresh installation offered a restore button")
	}
}

func TestApplyLatestReportsANewerRelease(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	front.applyLatest("v1.1.0", nil)

	if got := front.updateStatus.Text; !strings.Contains(got, "v1.1.0 is available") {
		t.Errorf("status = %q, want the newer release named", got)
	}

	if front.latest != "v1.1.0" {
		t.Errorf("latest = %q, want v1.1.0", front.latest)
	}

	if !front.installButton.Visible() {
		t.Error("the install action did not appear for an available release")
	}

	if !strings.Contains(front.installButton.Text, "v1.1.0") {
		t.Errorf("install button = %q, want the version named", front.installButton.Text)
	}

	if front.restartButton.Visible() {
		t.Error("restart appeared before anything was installed")
	}
}

func TestApplyLatestReportsTheNewestVersion(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	front.applyLatest("v1.0.0", nil)

	if got := front.updateStatus.Text; !strings.Contains(got, "is the newest release") {
		t.Errorf("status = %q, want it to say the running version is newest", got)
	}

	if front.installButton.Visible() {
		t.Error("the install action appeared although the running version is newest")
	}
}

func TestApplyLatestReportsALookupFailure(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	front.applyLatest("", errors.New("network down"))

	if got := front.updateStatus.Text; !strings.Contains(got, "Could not check for updates") {
		t.Errorf("status = %q, want the failure named", got)
	}

	if front.latest != "" {
		t.Errorf("latest = %q, want it cleared after a failure", front.latest)
	}

	if front.installButton.Visible() {
		t.Error("the install action appeared after a failed lookup")
	}
}

// A finished install replaces the install action with the restart, and drops
// the now-stale restore.
func TestAfterInstallOffersTheRestart(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0", previous: "v0.9.0", hasPrevious: true}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	front.applyLatest("v1.1.0", nil)
	front.afterInstall()

	if front.installButton.Visible() {
		t.Error("the install action stayed after the install finished")
	}

	if front.restoreButton.Visible() {
		t.Error("a restore action built from the old previous version stayed visible")
	}

	if !front.restartButton.Visible() {
		t.Error("the restart action did not appear after an install")
	}
}

func TestInstallNowRecordsTheVersion(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	if err := front.installNow("v1.1.0"); err != nil {
		t.Fatalf("installNow: %v", err)
	}

	if len(updater.installed) != 1 || updater.installed[0] != "v1.1.0" {
		t.Fatalf("installed = %v, want [v1.1.0]", updater.installed)
	}

	if got := front.updateStatus.Text; !strings.Contains(got, "Restart to use it") {
		t.Errorf("status = %q, want the restart hint", got)
	}
}

func TestInstallNowReportsAFailure(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0", installErr: errors.New("checksum mismatch")}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	if err := front.installNow("v1.1.0"); err == nil {
		t.Fatal("installNow accepted a failing installer")
	}

	if got := front.updateStatus.Text; !strings.Contains(got, "Install failed") || !strings.Contains(got, "checksum mismatch") {
		t.Errorf("status = %q, want the failure named", got)
	}
}

func TestRestartSpawnsTheNewBinary(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)

	front.restart()

	if updater.restarts != 1 {
		t.Errorf("restarts = %d, want 1", updater.restarts)
	}
}

// A narrow window must wrap the update actions instead of pushing them past the
// edge: the old single HBox overflowed as soon as Check, Install and Restore
// were all visible.
func TestAboutActionsFitANarrowWindow(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0", previous: "v0.9.0", hasPrevious: true}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)

	// Reveal Install before the first real layout: the software canvas never
	// lays out a widget that was shown after the window settled.
	content := front.aboutTab()
	front.applyLatest("v1.1.0", nil)

	window.SetContent(content)
	window.Resize(fyne.NewSize(360, 640))

	visible := 0
	rows := make(map[float32]bool)

	for _, button := range []*widget.Button{front.checkButton, front.installButton, front.restartButton, front.restoreButton} {
		if button == nil || !button.Visible() {
			continue
		}

		visible++

		if button.Size().Width == 0 {
			t.Errorf("%q was never laid out", button.Text)
		}

		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(button)
		if right := pos.X + button.Size().Width; right > 360 {
			t.Errorf("%q right edge = %.1f, want it inside the 360 px window", button.Text, right)
		}

		rows[pos.Y] = true
	}

	if visible < 3 {
		t.Fatalf("only %d actions visible, want Check, Install and Restore", visible)
	}

	if len(rows) < 2 {
		t.Errorf("actions share one row (%v), want them to wrap at 360 px", rows)
	}
}

// Hiding an action re-runs the flow and the body above it: the flow reports the
// height of its last layout, so without both the card would keep a hidden
// button's row until the window was resized.
func TestHidingAnAboutActionReclaimsItsRow(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0", previous: "v0.9.0", hasPrevious: true}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)

	content := front.aboutTab()
	front.applyLatest("v1.1.0", nil)

	window.SetContent(content)
	window.Resize(fyne.NewSize(360, 640))

	before := front.aboutUpdates.Size().Height

	front.showInstall("")

	after := front.aboutUpdates.Size().Height
	if after >= before {
		t.Errorf("card height = %.1f after hiding Install, want less than %.1f", after, before)
	}

	if front.installButton.Visible() {
		t.Error("install is still visible")
	}
}
