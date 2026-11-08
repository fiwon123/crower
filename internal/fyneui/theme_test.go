//go:build fyne

package fyneui

import (
	"testing"

	"fyne.io/fyne/v2/theme"
)

// The window has one spacing scale, and it is installed as the theme's padding:
// every Fyne box, border and form layout then separates its children by the same
// unit (#516).
func TestAppThemeInstallsOneSpacingScale(t *testing.T) {
	useAppTheme(t)

	if got := theme.Padding(); got != appPadding {
		t.Errorf("theme padding = %.1f, want the base unit %.1f", got, appPadding)
	}

	if got := theme.Size(theme.SizeNameText); got != appTextSize {
		t.Errorf("text size = %.1f, want %.1f", got, appTextSize)
	}
}

// The action gap is the base unit, not half of it: a button row uses the same
// space as every other pair of related controls (#516).
func TestActionGapIsTheBaseUnit(t *testing.T) {
	if got := actionGap(); got != appPadding {
		t.Errorf("action gap = %.1f, want the base unit %.1f", got, appPadding)
	}
}
