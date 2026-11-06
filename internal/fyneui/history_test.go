//go:build fyne

package fyneui

import (
	"slices"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// historyUI lays out the window over a log and opens History.
func historyUI(t *testing.T, exec desktop.Executor, commands []catalog.Command, history timeline.Log) (*ui, fyne.Window) {
	t.Helper()

	if commands == nil {
		commands = []catalog.Command{{Name: "greet", Exec: "echo hi"}}
	}

	model := desktop.NewModel(catalog.New(nil, commands), history, exec)

	front, window := newTestUI(t, model, stubEditor{}, stubReverter{})
	window.SetContent(front.content())
	window.Resize(fyneSize())
	front.tabs.SelectIndex(1)

	return front, window
}

// historyEntryUI lays out the window over one recorded entry.
func historyEntryUI(t *testing.T, exec desktop.Executor, commands []catalog.Command, op operation.Operation, command, note string) (*ui, fyne.Window) {
	t.Helper()

	var history timeline.Log
	history.Add(op, command, note)

	return historyUI(t, exec, commands, history)
}

// openHistoryCard opens the card at index and returns it.
func openHistoryCard(t *testing.T, front *ui, index int) *historyCard {
	t.Helper()

	if index < 0 || index >= len(front.historyCards) {
		t.Fatalf("card %d does not exist; %d cards", index, len(front.historyCards))
	}

	card := front.historyCards[index]
	front.toggleHistoryCard(card)

	return card
}

// Opening an entry builds its detail inline; nothing is built while collapsed.
func TestHistoryEntryOpensInline(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, []string{"a b"})
	front, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", note)

	if len(front.historyCards) != 1 {
		t.Fatalf("cards = %d, want 1", len(front.historyCards))
	}

	card := front.historyCards[0]
	if card.isOpen() {
		t.Fatal("the entry is open before it was tapped")
	}

	if len(card.body.Objects) != 0 {
		t.Fatal("a collapsed card built its body")
	}

	card.header.OnTapped()

	if !card.isOpen() || !card.body.Visible() {
		t.Fatal("tapping the header did not open the card")
	}

	if len(card.body.Objects) == 0 {
		t.Fatal("the opened card has no body")
	}
}

// Opening another entry closes the first, so the page holds one detail at a time.
func TestHistoryOpeningAnotherCardClosesTheFirst(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)

	var history timeline.Log
	history.Add(operation.Execute, "greet", note)
	history.Add(operation.Execute, "greet", note)

	front, _ := historyUI(t, stubExecutor{}, nil, history)

	first := openHistoryCard(t, front, 0)
	second := openHistoryCard(t, front, 1)

	if first.isOpen() {
		t.Error("the first card stayed open")
	}

	if !second.isOpen() {
		t.Error("the second card did not open")
	}
}

// The note is shown as labelled fields, with the raw line behind a toggle.
func TestHistoryCardLabelsTheNoteFields(t *testing.T) {
	command := catalog.Command{Name: "greet", Aliases: []string{"hi"}, Exec: "echo hi"}
	note := notes.ExecuteCommandNote(command, []string{"a b"})

	front, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", note)
	card := openHistoryCard(t, front, 0)

	labels := labelsIn(card.body)
	for _, want := range []string{"execute", "command", "greet", "hi", "echo hi", "a b"} {
		if !slices.Contains(labels, want) {
			t.Errorf("fields %q do not contain %q", labels, want)
		}
	}

	if card.raw.Visible() {
		t.Error("the raw note is open before the toggle was pressed")
	}

	if !card.raw.Selectable {
		t.Error("the raw note is not selectable")
	}

	if !strings.Contains(card.raw.Text, "command_exec=echo hi") {
		t.Errorf("raw note = %q, want the key=value line", card.raw.Text)
	}
}

// An entry whose note is empty says so.
func TestHistoryCardShowsNoNotePlaceholder(t *testing.T) {
	front, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", "")
	card := openHistoryCard(t, front, 0)

	if labels := labelsIn(card.body); !slices.Contains(labels, "(no note)") {
		t.Errorf("fields %q do not say the entry has no note", labels)
	}
}

// Run replays the entry's command with the arguments it recorded.
func TestHistoryRunReplaysTheRecordedParams(t *testing.T) {
	exec := &recordingExecutor{}
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, []string{"a b", "-x"})

	front, _ := historyEntryUI(t, exec, nil, operation.Execute, "greet", note)
	card := openHistoryCard(t, front, 0)

	if card.run.Disabled() {
		t.Fatal("run is disabled for an execute entry whose command exists")
	}

	card.run.OnTapped()

	if !slices.Equal(exec.keys, []string{"greet"}) {
		t.Fatalf("executed %v, want [greet]", exec.keys)
	}

	if !slices.Equal(exec.args[0], []string{"a b", "-x"}) {
		t.Errorf("params = %v, want the recorded arguments", exec.args[0])
	}
}

// An execute entry written before command_params existed asks for them when the
// command declares any.
func TestHistoryRunAsksWhenParamsWereNotRecorded(t *testing.T) {
	exec := &recordingExecutor{}
	commands := []catalog.Command{{Name: "greet", Exec: "echo {1}", Params: []string{"who"}}}
	legacy := "main_operation=execute;sub_operation=command;command_name=greet;command_exec=echo {1}"

	front, window := historyEntryUI(t, exec, commands, operation.Execute, "greet", legacy)
	card := openHistoryCard(t, front, 0)

	if card.run.Disabled() {
		t.Fatal("run is disabled for an execute entry whose command exists")
	}

	card.run.OnTapped()

	if len(exec.keys) != 0 {
		t.Errorf("executed %v before the parameters were supplied", exec.keys)
	}

	if front.model.State != desktop.ViewParams {
		t.Errorf("state = %v, want the params prompt", front.model.State)
	}

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("run did not open the parameters dialog")
	}
}

// Edit clears the Commands filters, selects the entry's command and opens its
// dialog.
func TestHistoryEditOpensTheCommand(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)
	commands := []catalog.Command{
		{Name: "alpha", Exec: "echo a"},
		{Name: "greet", Exec: "echo hi", Group: "work"},
	}

	front, window := historyEntryUI(t, stubExecutor{}, commands, operation.Execute, "greet", note)

	front.model.Group = "work"
	front.refreshGroups()

	front.filterEntry.SetText("zzz")

	card := openHistoryCard(t, front, 0)
	card.edit.OnTapped()

	if front.tabs.SelectedIndex() != 0 {
		t.Errorf("selected tab = %d, want Commands", front.tabs.SelectedIndex())
	}

	if front.model.Filter != "" || front.filterEntry.Text != "" {
		t.Errorf("filter = %q/%q, want it cleared", front.model.Filter, front.filterEntry.Text)
	}

	if front.model.Group != "" || front.groupFilter.Selected != allGroupsLabel {
		t.Errorf("group = %q/%q, want the all-groups entry", front.model.Group, front.groupFilter.Selected)
	}

	if got := front.model.Selected(); got != "greet" {
		t.Errorf("selected command = %q, want greet", got)
	}

	if front.commandOpen == nil || front.commandOpen.key != "greet" {
		t.Error("the greet card is not open after Edit")
	}

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("edit did not open a dialog")
	}
}

// An action is offered only where it can work, and a disabled one says why.
func TestHistoryActionsFollowTheEntry(t *testing.T) {
	tests := []struct {
		name    string
		op      operation.Operation
		command string
		note    string
		run     bool
		edit    bool
		hint    string
	}{
		{
			"execute with a live command",
			operation.Execute, "greet",
			"main_operation=execute;sub_operation=command;command_name=greet;command_params=[]",
			true, true, "",
		},
		{
			"copy naming a live command",
			operation.Copy, "greet",
			"main_operation=copy",
			false, true, "execute entries",
		},
		{
			"a command that no longer exists",
			operation.Execute, "gone",
			"main_operation=execute;sub_operation=command;command_name=gone",
			false, false, "not in the current setup",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			front, _ := historyEntryUI(t, stubExecutor{}, nil, tt.op, tt.command, tt.note)
			card := openHistoryCard(t, front, 0)

			if !card.hint.Selectable {
				t.Error("the disabled reason is not selectable")
			}

			if got := !card.run.Disabled(); got != tt.run {
				t.Errorf("run enabled = %v, want %v", got, tt.run)
			}

			if got := !card.edit.Disabled(); got != tt.edit {
				t.Errorf("edit enabled = %v, want %v", got, tt.edit)
			}

			if tt.hint == "" {
				if card.hint.Text != "" || card.hint.Visible() {
					t.Errorf("hint = %q, want none when all actions apply", card.hint.Text)
				}

				return
			}

			if !strings.Contains(card.hint.Text, tt.hint) {
				t.Errorf("hint = %q, want it to contain %q", card.hint.Text, tt.hint)
			}

			if !card.hint.Visible() {
				t.Error("a disabled action left no visible reason")
			}
		})
	}
}

// The Raw note toggle opens and closes the stored line, and the page lays the
// opened label out.
func TestHistoryRawNoteToggle(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)

	front, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", note)
	card := openHistoryCard(t, front, 0)

	if card.raw.Visible() {
		t.Fatal("the raw note starts open")
	}

	card.rawToggle.OnTapped()
	if !card.raw.Visible() || !card.rawShown {
		t.Error("the toggle did not open the raw note")
	}

	if card.raw.Size().Height == 0 {
		t.Error("the opened raw note has no height: the page did not lay it out")
	}

	card.rawToggle.OnTapped()
	if card.raw.Visible() || card.rawShown {
		t.Error("the toggle did not close the raw note")
	}
}

// The page renders the newest hundred entries and loads the rest on request.
func TestHistoryPagesOlderEntries(t *testing.T) {
	var history timeline.Log
	for i := 0; i < 150; i++ {
		history.Add(operation.Create, "greet", "note")
	}

	front, _ := historyUI(t, stubExecutor{}, nil, history)

	if len(front.historyCards) != historyPageSize {
		t.Fatalf("cards = %d, want the first page (%d)", len(front.historyCards), historyPageSize)
	}

	if !strings.Contains(front.historyCards[0].header.Text, "#150") {
		t.Errorf("first card = %q, want the newest entry (#150)", front.historyCards[0].header.Text)
	}

	if !strings.Contains(front.historyCards[99].header.Text, "#51") {
		t.Errorf("last card on page 1 = %q, want #51", front.historyCards[99].header.Text)
	}

	if !front.historyLoadMore.Visible() {
		t.Error("Load older is hidden while entries remain")
	}

	front.loadOlderHistory()

	if len(front.historyCards) != 150 {
		t.Errorf("cards = %d, want all 150 after loading older", len(front.historyCards))
	}

	if !strings.Contains(front.historyCards[149].header.Text, "#1") {
		t.Errorf("last card = %q, want the oldest entry (#1)", front.historyCards[149].header.Text)
	}

	if front.historyLoadMore.Visible() {
		t.Error("Load older is still visible with nothing left")
	}
}

// A reload (an edit, a revert or a setup switch) rebuilds the rendered cards.
func TestHistoryReloadRebuildsTheCards(t *testing.T) {
	front, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Create, "greet", "note")

	if len(front.historyCards) != 1 {
		t.Fatalf("cards = %d, want 1", len(front.historyCards))
	}

	var next timeline.Log
	next.Add(operation.Copy, "", "note")
	next.Add(operation.Delete, "greet", "note")

	front.reload(catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}), next)

	if len(front.historyCards) != 2 {
		t.Fatalf("cards = %d, want the new page", len(front.historyCards))
	}

	if !strings.Contains(front.historyCards[0].header.Text, "delete") {
		t.Errorf("first card = %q, want the newest entry", front.historyCards[0].header.Text)
	}
}

// Restore counts the entry's distance from the newest; the newest is one step.
func TestHistoryRestoreCountsFromTheNewest(t *testing.T) {
	note := notes.ExecuteCommandNote(catalog.Command{Name: "greet", Exec: "echo hi"}, nil)

	var history timeline.Log
	for i := 0; i < 3; i++ {
		history.Add(operation.Execute, "greet", note)
	}

	front, window := historyUI(t, stubExecutor{}, nil, history)

	newest := openHistoryCard(t, front, 0)
	if got := front.revertSteps(newest); got != 1 {
		t.Errorf("steps for the newest = %d, want 1", got)
	}

	second := openHistoryCard(t, front, 1)
	if got := front.revertSteps(second); got != 2 {
		t.Errorf("steps for the second = %d, want 2", got)
	}

	second.restore.OnTapped()

	if len(window.Canvas().Overlays().List()) == 0 {
		t.Error("restore did not ask for confirmation")
	}
}

// The summary uses the operation icon and the clock, not the full date; the
// expanded fields carry the date and version.
func TestHistorySummaryTitleIsCompact(t *testing.T) {
	title := historySummaryTitle(desktop.HistorySummary{
		Version:   7,
		When:      "2026-11-06 12:00",
		Operation: "execute",
		Command:   "greet",
	})

	for _, want := range []string{"▶️", "execute", "greet", "12:00", "#7"} {
		if !strings.Contains(title, want) {
			t.Errorf("title = %q, want it to contain %q", title, want)
		}
	}

	if strings.Contains(title, "2026-11-06") {
		t.Errorf("title = %q, want the clock only", title)
	}
}

// A long command is shortened in the summary so the header cannot widen the
// window's minimum.
func TestHistorySummaryTitleTruncatesLongCommands(t *testing.T) {
	title := historySummaryTitle(desktop.HistorySummary{
		Version:   7,
		When:      "2026-11-06 12:00",
		Operation: "execute",
		Command:   strings.Repeat("a", 60),
	})

	if !strings.Contains(title, "…") {
		t.Errorf("title = %q, want an ellipsis", title)
	}

	if len([]rune(title)) > 90 {
		t.Errorf("title = %q, want it bounded", title)
	}
}

// An empty log says so instead of showing an empty page.
func TestHistoryEmptyState(t *testing.T) {
	front, _ := historyUI(t, stubExecutor{}, nil, timeline.Log{})

	if len(front.historyCards) != 0 {
		t.Errorf("cards = %d, want none", len(front.historyCards))
	}

	if !front.historyEmpty.Visible() {
		t.Error("the empty state is hidden")
	}
}

// The History tab keeps its header and cards off the window edge.
func TestHistoryHasALeftMargin(t *testing.T) {
	_, window := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", "main_operation=execute")

	pos, found := labelPosition(window.Canvas().Content(), historyHeader)
	if !found {
		t.Fatal("the history header is not in the rendered tree")
	}

	if pos.X < theme.Padding()*2 {
		t.Errorf("header x = %.1f, want at least %.1f from the window edge", pos.X, theme.Padding()*2)
	}
}

// The stored list syntax is not what a reader should see.
func TestHistoryFieldValueFormatsLists(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  string
	}{
		{"aliases", "command_alias", "[a b]", "a b"},
		{"steps", "command_steps", `["echo a","echo b"]`, "echo a\necho b"},
		{"params", "command_params", `["a b","-x"]`, "a b\n-x"},
		{"corrupt json is shown as stored", "command_params", "[bad", "[bad"},
		{"a plain value is unchanged", "command_exec", "echo hi", "echo hi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := historyFieldValue(tt.key, tt.value); got != tt.want {
				t.Errorf("historyFieldValue(%q, %q) = %q, want %q", tt.key, tt.value, got, tt.want)
			}
		})
	}
}

// Unknown keys are humanised rather than dropped, so a newer note still reads.
func TestHistoryFieldLabel(t *testing.T) {
	tests := map[string]string{
		"main_operation": "Operation",
		"command_name":   "Name",
		"future_field":   "Future Field",
		"é_field":        "É Field",
	}

	for key, want := range tests {
		if got := historyFieldLabel(key); got != want {
			t.Errorf("historyFieldLabel(%q) = %q, want %q", key, got, want)
		}
	}
}
