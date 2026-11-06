//go:build fyne && !windows

package cmd

import (
	"fmt"
	"os"
)

// A graphical session still has a stderr to inherit, so the message goes there.
func showFyneStartupError(message string) {
	fmt.Fprintln(os.Stderr, message)
}
