package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/internal/errors"

	"github.com/spf13/cobra"
)

const (
	flagConfig  = "config"
	flagNoInput = "no-input"
)

// streams is the terminal this run reads and writes. It is a variable so a test
// can capture results and diagnostics separately instead of scraping one pipe.
var streams = app.DefaultStreams()

// defaultConfigPath is where crower keeps its state when --config is not given.
//
// The file keeps its historical .yaml name even though it holds TOML: renaming
// it without a migration would orphan the configuration of every existing user.
func defaultConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.NotUserHome()
	}

	return filepath.Join(home, "crower", "crower.yaml"), nil
}

// newApp builds the application from the global flags of cmd.
//
// Every command goes through here, so --config and --no-input are resolved in
// exactly one place. The commands used to call helper.GetPersistentConfigFlag
// individually and discard its error twenty times over.
//
// It also marks the command as past argument parsing: an operational failure
// must not print the whole help text after "file not found", because that
// buries the one line a script needs.
func newApp(cmd *cobra.Command) (*app.Config, error) {
	cmd.SilenceUsage = true

	cfgPath, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return nil, err
	}

	if cfgPath == "" {
		cfgPath, err = defaultConfigPath()
		if err != nil {
			return nil, err
		}
	}

	noInput, err := cmd.Flags().GetBool(flagNoInput)
	if err != nil {
		return nil, err
	}

	return app.InitApp(cfgPath, app.Options{
		Streams: streams,
		NoInput: noInput,
	})
}

// missingArgs reports a command run without the arguments it needs.
//
// The usage text goes to the diagnostics stream, next to the error, so a script
// that is piping the results does not receive it. The error keeps the exit
// status non-zero and the usage keeps the fix visible.
func missingArgs(cmd *cobra.Command, what string) error {
	fmt.Fprint(streams.Err, cmd.UsageString())
	return errors.NotEnoughArgs(what)
}
