//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// expandableCard is a header button with a body that opens below it.
//
// The body is built once, on the first open, so a collapsed card carries only a
// header. The owner assigns header.OnTapped, which is how a tab keeps one card
// open at a time.
type expandableCard struct {
	widget.BaseWidget

	header *cardHeader
	body   *fyne.Container
	box    *fyne.Container
	build  func() fyne.CanvasObject
	built  bool
	open   bool
}

func newExpandableCard(emoji, title string, build func() fyne.CanvasObject) *expandableCard {
	card := &expandableCard{build: build}

	card.header = newCardHeader(emoji, title)

	card.body = container.NewVBox()
	card.body.Hide()

	background := canvas.NewRectangle(theme.Color(theme.ColorNameInputBackground))
	background.CornerRadius = theme.Size(theme.SizeNameInputRadius)
	background.StrokeColor = theme.Color(theme.ColorNameSeparator)
	background.StrokeWidth = theme.Size(theme.SizeNameSeparatorThickness)

	card.box = container.NewStack(
		background,
		container.NewPadded(container.NewVBox(card.header, card.body)),
	)
	card.ExtendBaseWidget(card)

	return card
}

// CreateRenderer hands the header and body container to the canvas.
func (c *expandableCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.box)
}

// setOpen opens or closes the card, building the body on the first open.
func (c *expandableCard) setOpen(open bool) {
	if open && !c.built {
		// The body is inset from the header text and separated from it by a
		// little air, so an expanded card reads as a panel under its title
		// rather than a continuation of it.
		c.body.Add(container.NewVBox(gap(0, theme.Padding()), indented(c.build())))
		c.built = true
	}

	c.open = open

	if open {
		c.body.Show()
	} else {
		c.body.Hide()
	}

	c.header.setOpen(open)
	c.box.Refresh()
}

// isOpen reports whether the card is showing its body.
func (c *expandableCard) isOpen() bool {
	return c.open
}
