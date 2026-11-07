//go:build fyne

package fyneui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/domain/catalog"
)

// commandCard is one command: a name header that opens into its fields,
// parameters and actions.
type commandCard struct {
	*expandableCard

	key    string
	params *widget.Entry
	run    *widget.Button
	edit   *widget.Button
	remove *widget.Button
}

func (u *ui) newCommandCard(command catalog.Command) *commandCard {
	card := &commandCard{key: command.Name}
	card.expandableCard = newExpandableCard(command.Icon, shortLabel(command.Name, 40), func() fyne.CanvasObject {
		return u.buildCommandBody(card)
	})
	card.header.OnTapped = func() { u.toggleCommandCard(card) }

	return card
}

// buildCommandBody builds the fields, the parameter entry and the actions on
// first open.
func (u *ui) buildCommandBody(card *commandCard) fyne.CanvasObject {
	command, found := u.model.Commands.Get(card.key)
	if !found {
		return widget.NewLabel("This command is no longer in the current setup.")
	}

	form := newFieldList(commandDetailFields(command))

	card.params = newScrollFreeEntry()
	card.params.SetPlaceHolder(paramsHint(command))
	card.params.OnSubmitted = func(string) { u.runCommandCard(card) }

	card.run = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() { u.runCommandCard(card) })
	card.run.Importance = widget.SuccessImportance

	card.edit = widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() { u.showEditFor(card.key) })
	card.edit.Importance = widget.HighImportance

	card.remove = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() { u.showDeleteFor(card.key) })
	card.remove.Importance = widget.DangerImportance

	// Run sits with Edit and Delete, not inside the parameter row: a row of
	// actions reads as one group, and the field gets the full card width.
	actions := spacedActions(card.run, card.edit, card.remove)

	return container.NewVBox(
		form,
		widget.NewSeparator(),
		container.NewVBox(selectableLabel("Parameters"), card.params),
		gap(0, theme.Padding()),
		actions,
		gap(0, theme.Padding()),
	)
}

// toggleCommandCard opens a card and closes the one that was open. Opening also
// makes the command the model's selection, which the Edit/Run paths expect.
func (u *ui) toggleCommandCard(card *commandCard) {
	if card.isOpen() {
		card.setOpen(false)
		u.commandOpen = nil
		u.commandPage.Refresh()

		if u.commandScroll != nil {
			u.commandScroll.Refresh()
		}

		return
	}

	if u.commandOpen != nil && u.commandOpen != card {
		u.commandOpen.setOpen(false)
	}

	card.setOpen(true)
	u.commandOpen = card

	for index, candidate := range u.commandCards {
		if candidate == card {
			u.model.Index = index
			break
		}
	}

	u.commandPage.Refresh()
	u.scrollCommandTo(card)
}

// scrollCommandTo brings an opened card just into view, moving the minimum so
// the page does not jump to the top.
func (u *ui) scrollCommandTo(card *commandCard) {
	if u.commandScroll == nil || u.commandPage == nil || u.commandList == nil {
		return
	}

	y := u.commandPage.Position().Y + u.commandList.Position().Y + card.Position().Y
	offset := revealOffset(y, card.Size().Height, u.commandScroll.Size().Height, u.commandScroll.Offset.Y)
	if offset == u.commandScroll.Offset.Y {
		return
	}

	u.commandScroll.ScrollToOffset(fyne.NewPos(0, offset))
}

// runCommandCard runs the command with the card's parameter entry as its
// arguments. An empty entry runs it without parameters.
func (u *ui) runCommandCard(card *commandCard) {
	if card.params == nil {
		return
	}

	u.model.StartRunKey(card.key, true)
	u.model.ParamInput = card.params.Text
	u.model.SubmitParams()
	u.showOutput()
}
