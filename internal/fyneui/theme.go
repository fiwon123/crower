//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// appTextSize is the body text size. The toolkit's default of 14 leaves an
// emoji in a dropdown and a command's summary at the edge of legibility on a
// 1080p screen; every label, button and select grows with it.
const appTextSize float32 = 16

// appTheme is the default theme with a larger body text.
type appTheme struct {
	fyne.Theme
}

func newAppTheme() fyne.Theme {
	return appTheme{Theme: theme.DefaultTheme()}
}

func (t appTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return appTextSize
	}

	return t.Theme.Size(name)
}
