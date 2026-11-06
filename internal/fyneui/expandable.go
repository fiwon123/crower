//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
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

	header *widget.Button
	body   *fyne.Container
	box    *fyne.Container
	build  func() fyne.CanvasObject
	built  bool
	open   bool
}

func newExpandableCard(title string, build func() fyne.CanvasObject) *expandableCard {
	card := &expandableCard{build: build}

	card.header = widget.NewButton(title, nil)
	card.header.Alignment = widget.ButtonAlignLeading
	card.header.IconPlacement = widget.ButtonIconTrailingText
	card.header.Importance = widget.LowImportance

	card.body = container.NewVBox()
	card.body.Hide()

	card.box = container.NewVBox(card.header, card.body)
	card.ExtendBaseWidget(card)
	card.updateIcon()

	return card
}

// CreateRenderer hands the header and body container to the canvas.
func (c *expandableCard) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(c.box)
}

// setOpen opens or closes the card, building the body on the first open.
func (c *expandableCard) setOpen(open bool) {
	if open && !c.built {
		c.body.Add(c.build())
		c.built = true
	}

	c.open = open

	if open {
		c.body.Show()
	} else {
		c.body.Hide()
	}

	c.updateIcon()
	c.box.Refresh()
}

// isOpen reports whether the card is showing its body.
func (c *expandableCard) isOpen() bool {
	return c.open
}

func (c *expandableCard) updateIcon() {
	if c.open {
		c.header.SetIcon(theme.MenuDropUpIcon())
		return
	}

	c.header.SetIcon(theme.MenuDropDownIcon())
}
