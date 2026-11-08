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

// cardTextSize is the size of a card's title line, shared by the expandable
// cards and the simple ones. Emoji and summary are one text run at one size, so
// they share a baseline: two objects of different sizes never line up, whatever
// their boxes say.
const cardTextSize float32 = 18

// cardHeader is an expandable card's title row: the emoji and the summary as one
// line of text, and the chevron.
//
// Emoji and summary share a text run, which is what a dropdown does and why the
// emoji reads correctly there. Drawing them as two objects of different sizes
// leaves the emoji's baseline below the text's no matter how the boxes are
// centred (#480).
type cardHeader struct {
	widget.BaseWidget

	// Text is the summary without the emoji. Tests and the cards read it.
	Text string

	OnTapped func()

	emoji   string
	text    *canvas.Text
	chevron *widget.Icon
	open    bool
	box     *fyne.Container
}

func newCardHeader(emoji, title string) *cardHeader {
	header := &cardHeader{Text: title, emoji: emoji}

	// canvas.Text, not a Label: a Label is a Tappable in the canvas and would
	// swallow the click that opens the card.
	header.text = canvas.NewText(header.line(), theme.Color(theme.ColorNameForeground))
	header.text.TextSize = cardTextSize
	header.text.TextStyle = fyne.TextStyle{Bold: true}

	header.chevron = widget.NewIcon(theme.MenuDropDownIcon())

	// Vertical air only, so the title clears the card's top and bottom and the
	// tap target stays comfortable. No horizontal inset: the card's own padding
	// puts the title off the edge, and the body is added at this same left edge
	// so the two align (#516).
	header.box = container.New(
		layout.NewCustomPaddedLayout(appPadding, appPadding, 0, 0),
		container.NewBorder(
			nil,
			nil,
			nil,
			container.NewCenter(header.chevron),
			centeredVertically(header.text),
		),
	)

	header.ExtendBaseWidget(header)

	return header
}

// line is what the header draws: the emoji, a gap, and the summary. One string
// means one baseline.
func (h *cardHeader) line() string {
	if h.emoji == "" {
		return h.Text
	}

	return h.emoji + "  " + h.Text
}

// CreateRenderer hands the row to the canvas.
func (h *cardHeader) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(h.box)
}

// Tapped opens or closes the card. The label and the chevron are not tappable,
// so a click anywhere in the row arrives here.
func (h *cardHeader) Tapped(*fyne.PointEvent) {
	if h.OnTapped != nil {
		h.OnTapped()
	}
}

// setOpen points the chevron the way the card is going.
func (h *cardHeader) setOpen(open bool) {
	h.open = open

	if open {
		h.chevron.SetResource(theme.MenuDropUpIcon())
		return
	}

	h.chevron.SetResource(theme.MenuDropDownIcon())
}

// SetText replaces the summary, keeping the emoji.
func (h *cardHeader) SetText(title string) {
	h.Text = title
	h.text.Text = h.line()
	h.text.Refresh()
}
