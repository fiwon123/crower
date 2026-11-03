package utils

import (
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// UnknownProcessField is what a row carries when the OS will not report the
// name or the path of a running process.
const UnknownProcessField = "Unknown"

// ProcessInfo is one row of the process table: the data `list --process` shows,
// and the data its JSON result carries.
type ProcessInfo struct {
	PID  int32  `json:"pid"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// ListProcesses returns the running processes, optionally filtered by a
// case-insensitive part of the name.
//
// It returns rows rather than a rendered table because the table is only one
// presentation of this data: a JSON result needs the data itself, and keeping
// the layout here would have forced `list --process --json` to parse its own
// output back into fields.
func ListProcesses(partName string, skipUnknown bool) ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	rows := make([]ProcessInfo, 0, len(processes))

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			name = UnknownProcessField
		}

		// skip if its not the process
		if partName != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(partName)) {
			continue
		}

		exe, err := p.Exe()
		if err != nil {
			exe = UnknownProcessField
		}

		if skipUnknown && (name == UnknownProcessField || exe == UnknownProcessField) {
			continue
		}

		rows = append(rows, ProcessInfo{PID: p.Pid, Name: name, Path: exe})
	}

	return rows, nil
}

func GetProcessNameByID(id int32) (string, error) {
	processes, err := process.Processes()
	if err != nil {
		return "", err
	}

	for _, p := range processes {
		pid := p.Pid

		name, err := p.Name()
		if err != nil {
			name = ""
		}

		if pid == id {
			return name, nil
		}
	}

	return "", nil
}

func GetProcessPathByID(id int32) (string, error) {
	processes, err := process.Processes()
	if err != nil {
		return "", err
	}

	for _, p := range processes {
		pid := p.Pid

		exe, err := p.Exe()
		if err != nil {
			exe = ""
		}

		if pid == id {
			return exe, nil
		}
	}

	return "", nil
}

func GetProcessPathByName(findName string) (string, error) {
	processes, err := process.Processes()
	if err != nil {
		return "", err
	}

	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			name = ""
		}

		exe, err := p.Exe()
		if err != nil {
			exe = ""
		}

		if name == findName {
			return exe, nil
		}
	}

	return "", nil
}
