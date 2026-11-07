//go:build fyne

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fiwon123/crower/internal/fyneui"
	"github.com/fiwon123/crower/internal/session"
)

// windowHelp is what the window binary answers to --help. The window itself has
// no subcommands, so this is the whole grammar.
const windowHelp = `crower-gui opens the crower window.

	crower-gui                 open the window
	crower-gui --config FILE   use another configuration file
	crower-gui --setup NAME    use a named setup

The window starts whatever the arguments are: a flag it does not know is
ignored, so a double-click always reaches it. The terminal interface is
"crower gui"; the command-line operations are in "crower --help".
`

// RunWindow is the crower-gui entry point.
//
// It is not a command tree. A window has no subcommands, and a double-click must
// not depend on argument parsing reaching one: the reference is the sibling
// project's #382, where the window only opened when it was typed with a
// subcommand (#465). Two flags are still read because they were useful; every
// other argument is ignored.
func RunWindow(args []string) {
	if hasWindowArg(args, "--version", "-v") {
		fmt.Fprintln(streams.Out, Version)
		return
	}

	if hasWindowArg(args, "--help", "-h") {
		fmt.Fprint(streams.Out, windowHelp)
		return
	}

	config, setup := windowFlags(args)

	path, err := resolveWindowPath(config, setup)

	folder := filepath.Dir(path)
	if err != nil {
		folder = fallbackConfigFolder()
	}

	// The log exists before anything else can fail: a window that cannot start
	// has no console, and "nothing happens" is not a report.
	windowLog := startWindowLog(folder)

	defer func() {
		if recovered := recover(); recovered != nil {
			windowLog.recordPanic(recovered)
			windowLog.keep()
			notifyStartupError(fmt.Sprintf("crower-gui stopped unexpectedly.\n\n%v\n\nDetails: %s", recovered, windowLog.path))
			os.Exit(1)
		}
	}()

	if err != nil {
		reportWindowFailure(windowLog, err)
		return
	}

	cfg, err := session.Load(path, session.Options{Streams: streams})
	if err != nil {
		reportWindowFailure(windowLog, err)
		return
	}

	startWindow(cfg)

	windowLog.clean()
}

// startWindow composes the window's ports and blocks until it closes.
func startWindow(cfg *session.Session) {
	fyneui.Run(cfg.Catalog(), cfg.History, fyneui.Deps{
		Exec:     newExecutor(cfg),
		Editor:   newFyneEditor(cfg),
		Reverter: newFyneReverter(cfg),
		Updater:  newFyneUpdater(cfg, Version),
		Tools:    newFyneTools(cfg),
		Groups:   newFyneGroups(cfg),
		Setups:   newFyneSetups(cfg),
		Window:   newFyneWindow(cfg.ConfigFolder()),
	})
}

// windowFlags reads the flags the window honours. Everything else is ignored on
// purpose: an unexpected argument must not stand between a click and the window.
func windowFlags(args []string) (config, setup string) {
	for index := 0; index < len(args); index++ {
		name, value, hasValue := strings.Cut(args[index], "=")

		switch name {
		case "--config", "-config":
			if hasValue {
				config = value
				continue
			}

			if index+1 < len(args) {
				index++
				config = args[index]
			}
		case "--setup", "-setup":
			if hasValue {
				setup = value
				continue
			}

			if index+1 < len(args) {
				index++
				setup = args[index]
			}
		}
	}

	return config, setup
}

// hasWindowArg reports whether one of the names appears, with or without a
// value.
func hasWindowArg(args []string, names ...string) bool {
	for _, arg := range args {
		name, _, _ := strings.Cut(arg, "=")

		for _, want := range names {
			if name == want {
				return true
			}
		}
	}

	return false
}

// resolveWindowPath resolves the configuration file the same way the CLI does.
func resolveWindowPath(config, setup string) (string, error) {
	if setup == "" {
		setup = os.Getenv("CROWER_SETUP")
	}

	if setup == "" {
		setup = persistedSetup()
	}

	return resolveConfigPath(config, setup)
}

// fallbackConfigFolder is where the log goes when the configuration path itself
// cannot be resolved.
func fallbackConfigFolder() string {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return filepath.Join(dir, "crower")
	}

	return os.TempDir()
}

// reportWindowFailure records a start that never reached the window and says so
// in a dialog naming the log.
func reportWindowFailure(windowLog *windowLog, err error) {
	windowLog.logError(err)
	windowLog.keep()

	message := "crower-gui could not start.\n\n" + err.Error()
	if windowLog.path != "" {
		message += "\n\nLog: " + windowLog.path
	}

	notifyStartupError(message)
}
