package cmd

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// runWith executes the root command against buffers instead of a terminal, and
// returns what reached each stream.
//
// The configuration path is always explicit: a command that fell back to $HOME
// would write into the machine running the test.
func runWith(t *testing.T, cfgPath string, stdin string, args ...string) (stdout string, stderr string, err error) {
	t.Helper()

	resetFlags()

	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	previous := streams
	streams = app.Streams{In: strings.NewReader(stdin), Out: out, Err: errOut}
	t.Cleanup(func() { streams = previous })

	rootCmd.SetArgs(append([]string{"--config", cfgPath}, args...))
	t.Cleanup(func() { rootCmd.SetArgs(nil) })

	err = rootCmd.Execute()

	return out.String(), errOut.String(), err
}

// resetFlags puts every flag back to its default before a run.
//
// The binary runs once per process, so cobra's package-level flag variables are
// only ever set by the command line. A test runs the command tree many times in
// one process, and a flag one test set would otherwise still be set in the
// next -- `delete --file x` followed by `create a b` would try to create a file
// called "a b".
func resetFlags() {
	var reset func(cmd *cobra.Command)

	reset = func(cmd *cobra.Command) {
		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if err := f.Value.Set(f.DefValue); err != nil {
				panic(err)
			}
			f.Changed = false
		})

		for _, child := range cmd.Commands() {
			reset(child)
		}
	}

	reset(rootCmd)
	resetPersistentFlags(rootCmd.PersistentFlags())

	versionFlag = false
	checkNewVersion = false
	upgradeFlag = false
	executeFlag = false
	fileFlag = false
	folderFlag = false
	sysPathFlag = false
	systemFlag = false
	createFlag = false
	updateFlag = false
	processFlag = false
	historyFlag = false
	historyFlagDelete = false
	last = false
	create = false
	execute = false
	browserFlag = false
	scriptFlag = false
	processName = ""
	name = ""
	exec = ""
	allAlias = nil
	outDirFlag = ""
}

// resetPersistentFlags resets the flags every command inherits. They live on
// the root command's persistent set, which is separate from its own flags.
func resetPersistentFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		if err := f.Value.Set(f.DefValue); err != nil {
			panic(err)
		}
		f.Changed = false
	})
}

// run is runWith with a fresh configuration in a fresh directory.
func run(t *testing.T, args ...string) (string, string, error) {
	t.Helper()

	return runWith(t, filepath.Join(t.TempDir(), "crower.yaml"), "", args...)
}

// The whole point of the change: a failure has to be visible to a script. The
// command used to return no error, so the process exited 0 and the script went
// on as if the work had happened.
func TestFailuresReturnAnError(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"unknown command to execute", []string{"execute", "no-such-command-xyz"}},
		{"unknown command to delete", []string{"delete", "--no-input", "no-such-command-xyz"}},
		{"copy without arguments", []string{"copy"}},
		{"move without arguments", []string{"move"}},
		{"extract without arguments", []string{"extract"}},
		{"search without arguments", []string{"search"}},
		{"open without arguments", []string{"open"}},
		{"delete a file that does not exist", []string{"delete", "--file", filepath.Join(t.TempDir(), "missing-xyz")}},
		{"execute --last with no history", []string{"execute", "--last"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := run(t, tt.args...)
			if err == nil {
				t.Fatalf("%v returned no error, so the process would exit 0", tt.args)
			}
		})
	}
}

// A prompt with no input available is an error, not a wait. Before this,
// `crower create </dev/null` re-prompted until it was killed.
func TestNoInputFailsInsteadOfWaiting(t *testing.T) {
	_, _, err := run(t, "create", "--no-input")
	if err == nil {
		t.Fatal("create --no-input returned no error")
	}

	if !strings.Contains(err.Error(), "input") {
		t.Errorf("error does not say that input is required: %v", err)
	}
}

// A closed input stream is the same failure as --no-input: nobody is there to
// answer, so the run must not continue with an empty answer.
func TestClosedInputFails(t *testing.T) {
	_, _, err := runWith(t, filepath.Join(t.TempDir(), "crower.yaml"), "", "create")
	if err == nil {
		t.Fatal("create with a closed stdin returned no error")
	}
}

// --version answers without reading or creating the configuration, so it works
// on a machine where crower has never run.
func TestVersionDoesNotNeedAConfiguration(t *testing.T) {
	out, _, err := run(t, "--version")
	if err != nil {
		t.Fatalf("--version returned an error: %v", err)
	}

	if !strings.Contains(out, Version) {
		t.Errorf("--version wrote %q, want it to contain %q", out, Version)
	}
}

// Results go to one stream and diagnostics to the other, so `crower list | wc
// -l` counts rows.
func TestResultsAndDiagnosticsAreSplit(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	stdout, stderr, err := runWith(t, cfgPath, "", "create", "greet", "echo hello")
	if err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if stdout != "" {
		t.Errorf("create wrote a result: %q", stdout)
	}

	stdout, stderr, err = runWith(t, cfgPath, "", "list")
	if err != nil {
		t.Fatalf("list returned an error: %v", err)
	}

	if !strings.Contains(stdout, "greet") {
		t.Errorf("list did not write the command to the result stream: %q", stdout)
	}

	if strings.Contains(stderr, "greet") {
		t.Errorf("list wrote the command to the diagnostics stream: %q", stderr)
	}
}

// Every list line used to end with a space, because the table was printed
// through a formatter that appended one.
func TestListLinesHaveNoTrailingSpace(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list")
	if err != nil {
		t.Fatalf("list returned an error: %v", err)
	}

	for _, line := range strings.Split(strings.TrimRight(stdout, "\n"), "\n") {
		if strings.HasSuffix(line, " ") {
			t.Errorf("list line ends with a space: %q", line)
		}
	}
}

// A user who declines is not a failure: nothing happened, and that was the
// answer. The run exits 0.
func TestDecliningAConfirmationSucceeds(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	_, stderr, err := runWith(t, cfgPath, "n\n", "delete", "greet")
	if err != nil {
		t.Fatalf("declining returned an error: %v", err)
	}

	if !strings.Contains(stderr, "Cancelling") {
		t.Errorf("the decline was not reported: %q", stderr)
	}
}

// A failure must not print the whole usage text to the result stream: the
// script piping the results would receive it as data.
func TestMissingArgumentsPrintUsageToTheDiagnosticsStream(t *testing.T) {
	stdout, stderr, err := run(t, "copy")
	if err == nil {
		t.Fatal("copy returned no error")
	}

	if stdout != "" {
		t.Errorf("usage reached the result stream: %q", stdout)
	}

	if !strings.Contains(stderr, "Usage:") {
		t.Errorf("usage did not reach the diagnostics stream: %q", stderr)
	}
}

// A prompt screen prints the command table, and a prompt is a diagnostic. The
// table used to go to the result stream, so a failed `crower delete nope`
// printed a table on stdout that a caller piping the results read as data.
func TestThePromptTableGoesToTheDiagnosticsStream(t *testing.T) {
	stdout, stderr, err := run(t, "delete", "--no-input", "no-such-command-xyz")
	if err == nil {
		t.Fatal("delete of an unknown command returned no error")
	}

	if stdout != "" {
		t.Errorf("the prompt table reached the result stream: %q", stdout)
	}

	if !strings.Contains(stderr, "Name") {
		t.Errorf("the prompt table did not reach the diagnostics stream: %q", stderr)
	}
}

// The two flags that make a run scriptable are persistent, so every command
// accepts them without declaring them again.
func TestGlobalFlagsArePersistent(t *testing.T) {
	for _, name := range []string{flagConfig, flagNoInput} {
		if rootCmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("root command has no persistent --%s flag", name)
		}

		for _, child := range rootCmd.Commands() {
			if child.Flags().Lookup(name) == nil && child.InheritedFlags().Lookup(name) == nil {
				t.Errorf("%s does not inherit --%s", child.Name(), name)
			}
		}
	}
}
