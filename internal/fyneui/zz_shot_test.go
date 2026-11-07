//go:build fyne

package fyneui

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/software"

	"github.com/fiwon123/crower/internal/adapters/notes"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/operation"
)

func captureTab(t *testing.T, content fyne.CanvasObject, name string) {
	t.Helper()

	canvas := software.NewCanvas()
	canvas.SetContent(content)
	canvas.Resize(fyne.NewSize(960, 720))

	image := canvas.Capture()

	path := filepath.Join("..", "..", "build", "shots", name+".png")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, image); err != nil {
		t.Fatalf("encode: %v", err)
	}
}

func TestScreenshots(t *testing.T) {
	front, _ := commandUI(t, stubExecutor{}, []catalog.Command{
		{Name: "greet", Exec: "echo hi", Icon: "🚀", Group: "deploy", Aliases: []string{"hi"}},
		{Name: "deploy", Exec: "make deploy", Icon: "📦", Group: "deploy"},
		{Name: "backup", Exec: "tar czf backup.tgz .", Group: "backup"},
	})

	for _, card := range front.commandCards {
		if card.key == "greet" {
			front.toggleCommandCard(card)
			break
		}
	}

	captureTab(t, front.window.Canvas().Content(), "commands")

	command := catalog.Command{Name: "greet", Exec: "echo hi", Aliases: []string{"hi"}}
	note := notes.ExecuteCommandNote(command, []string{"a b"})

	history, _ := historyEntryUI(t, stubExecutor{}, nil, operation.Execute, "greet", note)
	if len(history.historyCards) > 0 {
		history.toggleHistoryCard(history.historyCards[0])
	}

	captureTab(t, history.window.Canvas().Content(), "history")

	tools, _ := commandUI(t, stubExecutor{}, nil)
	captureTab(t, tools.window.Canvas().Content(), "tools")
}
