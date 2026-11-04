//go:build !windows

package shell

import "os/exec"

// prepareWindows is a no-op away from Windows. The Windows branch is still
// compiled on every platform so the command table can be tested with a goos
// parameter, but only a real Windows run uses CmdLine.
func prepareWindows(*exec.Cmd, string) {}
