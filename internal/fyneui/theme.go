//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	// appTextSize is the body text size. The toolkit's default of 14 leaves an
	// emoji in a dropdown and a command's summary at the edge of legibility on a
	// 1080p screen; every label, button and select grows with it.
	appTextSize float32 = 16

	// appPadding is the window's base spacing unit. It is installed as the
	// theme's padding, so every box, border and form layout that Fyne builds
	// separates its children by exactly this much: buttons in a row, form rows,
	// and the gap around a tab.
	//
	// The toolkit's own 4 px is tight enough that two controls look joined; 8 px
	// reads as two controls that belong together. A hand-inserted gap() spacer
	// must never be added on top of it - a spacer is a normal child to a box
	// layout and brings two more paddings with it (#516).
	appPadding float32 = 8

	// appSection separates groups rather than controls: a card's title from its
	// body, a page from the window's edge, one toolbar group from the next. Two
	// base units, so the break is visible without a rule or a heading.
	appSection float32 = appPadding * 2
)

// appTheme is the default theme with a larger body text and one spacing scale.
type appTheme struct {
	fyne.Theme
}

func newAppTheme() fyne.Theme {
	return appTheme{Theme: theme.DefaultTheme()}
}

func (t appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return appTextSize
	case theme.SizeNamePadding:
		return appPadding
	}

	return t.Theme.Size(name)
}
