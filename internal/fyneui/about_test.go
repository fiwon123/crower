//go:build fyne

package fyneui

import (
	"errors"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// stubUpdater answers the About tab without touching GitHub.
type stubUpdater struct {
	current        string
	latest         string
	latestErr      error
	previous       string
	hasPrevious    bool
	installErr     error
	installed      []string
	installStarted chan string
	restarts       int
	restartErr     error
}

func (u *stubUpdater) Current() string { return u.current }

func (u *stubUpdater) Latest() (string, error) { return u.latest, u.latestErr }

func (u *stubUpdater) Previous() (string, bool) { return u.previous, u.hasPrevious }

func (u *stubUpdater) Install(version string) error {
	if u.installErr != nil {
		return u.installErr
	}

	u.installed = append(u.installed, version)

	if u.installStarted != nil {
		u.installStarted <- version
	}

	return nil
}

func (u *stubUpdater) Restart() error {
	u.restarts++

	return u.restartErr
}

// findLabel returns the label with exactly this text, or nil.
func findLabel(object fyne.CanvasObject, want string) *widget.Label {
	switch typed := object.(type) {
	case *widget.Label:
		if typed.Text == want {
			return typed
		}
	case *fyne.Container:
		for _, child := range typed.Objects {
			if label := findLabel(child, want); label != nil {
				return label
			}
		}
	case *container.Scroll:
		if typed.Content != nil {
			return findLabel(typed.Content, want)
		}
	}

	return nil
}

// findText walks the widget tree for a canvas text with exactly this content.
func findText(object fyne.CanvasObject, want string) *canvas.Text {
	switch typed := object.(type) {
	case *canvas.Text:
		if typed.Text == want {
			return typed
		}
	case *fyne.Container:
		for _, child := range typed.Objects {
			if found := findText(child, want); found != nil {
				return found
			}
		}
	case *container.Scroll:
		if typed.Content != nil {
			return findText(typed.Content, want)
		}
	case fyne.Widget:
		renderer := test.WidgetRenderer(typed)
		if renderer == nil {
			return nil
		}

		for _, child := range renderer.Objects() {
			if found := findText(child, want); found != nil {
				return found
			}
		}
	}

	return nil
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

// Restoring a previous version is a deliberate action: Restore stays disabled
// until exactly yes is typed, cancelling installs nothing, and confirming
// installs the recorded previous version (#519).
func TestRestorePreviousNeedsATypedYes(t *testing.T) {
	updater := &stubUpdater{
		current:        "v1.0.0",
		previous:       "v0.9.0",
		hasPrevious:    true,
		installStarted: make(chan string, 1),
	}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)
	front.aboutTab()

	openRestore := func() fyne.CanvasObject {
		t.Helper()

		front.restorePrevious()

		overlays := window.Canvas().Overlays().List()
		if len(overlays) == 0 {
			t.Fatal("restore did not ask for confirmation")
		}

		return overlays[len(overlays)-1]
	}

	pop := openRestore()

	if text := strings.Join(labelsIn(pop), "\n"); !strings.Contains(text, "v0.9.0") {
		t.Errorf("dialog text = %q, want the version it restores", text)
	}

	action := findButton(pop, "Restore")
	if action == nil {
		t.Fatal("the dialog has no Restore button")
	}

	if !action.Disabled() {
		t.Error("Restore starts enabled")
	}

	entry := entryIn(pop)
	if entry == nil {
		t.Fatal("the dialog has no confirmation field")
	}

	for _, typed := range []string{"y", "YES", "yes please"} {
		entry.SetText(typed)

		if !action.Disabled() {
			t.Errorf("%q enabled Restore", typed)
		}
	}

	entry.SetText("yes")
	if action.Disabled() {
		t.Fatal("exactly yes did not enable Restore")
	}

	// Cancel is a no: it must not replace the binary.
	cancel := findButton(pop, "Cancel")
	if cancel == nil {
		t.Fatal("the dialog has no Cancel button")
	}

	cancel.OnTapped()

	if len(updater.installed) != 0 {
		t.Fatalf("cancelling installed %v", updater.installed)
	}

	// Confirming installs the previous version, not the newest release.
	pop = openRestore()

	entry = entryIn(pop)
	entry.SetText("yes")

	action = findButton(pop, "Restore")
	action.OnTapped()

	select {
	case got := <-updater.installStarted:
		if got != "v0.9.0" {
			t.Errorf("installed %q, want the previous version v0.9.0", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("confirming Restore never started the install")
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

// A narrow window must keep every update action inside the width: the old
// single HBox pushed the last one past the edge.
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

	for _, button := range []*appButton{front.checkButton, front.installButton, front.restartButton, front.restoreButton} {
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

// The update steps are one centred column with the base-unit gap (#500). The
// column must sit in the middle of the card, not merely agree with itself: a
// stack that was never stretched keeps every button on a common centre while
// hugging the card's left edge (#503).
func TestAboutActionsAreCentredAndSpaced(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0", previous: "v0.9.0", hasPrevious: true}
	front, window := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)

	content := front.aboutTab()
	front.applyLatest("v1.1.0", nil)

	window.SetContent(content)
	window.Resize(fyne.NewSize(960, 720))

	driver := fyne.CurrentApp().Driver()

	// The card insets its content symmetrically, so a centred column has the
	// card's centre for its own.
	cardPos := driver.AbsolutePositionForObject(front.aboutUpdates)
	cardCentre := cardPos.X + front.aboutUpdates.Size().Width/2

	var centre, lastBottom float32

	for index, button := range []*appButton{front.checkButton, front.installButton, front.restoreButton} {
		if button == nil || !button.Visible() {
			t.Fatalf("action %d is not visible", index)
		}

		// Measure the embedded button: it is what the canvas paints, so a
		// wrapper placed correctly over a drifting button is still a failure.
		pos := driver.AbsolutePositionForObject(button.Button)

		got := pos.X + button.Size().Width/2
		if index == 0 {
			centre = got
		} else if diff := got - centre; diff < -0.5 || diff > 0.5 {
			t.Errorf("%q centre = %.1f, want %.1f", button.Text, got, centre)
		}

		if diff := got - cardCentre; diff < -1 || diff > 1 {
			t.Errorf("%q centre = %.1f, want the card's centre %.1f", button.Text, got, cardCentre)
		}

		if index > 0 {
			if gap := pos.Y - lastBottom; gap != actionGap() {
				t.Errorf("gap above %q = %.1f, want %.1f", button.Text, gap, actionGap())
			}
		}

		lastBottom = pos.Y + button.Size().Height
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

// The About card's title is centred like the rest of the page, and drawn with
// the same one card-title treatment as every other card (#516).
func TestAboutUpdateTitleIsCentred(t *testing.T) {
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v0.0.0-test"})

	title := findText(front.aboutTab(), "Updates")
	if title == nil {
		t.Fatal("the Updates title is not a card title")
	}

	if title.Alignment != fyne.TextAlignCenter {
		t.Errorf("title alignment = %v, want centre", title.Alignment)
	}

	if title.TextSize != cardTextSize || !title.TextStyle.Bold {
		t.Errorf("title is %.0fpt bold=%v, want the card title %.0fpt bold", title.TextSize, title.TextStyle.Bold, cardTextSize)
	}
}

// The version and status can be selected and copied for a bug report.
func TestAboutTextIsSelectable(t *testing.T) {
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, &stubUpdater{current: "v1.0.0"})

	content := front.aboutTab()

	if !front.updateStatus.Selectable {
		t.Error("the update status is not selectable")
	}

	version := findLabel(content, "Version v1.0.0")
	if version == nil || !version.Selectable {
		t.Error("the version label is missing or not selectable")
	}
}

// Restart closes the window; Run spawns/execs after the window loop returns.
func TestRestartDefersTheSpawnUntilTheWindowCloses(t *testing.T) {
	updater := &stubUpdater{current: "v1.0.0"}
	front, _ := newTestUIWithUpdater(t, testModel(), stubEditor{}, stubReverter{}, updater)

	front.restart()

	if !front.restartRequested {
		t.Error("restart did not record the request")
	}

	if updater.restarts != 0 {
		t.Error("restart spawned before the window closed")
	}
}
