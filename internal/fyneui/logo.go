//go:build fyne

package fyneui

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

// logo.png is rendered from assets/logo.svg; the SVG is the source, the PNG is
// what the window manager accepts as an app icon.
//
//go:embed logo.png
var logoPNG []byte

func logoResource() fyne.Resource {
	return fyne.NewStaticResource("crower-logo.png", logoPNG)
}
