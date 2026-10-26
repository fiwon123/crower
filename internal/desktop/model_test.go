package desktop

import (
	"testing"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/data/history"
	"github.com/fiwon123/crower/pkg/logger"
)

func testApp(t *testing.T) *app.Config {
	t.Helper()
	a := &app.Config{}
	a.AllCommandsByName = command.NewMapData()
	a.AllCommandsByAlias = command.NewMapData()
	a.OrderKeys = []string{}
	a.History = history.New(t.TempDir() + "/hist")
	a.Logger = *logger.New(t.TempDir() + "/log")
	a.HistoryFilePath = t.TempDir() + "/hist.json"
	a.HistoryFolderPath = t.TempDir() + "/hist"
	c := command.New("hello", nil, "echo hello")
	a.AllCommandsByName.Add("hello", c)
	a.OrderKeys = append(a.OrderKeys, "hello")
	return a
}

func TestNewModel(t *testing.T) {
	a := testApp(t)
	m := NewModel(a)
	if m.State != ViewList {
		t.Fatalf("want list state")
	}
	if len(m.Keys) != 1 {
		t.Fatalf("want 1 key")
	}
}

func TestMove(t *testing.T) {
	a := testApp(t)
	a.AllCommandsByName.Add("b", command.New("b", nil, "echo b"))
	a.OrderKeys = append(a.OrderKeys, "b")
	m := NewModel(a)
	m.Move(1)
	if m.Selected() != "b" {
		t.Fatalf("want b")
	}
	m.Move(1)
	if m.Selected() != "hello" {
		t.Fatalf("want wrap")
	}
}

func TestFilter(t *testing.T) {
	a := testApp(t)
	a.AllCommandsByName.Add("build", command.New("build", nil, "make build"))
	a.OrderKeys = append(a.OrderKeys, "build")
	m := NewModel(a)
	m.Filter = "bui"
	keys := m.FilteredKeys()
	if len(keys) != 1 || keys[0] != "build" {
		t.Fatalf("bad filter")
	}
}
