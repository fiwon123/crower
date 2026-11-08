//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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

	// The title and the body are two groups inside one card, so they are
	// separated by the section gap rather than the base unit. The body starts on
	// the title's own left edge: its content is a continuation of the heading,
	// not a nested panel (#516).
	card.box = container.NewStack(
		background,
		container.NewPadded(container.New(
			layout.NewCustomPaddedVBoxLayout(appSection),
			card.header,
			card.body,
		)),
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
		// The body is a continuation of the title, so it is added at the same
		// left edge. The section gap above it comes from the card's own layout.
		c.body.Add(c.build())
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
