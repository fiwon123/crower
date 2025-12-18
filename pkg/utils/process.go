package utils

import (
	"fmt"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

func ListAllProcess(partName string, skipUnknown bool) error {
	const (
		unknow        = "Unknown"
		titleLayout   = "%-8s %-30s %s\n"
		contentLayout = "%-8d %-30s %s\n"
	)

	processes, err := process.Processes()
	if err != nil {
		return err
	}

	fmt.Printf(titleLayout, "PID", "Name", "Path")
	fmt.Println("-------------------------------------------------------------")

	for _, p := range processes {
		pid := p.Pid

		name, err := p.Name()
		if err != nil {
			name = unknow
		}

		// skip if its not the process
		if partName != "" && !strings.Contains(strings.ToLower(name), strings.ToLower(partName)) {
			continue
		}

		exe, err := p.Exe()
		if err != nil {
			exe = unknow
		}

		if skipUnknown && (name == unknow || exe == unknow) {
			continue
		}

		fmt.Printf(contentLayout, pid, name, exe)
	}

	return nil
}
