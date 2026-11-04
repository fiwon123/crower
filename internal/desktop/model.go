package desktop

import (
	"fmt"
	"strings"

	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/domain/timeline"
)

// Executor runs a managed command for the window.
//
// It is declared here, at the point of use, so this package depends on neither
// the use case nor its adapters; the command layer supplies one.
type Executor interface {
	// Execute runs key with params and returns the output and the command's
	// name.
	Execute(key string, params []string) (string, string, error)
}

type ViewState int

const (
	ViewList ViewState = iota
	ViewRun
	ViewParams
	ViewHistory
	ViewOutput
)

type Model struct {
	Commands catalog.Catalog
	History  timeline.Log
	Exec     Executor
	State    ViewState
	Keys     []string
	Index    int
	Filter   string

	// run
	SelectedKey string
	Params      []string
	ParamInput  string

	// output
	LastOut string
	LastErr error
	LastCmd string
}

// NewModel builds the window over the current command set and history.
func NewModel(commands catalog.Catalog, history timeline.Log, exec Executor) *Model {
	return &Model{
		Commands: commands,
		History:  history,
		Exec:     exec,
		State:    ViewList,
		Keys:     commands.Names(),
	}
}

func (m *Model) FilteredKeys() []string {
	if m.Filter == "" {
		return m.Keys
	}
	f := strings.ToLower(m.Filter)
	var out []string
	for _, k := range m.Keys {
		if strings.Contains(strings.ToLower(k), f) {
			out = append(out, k)
		}
	}
	if len(out) == 0 {
		return m.Keys
	}
	return out
}

func (m *Model) Move(d int) {
	keys := m.FilteredKeys()
	if len(keys) == 0 {
		return
	}
	m.Index = (m.Index + d) % len(keys)
	if m.Index < 0 {
		m.Index = len(keys) + m.Index
	}
}

func (m *Model) Selected() string {
	keys := m.FilteredKeys()
	if len(keys) == 0 {
		return ""
	}
	if m.Index < 0 || m.Index >= len(keys) {
		m.Index = 0
	}
	return keys[m.Index]
}

func (m *Model) StartRun(withParams bool) {
	sel := m.Selected()
	if sel == "" {
		return
	}
	m.SelectedKey = sel
	m.Params = nil
	m.ParamInput = ""
	if withParams {
		m.State = ViewParams
		return
	}
	m.execute()
}

func (m *Model) AddParamChar(c rune) {
	if c >= 32 && c < 127 {
		m.ParamInput += string(c)
	}
}

func (m *Model) BackspaceParam() {
	if len(m.ParamInput) > 0 {
		m.ParamInput = m.ParamInput[:len(m.ParamInput)-1]
	}
}

func (m *Model) SubmitParams() {
	if strings.TrimSpace(m.ParamInput) != "" {
		m.Params = strings.Fields(m.ParamInput)
	} else {
		m.Params = nil
	}
	m.execute()
}

func (m *Model) execute() {
	out, name, err := m.Exec.Execute(m.SelectedKey, m.Params)
	m.LastOut = out
	m.LastErr = err
	if name != "" {
		m.LastCmd = name
	} else {
		m.LastCmd = m.SelectedKey
	}
	m.State = ViewOutput
}

// historyText renders the history screen. It is a method so a test can check
// that a command operation shows which command it touched.
func (m *Model) historyText() string {
	var sb strings.Builder

	sb.WriteString("[HISTORY] Enter/Esc=back\n\n")

	if len(m.History.AllData) == 0 {
		sb.WriteString("No history\n")
		return sb.String()
	}

	for i, entry := range m.History.AllData {
		if i > 20 {
			break
		}

		commandName := entry.CommandName
		if commandName == "" {
			commandName = "-"
		}

		sb.WriteString(fmt.Sprintf("%d. [%s] %s %s\n", i+1, entry.Operation, commandName, entry.Timestemp))
	}

	return sb.String()
}
