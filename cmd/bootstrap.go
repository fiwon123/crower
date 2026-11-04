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
	flagJSON    = "json"
	flagSetup   = "setup"
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

// setupConfigPath is where a named setup keeps its state.
//
// The default setup deliberately keeps the old path: a user who never asks for
// a setup sees no change and needs no migration. Each setup folder carries its
// own history, so reverting in one can never touch another.
func setupConfigPath(home, name string) (string, error) {
	if name == "" || name == "." || name == ".." {
		return "", fmt.Errorf("invalid setup name %q", name)
	}

	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '-', r == '_', r == '.':
		default:
			return "", fmt.Errorf("invalid setup name %q", name)
		}
	}

	return filepath.Join(home, "crower", "setups", name, "crower.yaml"), nil
}

// resolveConfigPath applies the selection order: an explicit --config wins,
// then --setup, then CROWER_SETUP, then the default installation.
func resolveConfigPath(explicit, setup string) (string, error) {
	if explicit != "" {
		return explicit, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.NotUserHome()
	}

	if setup == "" {
		return filepath.Join(home, "crower", "crower.yaml"), nil
	}

	return setupConfigPath(home, setup)
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

	setup, err := cmd.Flags().GetString(flagSetup)
	if err != nil {
		return nil, err
	}

	if setup == "" {
		setup = os.Getenv("CROWER_SETUP")
	}

	cfgPath, err = resolveConfigPath(cfgPath, setup)
	if err != nil {
		return nil, err
	}

	noInput, err := cmd.Flags().GetBool(flagNoInput)
	if err != nil {
		return nil, err
	}

	return app.InitApp(cfgPath, app.Options{
		Streams: streams,
		NoInput: noInput,
		Format:  requestedFormat(cmd),
	})
}

// requestedFormat reads the --json flag the running command declared, if it
// declared one.
//
// --json is registered per command, on the ones whose result has a structured
// shape. A command that never declared it cannot be asked for JSON, so an
// unsupported `--json` is an unknown-flag error rather than a silent fall back
// to text -- exactly the "reported success without doing the work" failure the
// exit-status contract was added to prevent.
func requestedFormat(cmd *cobra.Command) app.Format {
	if cmd.Flags().Lookup(flagJSON) == nil {
		return app.FormatText
	}

	asJSON, err := cmd.Flags().GetBool(flagJSON)
	if err != nil || !asJSON {
		return app.FormatText
	}

	return app.FormatJSON
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
