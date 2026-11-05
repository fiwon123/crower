//go:build fyne

// Command crowergui is the Fyne desktop application. It shares the CLI's
// configuration and catalog through the root command, so --config and --setup
// work the same way.
package main

import (
	"os"

	"github.com/fiwon123/crower/cmd"
)

func main() {
	cmd.ExecuteFyne(os.Args[1:])
}
