package desktop

import (
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/domain/executedom"
)

type ViewState int

const (
	ViewList ViewState = iota
	ViewRun
	ViewParams
	ViewHistory
	ViewOutput
)

type Model struct {
	App    *app.Config
	Exec   *executedom.Handler
	State  ViewState
	Keys   []string
	Index  int
	Filter string

	// run
	SelectedKey string
	Params      []string
	ParamInput  string

	// output
	LastOut string
	LastErr error
	LastCmd string
}

func NewModel(a *app.Config) *Model {
	m := &Model{
		App:   a,
		Exec:  executedom.NewHandler(a),
		State: ViewList,
		Keys:  append([]string(nil), a.OrderKeys...),
	}
	return m
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
	out, cmd, err := m.Exec.Execute(m.SelectedKey, m.Params)
	m.LastOut = out
	m.LastErr = err
	if cmd != nil {
		m.LastCmd = cmd.Name
	} else {
		m.LastCmd = m.SelectedKey
	}
	m.State = ViewOutput
}
