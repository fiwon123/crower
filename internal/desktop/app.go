package desktop

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	model *Model
}

func NewGame(m *Model) *Game {
	return &Game{model: m}
}

func (g *Game) Update() error {
	m := g.model
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		if m.State == ViewList {
			return ebiten.Termination
		}
		m.State = ViewList
		return nil
	}
	if m.State == ViewList {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) || inpututil.IsKeyJustPressed(ebiten.KeyJ) {
			m.Move(1)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) || inpututil.IsKeyJustPressed(ebiten.KeyK) {
			m.Move(-1)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.StartRun(false)
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
			m.StartRun(true)
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyH) {
			m.State = ViewHistory
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySlash) || inpututil.IsKeyJustPressed(ebiten.KeyF) {
			// simple clear filter on key press for MVP
		}
		return nil
	}
	if m.State == ViewParams {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			m.SubmitParams()
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			m.BackspaceParam()
			return nil
		}
		// append chars
		runes := ebiten.AppendInputChars(nil)
		for _, r := range runes {
			m.AddParamChar(r)
		}
		return nil
	}
	if m.State == ViewHistory || m.State == ViewOutput {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			m.State = ViewList
			return nil
		}
		return nil
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	m := g.model
	var sb strings.Builder
	sb.WriteString("CROWER DESKTOP\n\n")
	if m.State == ViewList {
		sb.WriteString("[LIST] Enter=run, P=run with params, H=history, Esc/q=quit\n\n")
		keys := m.FilteredKeys()
		if len(keys) == 0 {
			sb.WriteString("No commands\n")
		}
		for i, k := range keys {
			cmd, found := m.Commands.Get(k)
			prefix := "  "
			if i == m.Index && len(keys) > 0 {
				// find index in filtered
				fk := m.FilteredKeys()
				if m.Index < len(fk) && fk[m.Index] == k {
					prefix = "> "
				}
			}
			if found {
				sb.WriteString(fmt.Sprintf("%s%d. %s - %s\n", prefix, i+1, cmd.Name, cmd.Exec))
			} else {
				sb.WriteString(fmt.Sprintf("%s%d. %s\n", prefix, i+1, k))
			}
		}
	} else if m.State == ViewParams {
		sb.WriteString("[PARAMS] Enter=execute, Backspace=edit, Esc=cancel\n\n")
		sb.WriteString(fmt.Sprintf("Command: %s\n\n", m.SelectedKey))
		sb.WriteString("Parameters: " + m.ParamInput + "_\n")
	} else if m.State == ViewHistory {
		sb.WriteString(m.historyText())
	} else if m.State == ViewOutput {
		sb.WriteString("[OUTPUT] Enter/Esc=back\n\n")
		sb.WriteString(fmt.Sprintf("Cmd: %s\n\n", m.LastCmd))
		if m.LastOut != "" {
			sb.WriteString(m.LastOut)
			if !strings.HasSuffix(m.LastOut, "\n") {
				sb.WriteString("\n")
			}
		}
		if m.LastErr != nil {
			sb.WriteString(fmt.Sprintf("\nError: %v\n", m.LastErr))
		}
	}
	ebitenutil.DebugPrint(screen, sb.String())
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

func Run(commands catalog.Catalog, history timeline.Log, exec Executor) error {
	m := NewModel(commands, history, exec)
	game := NewGame(m)
	ebiten.SetWindowTitle("Crower Desktop")
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowResizable(true)
	return ebiten.RunGame(game)
}
