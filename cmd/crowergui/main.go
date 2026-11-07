//go:build fyne

// Command crowergui is the Fyne desktop application.
//
// It is its own entry point: no command tree sits between a double-click and the
// window. --config and --setup still work; every other argument is ignored.
package main

import (
	"os"

	"github.com/fiwon123/crower/cmd"
)

func main() {
	cmd.RunWindow(os.Args[1:])
}
