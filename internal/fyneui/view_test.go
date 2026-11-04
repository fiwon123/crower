//go:build fyne

package fyneui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"

	"github.com/fiwon123/crower/internal/desktop"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

type stubExecutor struct{}

func (stubExecutor) Execute(string, []string) (string, string, error) {
	return "", "", nil
}

// The view builds over the headless model and attaches to a window. The model
// itself is covered by internal/desktop; this is the adapter smoke test.
func TestNewContentBuildsOverTheModel(t *testing.T) {
	test.NewApp()

	model := desktop.NewModel(
		catalog.New(nil, []catalog.Command{{Name: "greet", Exec: "echo hi"}}),
		timeline.Log{},
		stubExecutor{},
	)

	window := test.NewWindow(nil)
	content := NewContent(model, window)
	if content == nil {
		t.Fatal("NewContent returned nil")
	}

	window.SetContent(content)
	window.Resize(fyne.NewSize(720, 520))

	if window.Canvas().Content() == nil {
		t.Fatal("window has no content")
	}
}
