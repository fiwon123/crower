//go:build fyne

package fyneui

import (
	"fmt"

	"fyne.io/fyne/v2"
)

// handleDrop fills the copy/move sources from a drop anywhere on the window.
//
// Fyne v2.8 offers a window-wide drop callback only, not a per-widget one, so
// the rule is simple and promises nothing it cannot know: the paths are
// appended to the Sources field and the Tools tab opens, where the user picks a
// destination. A drop never runs anything.
func (u *ui) handleDrop(_ fyne.Position, uris []fyne.URI) {
	if u.sourcesEntry == nil || len(uris) == 0 {
		return
	}

	count := 0

	for _, uri := range uris {
		path := uriPath(uri)
		if path == "" {
			continue
		}

		appendPath(u.sourcesEntry, path)
		count++
	}

	if count == 0 {
		return
	}

	u.selectTools()

	if u.transferResult != nil {
		u.transferResult.report(fmt.Sprintf("Added %d path(s) from the drop. Pick a destination, then Copy or Move.", count))
	}
}

// selectTools opens the Tools tab, so a drop shows where the paths landed.
func (u *ui) selectTools() {
	if u.tabs == nil {
		return
	}

	for _, item := range u.tabs.Items {
		if item.Text == "Tools" {
			u.tabs.Select(item)
			return
		}
	}
}
