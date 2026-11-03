package cmd

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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

// The JSON flag has to be part of the command that can honour it, not a global
// switch a command silently ignores: an unsupported --json is a usage error.
func TestJSONIsDeclaredOnlyWhereSupported(t *testing.T) {
	supported := map[string]bool{"list": true, "search": true}
	seen := 0

	for _, child := range rootCmd.Commands() {
		hasFlag := child.Flags().Lookup(flagJSON) != nil

		if supported[child.Name()] {
			seen++
		}

		if supported[child.Name()] != hasFlag {
			t.Errorf("%s: supports json = %v, has --%s = %v", child.Name(), supported[child.Name()], flagJSON, hasFlag)
		}
	}

	if seen != len(supported) {
		t.Errorf("found %d of the %d commands that support --%s", seen, len(supported), flagJSON)
	}
}

// One JSON object on stdout, nothing on stderr, and the command is still there.
func TestListJSON(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	stdout, stderr, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	if stderr != "" {
		t.Errorf("list --json wrote diagnostics: %q", stderr)
	}

	if strings.Count(stdout, "\n") != 1 {
		t.Fatalf("result is not one line: %q", stdout)
	}

	var payload struct {
		Commands []struct {
			Row     int      `json:"row"`
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
			Exec    string   `json:"exec"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if len(payload.Commands) != 1 {
		t.Fatalf("commands = %+v, want the created command", payload.Commands)
	}

	command := payload.Commands[0]
	if command.Name != "greet" || command.Exec != "echo hello" || command.Row != 0 {
		t.Errorf("command = %+v", command)
	}

	if command.Aliases == nil {
		t.Error("aliases encoded as null")
	}
}

// An empty store still has to be a collection, not null.
func TestListJSONEmptyStore(t *testing.T) {
	stdout, _, err := run(t, "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	if strings.TrimSpace(stdout) != `{"commands":[]}` {
		t.Errorf("empty store = %q", stdout)
	}
}

func TestListHistoryJSON(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--history", "--json")
	if err != nil {
		t.Fatalf("list --history --json returned an error: %v", err)
	}

	var payload struct {
		History []struct {
			Version   int    `json:"version"`
			Timestamp string `json:"timestamp"`
			Operation string `json:"operation"`
			Command   string `json:"command"`
			Note      string `json:"note"`
		} `json:"history"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if len(payload.History) == 0 {
		t.Fatal("history is empty after a create")
	}

	// Newest first: the create is the only entry when this list is rendered,
	// so it must be first. Accepting "list" here would hide an off-by-one.
	newest := payload.History[0]
	if newest.Operation != "create" {
		t.Errorf("newest operation = %q, want create", newest.Operation)
	}

	if newest.Timestamp == "" || newest.Note == "" {
		t.Errorf("newest entry = %+v, want a timestamp and a note", newest)
	}

	if newest.Command != "greet" {
		t.Errorf("newest entry command = %q, want greet", newest.Command)
	}
}

func TestListSysPathJSON(t *testing.T) {
	stdout, _, err := run(t, "list", "--syspath", "--json")
	if err != nil {
		t.Fatalf("list --syspath --json returned an error: %v", err)
	}

	var payload struct {
		Variable struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
		} `json:"variable"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v", err)
	}

	if payload.Variable.Name != "PATH" {
		t.Errorf("variable name = %q, want PATH", payload.Variable.Name)
	}

	if len(payload.Variable.Values) == 0 {
		t.Error("PATH has no entries")
	}
}

// The system listing is environment values, so this only asserts the shape --
// never print the payload, and never copy the values into a fixture.
func TestListSystemJSON(t *testing.T) {
	stdout, _, err := run(t, "list", "--system", "--json")
	if err != nil {
		t.Fatalf("list --system --json returned an error: %v", err)
	}

	var payload struct {
		Variables []struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
		} `json:"variables"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v", err)
	}

	if len(payload.Variables) == 0 {
		t.Fatal("the system listing carried no variables")
	}

	foundPath := false
	for _, variable := range payload.Variables {
		if variable.Name == "PATH" {
			foundPath = true
		}
	}

	if !foundPath {
		t.Error("the system listing has no PATH entry")
	}
}

func TestListProcessJSON(t *testing.T) {
	stdout, _, err := run(t, "list", "--process", "--json")
	if err != nil {
		t.Fatalf("list --process --json returned an error: %v", err)
	}

	var payload struct {
		Processes []struct {
			PID  int32  `json:"pid"`
			Name string `json:"name"`
			Path string `json:"path"`
		} `json:"processes"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v", err)
	}

	if len(payload.Processes) == 0 {
		t.Fatal("the process listing is empty")
	}

	for _, process := range payload.Processes {
		if process.PID <= 0 || process.Name == "" {
			t.Errorf("process row = %+v", process)
		}
	}
}

func TestSearchFolderJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the search commands are shell-backed and this suite runs them on POSIX")
	}

	dir := t.TempDir()

	if err := os.Mkdir(filepath.Join(dir, "needle-dir"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	stdout, _, err := run(t, "search", "--folder", "needle-dir", dir, "--json")
	if err != nil {
		t.Fatalf("search --folder --json returned an error: %v", err)
	}

	var payload struct {
		Search struct {
			Kind    string   `json:"kind"`
			Results []string `json:"results"`
		} `json:"search"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if payload.Search.Kind != "folder" {
		t.Errorf("kind = %q, want folder", payload.Search.Kind)
	}

	if len(payload.Search.Results) != 1 || !strings.HasSuffix(payload.Search.Results[0], "needle-dir") {
		t.Errorf("results = %v", payload.Search.Results)
	}
}

func TestSearchDefaultJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the search commands are shell-backed and this suite runs them on POSIX")
	}

	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "needle.txt"), nil, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	stdout, _, err := run(t, "search", "needle.txt", dir, "--json")
	if err != nil {
		t.Fatalf("search --json returned an error: %v", err)
	}

	var payload struct {
		Search struct {
			Kind    string   `json:"kind"`
			Results []string `json:"results"`
		} `json:"search"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if payload.Search.Kind != "file_and_folder" {
		t.Errorf("kind = %q, want file_and_folder", payload.Search.Kind)
	}

	if len(payload.Search.Results) != 1 {
		t.Errorf("results = %v", payload.Search.Results)
	}
}

func TestListFolderJSON(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("abc"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	stdout, _, err := run(t, "list", "--folder", dir, "--json")
	if err != nil {
		t.Fatalf("list --folder --json returned an error: %v", err)
	}

	var payload struct {
		Folder  string `json:"folder"`
		Entries []struct {
			Name string `json:"name"`
			Dir  bool   `json:"dir"`
			Size int64  `json:"size"`
		} `json:"entries"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if payload.Folder != dir {
		t.Errorf("folder = %q, want %q", payload.Folder, dir)
	}

	if len(payload.Entries) != 1 || payload.Entries[0].Name != "a.txt" || payload.Entries[0].Dir || payload.Entries[0].Size != 3 {
		t.Errorf("entries = %+v", payload.Entries)
	}
}

// A folder that does not exist is a failure in JSON too: a result must not be
// an empty object that reads as "no files".
func TestListFolderJSONMissingFolderFails(t *testing.T) {
	stdout, _, err := run(t, "list", "--folder", filepath.Join(t.TempDir(), "nope"), "--json")
	if err == nil {
		t.Fatal("list --folder on a missing folder returned no error")
	}

	if stdout != "" {
		t.Errorf("a failed listing wrote a result: %q", stdout)
	}
}

func TestSearchJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the search commands are shell-backed and this suite runs them on POSIX")
	}

	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "needle.txt"), nil, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	stdout, _, err := run(t, "search", "--file", "needle.txt", dir, "--json")
	if err != nil {
		t.Fatalf("search --json returned an error: %v", err)
	}

	var payload struct {
		Search struct {
			Kind    string   `json:"kind"`
			Query   string   `json:"query"`
			Path    string   `json:"path"`
			Results []string `json:"results"`
		} `json:"search"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if payload.Search.Kind != "file" || payload.Search.Query != "needle.txt" {
		t.Errorf("search = %+v", payload.Search)
	}

	if len(payload.Search.Results) != 1 || !strings.HasSuffix(payload.Search.Results[0], "needle.txt") {
		t.Errorf("results = %v", payload.Search.Results)
	}
}

// A command without a JSON rendering must reject the flag: silently returning
// the text would be the same class of failure as exiting 0 without doing the
// work.
func TestJSONOnARawCommandIsRejected(t *testing.T) {
	stdout, _, err := run(t, "execute", "--json", "echo", "hello")
	if err == nil {
		t.Fatal("execute --json returned no error")
	}

	if !strings.Contains(err.Error(), "unknown flag") {
		t.Errorf("error = %v, want an unknown flag", err)
	}

	if stdout != "" {
		t.Errorf("the rejected run wrote a result: %q", stdout)
	}
}

// The "last operation" family acts on the command recorded with the operation.
// It used to look up an empty name, so `execute --last` always failed with
// `command not found: ""`.
func TestExecuteLastRunsTheLastExecutedCommand(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "execute", "greet"); err != nil {
		t.Fatalf("execute returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "execute", "--last")
	if err != nil {
		t.Fatalf("execute --last returned an error: %v", err)
	}

	if !strings.Contains(stdout, "hello") {
		t.Errorf("execute --last wrote %q, want the command's output", stdout)
	}
}

// create --execute re-creates the line the last execute ran under a new name.
// The reader looked for a bare "exec" note field that no note builder ever
// wrote, so the command failed every time with "carries no exec to reuse".
func TestCreateExecuteReusesTheLastExecutedCommand(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "execute", "greet"); err != nil {
		t.Fatalf("execute returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "create", "--execute", "greet2"); err != nil {
		t.Fatalf("create --execute returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name string `json:"name"`
			Exec string `json:"exec"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	found := false
	for _, command := range payload.Commands {
		if command.Name == "greet2" {
			found = true
			if command.Exec != "echo hello" {
				t.Errorf("greet2 exec = %q, want the last executed line", command.Exec)
			}
		}
	}

	if !found {
		t.Fatalf("create --execute did not add greet2: %q", stdout)
	}

	stdout, _, err = runWith(t, cfgPath, "", "execute", "greet2")
	if err != nil {
		t.Fatalf("execute greet2 returned an error: %v", err)
	}

	if !strings.Contains(stdout, "hello") {
		t.Errorf("execute greet2 wrote %q, want the reused command's output", stdout)
	}
}

// create --execute takes one name and the flags any other create form takes.
// The aliases were hard-coded to an empty slice, so --alias succeeded and did
// nothing, and a trailing word was ignored instead of rejected.
func TestCreateExecuteHonoursAliasAndRejectsExtraArguments(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	// Execute a *different* command, so the reused line cannot be mistaken for
	// greet's when the alias is resolved below.
	if _, _, err := runWith(t, cfgPath, "", "create", "other", "echo other"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "execute", "other"); err != nil {
		t.Fatalf("execute returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "create", "--execute", "greet2", "-a", "g2"); err != nil {
		t.Fatalf("create --execute with an alias returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	hasAlias := false
	for _, command := range payload.Commands {
		if command.Name == "greet2" {
			for _, alias := range command.Aliases {
				if alias == "g2" {
					hasAlias = true
				}
			}
		}
	}

	if !hasAlias {
		t.Errorf("create --execute dropped the alias: %q", stdout)
	}

	stdout, _, err = runWith(t, cfgPath, "", "execute", "g2")
	if err != nil {
		t.Fatalf("execute by alias returned an error: %v", err)
	}

	if !strings.Contains(stdout, "other") {
		t.Errorf("execute g2 wrote %q, want the reused command's output", stdout)
	}

	_, _, err = runWith(t, cfgPath, "", "create", "--execute", "greet3", "extra")
	if err == nil {
		t.Fatal("create --execute with a trailing argument returned no error")
	}

	if !strings.Contains(err.Error(), "too many arguments") {
		t.Errorf("error = %v, want a too-many-arguments error", err)
	}

	stdout, _, err = runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	if strings.Contains(stdout, "greet3") {
		t.Errorf("a rejected create --execute still added the command: %q", stdout)
	}

	stdout, _, err = runWith(t, cfgPath, "", "list", "--history", "--json")
	if err != nil {
		t.Fatalf("list --history --json returned an error: %v", err)
	}

	var history struct {
		History []struct {
			Command string `json:"command"`
		} `json:"history"`
	}

	if err := json.Unmarshal([]byte(stdout), &history); err != nil {
		t.Fatalf("history is not valid json: %v (%q)", err, stdout)
	}

	for _, entry := range history.History {
		if entry.Command == "greet3" {
			t.Errorf("a rejected create --execute still wrote a history entry: %q", stdout)
		}
	}
}

// An exec may contain the note separator. It is escaped on write, decoded on
// read, so create --execute copies the whole line instead of the part before
// the first ";".
func TestCreateExecuteKeepsASemicolonInTheExec(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")
	exec := "echo one; echo two"

	if _, _, err := runWith(t, cfgPath, "", "create", "work", exec); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "execute", "work"); err != nil {
		t.Fatalf("execute returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "create", "--execute", "work2"); err != nil {
		t.Fatalf("create --execute returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name string `json:"name"`
			Exec string `json:"exec"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	found := false
	for _, command := range payload.Commands {
		if command.Name == "work2" {
			found = true
			if command.Exec != exec {
				t.Errorf("work2 exec = %q, want %q", command.Exec, exec)
			}
		}
	}

	if !found {
		t.Fatalf("create --execute did not add work2: %q", stdout)
	}

	// The stored note is escaped, but the history listing shows the text the
	// user wrote, not the escape sequences.
	stdout, _, err = runWith(t, cfgPath, "", "list", "--history", "--json")
	if err != nil {
		t.Fatalf("list --history --json returned an error: %v", err)
	}

	var history struct {
		History []struct {
			Command string `json:"command"`
			Note    string `json:"note"`
		} `json:"history"`
	}

	if err := json.Unmarshal([]byte(stdout), &history); err != nil {
		t.Fatalf("history is not valid json: %v (%q)", err, stdout)
	}

	shown := false
	for _, entry := range history.History {
		if entry.Command == "work2" && strings.Contains(entry.Note, "command_exec="+exec) {
			shown = true
		}
	}

	if !shown {
		t.Errorf("history does not show the decoded exec: %q", stdout)
	}

	if runtime.GOOS == "windows" {
		return
	}

	stdout, _, err = runWith(t, cfgPath, "", "execute", "work2")
	if err != nil {
		t.Fatalf("execute work2 returned an error: %v", err)
	}

	if !strings.Contains(stdout, "two") {
		t.Errorf("execute work2 wrote %q, want both commands to run", stdout)
	}
}

// create -p names the new command and takes the process as an argument. With
// no process it indexed an empty slice and panicked; it must fail like any
// other missing argument.
func TestCreateProcessWithoutAProcessFails(t *testing.T) {
	_, _, err := run(t, "create", "-p", "only-a-name")
	if err == nil {
		t.Fatal("create -p with no process returned no error")
	}

	if !strings.Contains(err.Error(), "process name or pid") {
		t.Errorf("error = %v, want it to name the missing process", err)
	}
}

// The process form takes --alias like every other create form. It used to
// register the command with no aliases while accepting the flag.
func TestCreateProcessHonoursAlias(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	// A pid that is certainly running: this test process.
	pid := strconv.Itoa(os.Getpid())

	if _, _, err := runWith(t, cfgPath, "", "create", "-p", "probe", pid, "-a", "p1"); err != nil {
		t.Fatalf("create -p returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	hasAlias := false
	for _, command := range payload.Commands {
		if command.Name == "probe" {
			for _, alias := range command.Aliases {
				if alias == "p1" {
					hasAlias = true
				}
			}
		}
	}

	if !hasAlias {
		t.Errorf("create -p dropped the alias: %q", stdout)
	}
}

// The script form takes --alias too. The opener is stubbed in PATH so the test
// never launches a desktop application, but the whole command path runs.
func TestCreateScriptHonoursAlias(t *testing.T) {
	opener := map[string]string{"linux": "xdg-open", "darwin": "open"}[runtime.GOOS]
	if opener == "" {
		t.Skip("the script form's opener has no stub on this platform")
	}

	binDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(binDir, opener), []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatalf("write opener stub: %v", err)
	}

	t.Setenv("PATH", binDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "--script", "probe-script", "-a", "s1"); err != nil {
		t.Fatalf("create --script returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	hasAlias := false
	for _, command := range payload.Commands {
		if command.Name == "probe-script" {
			for _, alias := range command.Aliases {
				if alias == "s1" {
					hasAlias = true
				}
			}
		}
	}

	if !hasAlias {
		t.Errorf("create --script dropped the alias: %q", stdout)
	}
}

// A new command may not take a name that is already an alias: lookups try the
// name map first, so the command the alias pointed at would become unreachable.
func TestCreateNameThatIsAnAliasFails(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "holder", "-a", "g2", "echo two"); err != nil {
		t.Fatalf("create holder returned an error: %v", err)
	}

	_, _, err := runWith(t, cfgPath, "", "create", "g2", "echo shadow")
	if err == nil {
		t.Fatal("create accepted a name that is already an alias")
	}

	if !strings.Contains(err.Error(), `already an alias of "holder"`) {
		t.Errorf("error = %v, want it to name the alias owner", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "execute", "g2")
	if err != nil {
		t.Fatalf("execute g2 returned an error: %v", err)
	}

	if !strings.Contains(stdout, "two") {
		t.Errorf("alias g2 no longer resolves to holder: %q", stdout)
	}

	stdout, _, err = runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name string `json:"name"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	if len(payload.Commands) == 0 {
		t.Fatal("list --json returned no commands, so the check below proves nothing")
	}

	hasHolder := false
	for _, command := range payload.Commands {
		if command.Name == "g2" {
			t.Errorf("the rejected name was still registered: %q", stdout)
		}

		if command.Name == "holder" {
			hasHolder = true
		}
	}

	if !hasHolder {
		t.Errorf("holder is missing, so the rejection cannot be about the alias: %q", stdout)
	}
}

// update must not create the shadow that create refuses: a new name that is an
// alias would win over it, and a new alias that is a name would never be
// reachable at all, because lookups try the name map first.
func TestUpdateCannotShadowANameOrAlias(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "holder", "-a", "g2", "echo two"); err != nil {
		t.Fatalf("create holder returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "create", "other", "echo other"); err != nil {
		t.Fatalf("create other returned an error: %v", err)
	}

	_, _, err := runWith(t, cfgPath, "y\n", "update", "other", "--name", "g2", "-a", "o1", "--exec", "echo other")
	if err == nil {
		t.Fatal("update accepted a name that is already an alias")
	}

	if !strings.Contains(err.Error(), `already an alias of "holder"`) {
		t.Errorf("error = %v, want it to name the alias owner", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "execute", "g2")
	if err != nil {
		t.Fatalf("execute g2 returned an error: %v", err)
	}

	if !strings.Contains(stdout, "two") {
		t.Errorf("alias g2 no longer resolves to holder: %q", stdout)
	}

	_, _, err = runWith(t, cfgPath, "y\n", "update", "other", "--name", "other", "-a", "holder", "--exec", "echo other")
	if err == nil {
		t.Fatal("update accepted an alias that is already a command name")
	}

	if !strings.Contains(err.Error(), "already a command name") {
		t.Errorf("error = %v, want it to name the conflict", err)
	}

	// A name-first lookup finds holder whether or not the alias was accepted,
	// so inspect other's alias list: the rejected alias must not be there.
	stdout, _, err = runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	found := false
	for _, command := range payload.Commands {
		if command.Name != "other" {
			continue
		}

		found = true
		for _, alias := range command.Aliases {
			if alias == "holder" {
				t.Errorf("the rejected alias was still registered: %q", stdout)
			}
		}
	}

	if !found {
		t.Errorf("other is missing, so the check above proves nothing: %q", stdout)
	}
}

// update takes a name or an alias, as its help says. The input path used to
// check only the name map, so every alias key was rejected before the handler
// could resolve it.
func TestUpdateByAlias(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "holder", "-a", "g2", "echo two"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "y\n", "update", "g2", "--name", "holder", "-a", "g2", "--exec", "echo changed"); err != nil {
		t.Fatalf("update by alias returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "execute", "g2")
	if err != nil {
		t.Fatalf("execute g2 returned an error: %v", err)
	}

	if !strings.Contains(stdout, "changed") {
		t.Errorf("update by alias did not change the command: %q", stdout)
	}

	stdout, _, err = runWith(t, cfgPath, "", "execute", "holder")
	if err != nil {
		t.Fatalf("execute holder returned an error: %v", err)
	}

	if !strings.Contains(stdout, "changed") {
		t.Errorf("the name stopped resolving after an alias-keyed update: %q", stdout)
	}

	_, _, err = runWith(t, cfgPath, "", "update", "no-such-command", "--name", "x", "-a", "y", "--exec", "z")
	if err == nil {
		t.Fatal("update of an unknown key returned no error")
	}

	if !strings.Contains(err.Error(), "command not found") {
		t.Errorf("error = %v, want a command-not-found error", err)
	}
}

// delete --create removes the command the last create added. With the empty
// name it failed before touching anything.
func TestDeleteCreatedLastDeletesTheCreatedCommand(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "delete", "--create"); err != nil {
		t.Fatalf("delete --create returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list")
	if err != nil {
		t.Fatalf("list returned an error: %v", err)
	}

	if strings.Contains(stdout, "greet") {
		t.Errorf("delete --create left the command behind: %q", stdout)
	}
}

func writeTestTar(t *testing.T, archivePath string) {
	t.Helper()

	file, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("create %s: %v", archivePath, err)
	}

	writer := tar.NewWriter(file)

	header := &tar.Header{Name: "x.txt", Mode: 0o644, Size: 1}
	if err := writer.WriteHeader(header); err != nil {
		t.Fatalf("write header: %v", err)
	}

	if _, err := writer.Write([]byte("x")); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("close %s: %v", archivePath, err)
	}
}

// `extract` records the extract operation, not an execute: it used to be filed
// as Execute with no command, so `execute --last` after an extract ran nothing
// and failed with `command not found: ""`.
func TestExtractDoesNotCountAsTheLastExecute(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("no gate runs a Windows shell: extraction there is unverified")
	}

	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	if _, _, err := runWith(t, cfgPath, "", "execute", "greet"); err != nil {
		t.Fatalf("execute returned an error: %v", err)
	}

	archive := filepath.Join(t.TempDir(), "x.tar")
	writeTestTar(t, archive)

	if _, _, err := runWith(t, cfgPath, "", "extract", archive); err != nil {
		t.Fatalf("extract returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "execute", "--last")
	if err != nil {
		t.Fatalf("execute --last after an extract returned an error: %v", err)
	}

	if !strings.Contains(stdout, "hello") {
		t.Errorf("execute --last wrote %q, want the last executed command's output", stdout)
	}
}

// A name beginning with "-" is creatable with the terminator; the recorded
// crower_exec uses the same shape, so this is the line that has to work.
func TestCreateNameBeginningWithADash(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "--", "-v", "echo x"); err != nil {
		t.Fatalf("create -- -v returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--json")
	if err != nil {
		t.Fatalf("list --json returned an error: %v", err)
	}

	var payload struct {
		Commands []struct {
			Name string `json:"name"`
			Exec string `json:"exec"`
		} `json:"commands"`
	}

	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		t.Fatalf("result is not valid json: %v (%q)", err, stdout)
	}

	found := false
	for _, command := range payload.Commands {
		if command.Name == "-v" && command.Exec == "echo x" {
			found = true
		}
	}

	if !found {
		t.Errorf("create -- -v did not add the command: %q", stdout)
	}
}

// The history table names the command column; the JSON already carried it.
func TestListHistoryShowsTheCommandColumn(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create returned an error: %v", err)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list", "--history")
	if err != nil {
		t.Fatalf("list --history returned an error: %v", err)
	}

	if !strings.Contains(stdout, "Command") {
		t.Errorf("history has no Command column: %q", stdout)
	}
}

// reset shows what it will erase, asks, and on a yes leaves no command behind.
func TestResetErasesEveryCommand(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	for _, name := range []string{"greet", "build"} {
		if _, _, err := runWith(t, cfgPath, "", "create", name, "echo "+name); err != nil {
			t.Fatalf("create %s: %v", name, err)
		}
	}

	stdout, stderr, err := runWith(t, cfgPath, "y\n", "reset")
	if err != nil {
		t.Fatalf("reset returned an error: %v", err)
	}

	if stdout != "" {
		t.Errorf("reset wrote a result: %q", stdout)
	}

	for _, want := range []string{"Name", "greet", "All commands will be erased...", "[Y]es/[N]o", "reset all commands"} {
		if !strings.Contains(stderr, want) {
			t.Errorf("the diagnostics stream does not contain %q: %q", want, stderr)
		}
	}

	stdout, _, err = runWith(t, cfgPath, "", "list")
	if err != nil {
		t.Fatalf("list returned an error: %v", err)
	}

	for _, name := range []string{"greet", "build"} {
		if strings.Contains(stdout, name) {
			t.Errorf("list still shows %q after a reset:\n%s", name, stdout)
		}
	}
}

// Declining keeps every command and is not a failure.
func TestDecliningResetKeepsTheCommands(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo hello"); err != nil {
		t.Fatalf("create: %v", err)
	}

	_, stderr, err := runWith(t, cfgPath, "n\n", "reset")
	if err != nil {
		t.Fatalf("a declined reset returned an error: %v", err)
	}

	if !strings.Contains(stderr, "Cancelling reset...") {
		t.Errorf("the decline was not reported: %q", stderr)
	}

	stdout, _, err := runWith(t, cfgPath, "", "list")
	if err != nil {
		t.Fatalf("list returned an error: %v", err)
	}

	if !strings.Contains(stdout, "greet") {
		t.Errorf("the command was erased by a decline:\n%s", stdout)
	}
}

// --no-input must fail the question rather than wait for an answer.
func TestResetNoInputFails(t *testing.T) {
	if _, _, err := run(t, "reset", "--no-input"); err == nil {
		t.Fatal("reset --no-input returned no error, so it would wait")
	}
}

// restore has to write the command it restores to the configuration file, not
// only to the map this process holds: Save() writes the history snapshot and
// index, and nothing else writes the command file (#247).
func TestRestoreWritesTheCommandFile(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	if _, _, err := runWith(t, cfgPath, "", "create", "greet", "echo one"); err != nil {
		t.Fatalf("create: %v", err)
	}

	// update --exec still offers the other fields, so answer them.
	if _, _, err := runWith(t, cfgPath, "n\nn\ny\n", "update", "greet", "--exec", "echo two"); err != nil {
		t.Fatalf("update: %v", err)
	}

	if got := configExec(t, cfgPath); got != "echo two" {
		t.Fatalf("exec after update = %q, want echo two", got)
	}

	// The history lists the newest entry first: row 1 is the create.
	stdout, stderr, err := runWith(t, cfgPath, "1\ny\n", "restore", "greet")
	if err != nil {
		t.Fatalf("restore: %v", err)
	}

	if stdout != "" {
		t.Errorf("restore wrote a result: %q", stdout)
	}

	if !strings.Contains(stderr, "restored command:") {
		t.Errorf("restore reported nothing: %q", stderr)
	}

	if got := configExec(t, cfgPath); got != "echo one" {
		t.Errorf("exec after restore = %q, want the restored value echo one", got)
	}

	stdout, _, err = runWith(t, cfgPath, "", "list")
	if err != nil {
		t.Fatalf("list: %v", err)
	}

	if !strings.Contains(stdout, "echo one") {
		t.Errorf("list does not show the restored command:\n%s", stdout)
	}
}

// configExec reads the exec of the one command in a test configuration.
func configExec(t *testing.T, path string) string {
	t.Helper()

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}

	for _, line := range strings.Split(string(body), "\n") {
		if value, ok := strings.CutPrefix(strings.TrimSpace(line), "Exec = "); ok {
			return strings.Trim(value, `"`)
		}
	}

	return ""
}
