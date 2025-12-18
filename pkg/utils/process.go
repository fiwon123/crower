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
