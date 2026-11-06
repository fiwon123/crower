//go:build fyne

package fyneui

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

// historyHeader is the sentence above the list.
const historyHeader = "Every recorded operation, newest first. Open an entry to see what it recorded and to restore that state."

// historyPageSize is how many entries a page renders. The log stays in memory,
// but only this many cards are rendered per page, which is what keeps a long
// history responsive.
const historyPageSize = 100

// historyCard is one entry: a summary header that opens into its detail.
type historyCard struct {
	*expandableCard

	version int
	detail  desktop.HistoryDetail

	raw       *widget.Label
	rawToggle *widget.Button
	rawShown  bool
	run       *widget.Button
	edit      *widget.Button
	restore   *widget.Button
	hint      *widget.Label
}

// historyTab renders the log as a scroll of expandable entry cards.
//
// It replaced a table with a global detail panel: selecting a row used to
// refresh every cell and re-lay out the tab twice, and every note in the log was
// parsed on each refresh. Now only the opened card is built, and only the
// rendered page is formatted.
func (u *ui) historyTab() fyne.CanvasObject {
	header := widget.NewLabel(historyHeader)
	header.Wrapping = fyne.TextWrapWord

	u.historyEmpty = widget.NewLabel("No history yet.")
	u.historyEmpty.Hide()

	u.historyList = container.NewVBox()

	u.historyLoadMore = widget.NewButton("Load older", func() {
		u.loadOlderHistory()
	})
	u.historyLoadMore.Importance = widget.LowImportance
	u.historyLoadMore.Hide()

	u.historyPage = container.NewVBox(header, u.historyEmpty, u.historyList, u.historyLoadMore)
	u.historyScroll = container.NewVScroll(page(u.historyPage))

	return u.historyScroll
}

// refreshHistory rebuilds the first page from the current log. It formats at
// most historyPageSize entries and parses no notes.
func (u *ui) refreshHistory() {
	if u.historyList == nil {
		return
	}

	u.historyTotal = len(u.model.History.AllData)
	u.historyLimit = min(historyPageSize, u.historyTotal)
	u.historyOpen = nil
	u.historyCards = nil

	u.historyList.RemoveAll()
	u.addHistoryCards(u.model.HistoryPage(0, u.historyLimit))

	if u.historyTotal == 0 {
		u.historyEmpty.Show()
	} else {
		u.historyEmpty.Hide()
	}

	u.updateHistoryLoad()
	u.historyPage.Refresh()

	// A restore drops entries, so the page can be shorter than the offset the
	// scroll was holding; go back to the top rather than show blank space.
	if u.historyScroll != nil {
		u.historyScroll.ScrollToOffset(fyne.NewPos(0, 0))
	}
}

// loadOlderHistory appends the next page without rebuilding the cards already
// on screen.
func (u *ui) loadOlderHistory() {
	offset := u.historyLimit

	summaries := u.model.HistoryPage(offset, historyPageSize)
	if len(summaries) == 0 {
		return
	}

	u.addHistoryCards(summaries)
	u.historyLimit += len(summaries)
	u.updateHistoryLoad()
	u.historyPage.Refresh()
}

// addHistoryCards builds the summaries' cards in one batch: Add lays the VBox
// out on every call, and a page is up to a hundred of them.
func (u *ui) addHistoryCards(summaries []desktop.HistorySummary) {
	if len(summaries) == 0 {
		return
	}

	objects := make([]fyne.CanvasObject, 0, len(summaries))

	for _, summary := range summaries {
		card := u.newHistoryCard(summary)
		u.historyCards = append(u.historyCards, card)
		objects = append(objects, card.expandableCard)
	}

	u.historyList.Objects = append(u.historyList.Objects, objects...)
	u.historyList.Refresh()
}

func (u *ui) updateHistoryLoad() {
	if u.historyLimit < u.historyTotal {
		u.historyLoadMore.Show()
		return
	}

	u.historyLoadMore.Hide()
}

func (u *ui) newHistoryCard(summary desktop.HistorySummary) *historyCard {
	card := &historyCard{version: summary.Version}
	card.expandableCard = newExpandableCard(historySummaryTitle(summary), func() fyne.CanvasObject {
		return u.buildHistoryBody(card)
	})
	card.header.OnTapped = func() { u.toggleHistoryCard(card) }

	return card
}

// historySummaryTitle is the collapsed row: when, what, which command, version.
// A long command is shortened so the header button's text cannot widen the
// window's minimum size.
func historySummaryTitle(summary desktop.HistorySummary) string {
	command := summary.Command
	if command == "" {
		command = "-"
	}

	return fmt.Sprintf("%s · %s · %s · #%d", summary.When, summary.Operation, shortLabel(command, 40), summary.Version)
}

// shortLabel shortens text to limit runes, keeping an ellipsis.
func shortLabel(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}

	return string(runes[:limit-1]) + "…"
}

// toggleHistoryCard opens a card and closes the one that was open, so the page
// never holds two detail bodies.
func (u *ui) toggleHistoryCard(card *historyCard) {
	if card.isOpen() {
		card.setOpen(false)
		u.historyOpen = nil
		u.historyPage.Refresh()

		if u.historyScroll != nil {
			u.historyScroll.Refresh()
		}

		return
	}

	if u.historyOpen != nil && u.historyOpen != card {
		u.historyOpen.setOpen(false)
	}

	card.setOpen(true)
	u.historyOpen = card
	u.historyPage.Refresh()
	u.scrollHistoryTo(card)
}

// scrollHistoryTo brings an opened card into view. The position is relative to
// the page, which is what the scroll wraps.
func (u *ui) scrollHistoryTo(card *historyCard) {
	if u.historyScroll == nil || u.historyPage == nil || u.historyList == nil {
		return
	}

	y := u.historyPage.Position().Y + u.historyList.Position().Y + card.Position().Y
	u.historyScroll.ScrollToOffset(fyne.NewPos(0, y))
}

// buildHistoryBody builds the detail for one card on first open.
func (u *ui) buildHistoryBody(card *historyCard) fyne.CanvasObject {
	detail, found := u.model.HistoryDetail(card.version)
	if !found {
		return widget.NewLabel("This entry is no longer in the history.")
	}

	card.detail = detail

	fields := make([]detailField, 0, len(detail.Fields))
	for _, field := range detail.Fields {
		fields = append(fields, detailField{
			Label: historyFieldLabel(field.Key),
			Value: historyFieldValue(field.Key, field.Value),
		})
	}

	if len(fields) == 0 {
		fields = append(fields, detailField{Value: "(no note)"})
	}

	grid := container.NewGridWithColumns(2)
	setFieldGrid(grid, fields)

	card.raw = widget.NewLabel(detail.Note)
	card.raw.Wrapping = fyne.TextWrapWord
	card.raw.Hide()

	card.rawToggle = widget.NewButton("Raw note", func() {
		card.rawShown = !card.rawShown
		if card.rawShown {
			card.raw.Show()
		} else {
			card.raw.Hide()
		}

		u.historyPage.Refresh()

		if u.historyScroll != nil {
			u.historyScroll.Refresh()
		}
	})
	card.rawToggle.Importance = widget.LowImportance

	card.run = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() {
		u.runHistoryEntry(card)
	})
	card.edit = widget.NewButtonWithIcon("Edit", theme.DocumentCreateIcon(), func() {
		u.editHistoryEntry(card)
	})
	card.restore = widget.NewButtonWithIcon("Restore this state…", theme.HistoryIcon(), func() {
		u.showRevert(card)
	})

	card.hint = widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
	card.hint.Wrapping = fyne.TextWrapWord

	u.updateHistoryActions(card, detail)

	return container.NewVBox(
		grid,
		card.rawToggle,
		card.raw,
		container.NewHBox(card.run, card.edit, card.restore),
		card.hint,
	)
}

// updateHistoryActions enables only what this entry can do: Run needs an execute
// entry whose command still exists, Edit needs the command to exist. A disabled
// action says why.
func (u *ui) updateHistoryActions(card *historyCard, detail desktop.HistoryDetail) {
	if card.run == nil || card.edit == nil {
		return
	}

	command, exists := u.model.Commands.Get(detail.Command)

	if exists && detail.Operation == string(operation.Execute) {
		card.run.Enable()
	} else {
		card.run.Disable()
	}

	if exists {
		card.edit.Enable()
	} else {
		card.edit.Disable()
	}

	card.restore.Enable()

	card.hint.SetText(historyActionHint(detail, exists, command))
	if card.hint.Text == "" {
		card.hint.Hide()
	} else {
		card.hint.Show()
	}
}

// historyActionHint explains a disabled action, or says nothing when all apply.
func historyActionHint(detail desktop.HistoryDetail, exists bool, command catalog.Command) string {
	if !exists {
		return "This command is not in the current setup, so it cannot be run or edited."
	}

	if detail.Operation != string(operation.Execute) {
		return "Run replays execute entries; this entry recorded a " + detail.Operation + "."
	}

	if !detail.ParamsRecorded && len(command.Params) > 0 {
		return "This entry did not record its parameters; Run asks for them."
	}

	return ""
}

// runHistoryEntry runs the managed command again with the parameters the entry
// recorded. It never re-creates or deletes.
func (u *ui) runHistoryEntry(card *historyCard) {
	command, found := u.model.Commands.Get(card.detail.Command)
	if !found {
		return
	}

	if !card.detail.ParamsRecorded && len(command.Params) > 0 {
		u.askParams(card.detail.Command, command)
		return
	}

	u.model.RunWithParams(card.detail.Command, card.detail.Params)
	u.showOutput()
}

// editHistoryEntry takes the entry's command to the Commands tab and opens its
// edit dialog.
func (u *ui) editHistoryEntry(card *historyCard) {
	if _, found := u.model.Commands.Get(card.detail.Command); !found {
		return
	}

	u.selectCommand(card.detail.Command)
	u.showEditFor(card.detail.Command)
}

// revertSteps is how many entries back the card is, which is what the Reverter
// expects (the newest entry is one step).
func (u *ui) revertSteps(card *historyCard) int {
	return card.detail.Position + 1
}

// showRevert asks before restoring, then applies the change through the
// Reverter port.
func (u *ui) showRevert(card *historyCard) {
	detail := card.detail
	steps := u.revertSteps(card)

	dropped := steps - 1

	detailText := "No newer entries are dropped."
	switch {
	case dropped == 1:
		detailText = "One newer entry is dropped."
	case dropped > 1:
		detailText = fmt.Sprintf("%d newer entries are dropped.", dropped)
	}

	question := fmt.Sprintf(
		"Restore the command set recorded at version %d (%s, %s)?\n\n%s",
		detail.Version, detail.Operation, detail.Command, detailText,
	)

	dialog.ShowConfirm("Restore command set", question, func(ok bool) {
		if !ok {
			return
		}

		u.applyRevert(steps)
	}, u.window)
}

// applyRevert is the confirmation's submit path.
func (u *ui) applyRevert(steps int) bool {
	commands, history, err := u.reverter.Revert(steps)
	if err != nil {
		showError(err.Error(), u.window)
		return false
	}

	u.reload(commands, history)
	return true
}

// historyFieldLabel is what the grid calls a stored key.
func historyFieldLabel(key string) string {
	if label, ok := historyFieldLabels[key]; ok {
		return label
	}

	words := strings.Split(key, "_")
	for i, word := range words {
		if word == "" {
			continue
		}

		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		words[i] = string(runes)
	}

	return strings.Join(words, " ")
}

var historyFieldLabels = map[string]string{
	"main_operation":    "Operation",
	"sub_operation":     "Kind",
	"command_name":      "Name",
	"command_alias":     "Aliases",
	"command_exec":      "Exec",
	"command_steps":     "Steps",
	"command_shell":     "Shell",
	"command_env":       "Env keys",
	"command_timeout":   "Timeout",
	"command_params":    "Params",
	"old_command_name":  "Old name",
	"old_command_alias": "Old aliases",
	"old_command_exec":  "Old exec",
	"crower_exec":       "Recorded as",
}

// historyFieldValue renders a stored value for a reader: lists are lists, not
// Go's bracket syntax or raw JSON.
func historyFieldValue(key, value string) string {
	switch key {
	case "command_alias", "old_command_alias":
		return strings.TrimSpace(strings.Trim(value, "[]"))
	case "command_steps", "command_params":
		var items []string
		if err := json.Unmarshal([]byte(value), &items); err == nil {
			return strings.Join(items, "\n")
		}
	}

	return value
}
