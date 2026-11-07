//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// cardEmojiSize is how large a command's or an operation's icon is drawn in a
// card header. The button it replaced drew it at the body text size, where an
// emoji is a smudge.
const cardEmojiSize float32 = 22

// cardHeader is an expandable card's title row: the emoji at a readable size,
// the summary text, and the chevron.
//
// It replaces the header button. A Button draws its whole label at one text
// size, so growing the emoji would have grown the summary with it; here the two
// have their own sizes.
type cardHeader struct {
	widget.BaseWidget

	// Text is the summary without the emoji. Tests and the cards read it.
	Text string

	OnTapped func()

	emoji   *canvas.Text
	title   *canvas.Text
	chevron *widget.Icon
	open    bool
	box     *fyne.Container
}

func newCardHeader(emoji, title string) *cardHeader {
	header := &cardHeader{Text: title}

	header.emoji = canvas.NewText(emoji, theme.Color(theme.ColorNameForeground))
	header.emoji.TextSize = cardEmojiSize

	if emoji == "" {
		header.emoji.Hide()
	}

	// The title is canvas.Text, not a Label: a Label is a Tappable in the
	// canvas, so it would swallow the click that opens the card.
	header.title = canvas.NewText(title, theme.Color(theme.ColorNameForeground))
	header.title.TextSize = theme.Size(theme.SizeNameText)

	header.chevron = widget.NewIcon(theme.MenuDropDownIcon())

	header.box = container.NewBorder(
		nil,
		nil,
		nil,
		container.NewCenter(header.chevron),
		container.NewHBox(header.emoji, header.title),
	)

	header.ExtendBaseWidget(header)

	return header
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
	h.title.Text = title
	h.title.Refresh()
}
