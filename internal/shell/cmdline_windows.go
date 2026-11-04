//go:build windows

package shell

import (
	"os/exec"
	"syscall"
)

// prepareWindows hands cmd the line exactly as crower built it.
//
// exec.Command("cmd", "/c", line) makes Go escape line for CreateProcess, and
// cmd then applies its own quoting rules to the escaped text. A path with a
// space or a value containing quotes therefore reaches the program split or
// mangled. CmdLine is the raw lpCommandLine CreateProcess receives, so cmd sees
// the quotes QuoteFor wrote (#185).
func prepareWindows(command *exec.Cmd, line string) {
	if len(command.Args) == 0 || command.Args[0] != "cmd" {
		return
	}

	command.SysProcAttr = &syscall.SysProcAttr{CmdLine: "cmd /c " + line}
}
