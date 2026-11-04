package session

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
	"github.com/fiwon123/crower/internal/domain/operation"
	"github.com/fiwon123/crower/internal/domain/timeline"

	"github.com/fiwon123/crower/pkg/utils"
)

// These tests pin the on-disk contract, not the in-memory model.
//
// Everything here is already written into users' home directories: the command
// map in `crower.yaml` (TOML despite the name), the history index in
// `historyhelper.json`, and one `history/NNNNN.yaml` snapshot per operation.
// A field rename that does not keep its tag orphans that data silently, so the
// refactor may move these types anywhere as long as these bytes stay readable.

func TestCommandMapTOMLRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")

	in := command.NewMapData()
	in.Add("greet", command.New("greet", []string{"hi", "hello"}, "echo hello"))
	in.Add("build", &command.Data{
		Name:     "build",
		AllAlias: []string{},
		Exec:     "go build ./...",
		Steps:    []string{"go vet ./...", "go test ./..."},
	})

	if err := utils.WriteToml(in, path); err != nil {
		t.Fatalf("WriteToml: %v", err)
	}

	out := command.NewMapData()
	if err := utils.ReadToml(path, &out); err != nil {
		t.Fatalf("ReadToml: %v", err)
	}

	if !reflect.DeepEqual(out, in) {
		t.Fatalf("round trip changed the map:\n got %#v\nwant %#v", out, in)
	}
}

// Steps carries `omitempty` on both tags. A command without steps must not
// grow an empty `Steps` key, and one with steps must keep them.
func TestStepsAreOmittedWhenAbsent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")

	in := command.NewMapData()
	in.Add("plain", command.New("plain", nil, "echo hi"))
	in.Add("multi", &command.Data{Name: "multi", AllAlias: []string{}, Exec: "", Steps: []string{"a", "b"}})

	if err := utils.WriteToml(in, path); err != nil {
		t.Fatalf("WriteToml: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written config: %v", err)
	}

	text := string(body)

	plain := tableSection(t, text, "plain")
	if strings.Contains(plain, "Steps") {
		t.Errorf("a command without steps wrote a Steps key:\n%s", plain)
	}

	multi := tableSection(t, text, "multi")
	if !strings.Contains(multi, "Steps") {
		t.Errorf("a command with steps lost them on disk:\n%s", multi)
	}
}

// tableSection returns the text of one TOML table, from its header to the next
// header. It keeps the assertions above about the shape of the file rather than
// about the whole document.
func tableSection(t *testing.T, document, name string) string {
	t.Helper()

	header := "[" + name + "]"
	start := strings.Index(document, header)
	if start < 0 {
		t.Fatalf("no %s table in:\n%s", header, document)
	}

	rest := document[start+len(header):]

	if next := strings.Index(rest, "\n["); next >= 0 {
		rest = rest[:next]
	}

	return rest
}

// The configuration keeps the order its file has: ReadKeysTomlInOrder is what
// makes `crower list` show commands in the order the user saved them rather
// than in map order. The encoder writes the file, the parser reads it back, and
// the two must agree on the sequence for every key.
func TestReadKeysPreservesFileOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")

	const document = `[zebra]
  Name = "zebra"
  AllAlias = []
  Exec = "echo z"

[alpha]
  Name = "alpha"
  AllAlias = []
  Exec = "echo a"

[middle]
  Name = "middle"
  AllAlias = []
  Exec = "echo m"
`

	if err := os.WriteFile(path, []byte(document), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	keys, err := utils.ReadKeysTomlInOrder(path)
	if err != nil {
		t.Fatalf("ReadKeysTomlInOrder: %v", err)
	}

	if want := []string{"zebra", "alpha", "middle"}; !slices.Equal(keys, want) {
		t.Errorf("keys = %v, want the file order %v", keys, want)
	}
}

// A written map is read back in a stable order: the next load of the same file
// must not reshuffle the listing.
func TestWrittenMapKeepsAStableKeyOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")

	in := command.NewMapData()
	for _, name := range []string{"zebra", "alpha", "middle"} {
		in.Add(name, command.New(name, nil, "echo "+name))
	}

	if err := utils.WriteToml(in, path); err != nil {
		t.Fatalf("WriteToml: %v", err)
	}

	first, err := utils.ReadKeysTomlInOrder(path)
	if err != nil {
		t.Fatalf("ReadKeysTomlInOrder: %v", err)
	}

	second, err := utils.ReadKeysTomlInOrder(path)
	if err != nil {
		t.Fatalf("ReadKeysTomlInOrder: %v", err)
	}

	if !slices.Equal(first, second) {
		t.Errorf("key order changed between reads: %v then %v", first, second)
	}

	sortedFirst := slices.Clone(first)
	slices.Sort(sortedFirst)

	if !slices.Equal(sortedFirst, []string{"alpha", "middle", "zebra"}) {
		t.Errorf("the file does not contain every added key: %v", first)
	}
}

// Timestemp (sic) is the persisted key of the timestamp field. Correcting the
// spelling in the struct would make every existing history file parse with an
// empty date, so the misspelling is part of the file format.
func TestHistoryJSONKeepsItsOnDiskKeys(t *testing.T) {
	path := filepath.Join(t.TempDir(), "historyhelper.json")

	data := timeline.Log{}
	data.Add(operation.Create, "greet", "main_operation=create;command_name=greet")
	data.Add(operation.Copy, "", "main_operation=copy;crower_exec=crower copy a b")

	if err := utils.WriteJson(data, path); err != nil {
		t.Fatalf("WriteJson: %v", err)
	}

	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written history: %v", err)
	}

	for _, key := range []string{
		`"AllData"`, `"HistoryFolderPath"`,
		`"Version"`, `"File"`, `"Timestemp"`,
		`"Operation"`, `"CommandName"`, `"Note"`,
	} {
		if !strings.Contains(string(body), key) {
			t.Errorf("history file lost the key %s:\n%s", key, body)
		}
	}

	var back timeline.Log
	if err := utils.ReadJson(path, &back); err != nil {
		t.Fatalf("ReadJson: %v", err)
	}

	if !reflect.DeepEqual(back, data) {
		t.Fatalf("round trip changed the history:\n got %#v\nwant %#v", back, data)
	}
}

// An empty history file is a fresh install, not a broken one: ReadJson leaves
// the zero value alone and Load continues with no entries.
func TestEmptyHistoryFileIsNotAnError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "historyhelper.json")

	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatalf("write empty file: %v", err)
	}

	var data timeline.Log
	if err := utils.ReadJson(path, &data); err != nil {
		t.Fatalf("ReadJson of an empty file: %v", err)
	}

	if len(data.AllData) != 0 {
		t.Errorf("AllData = %v, want empty", data.AllData)
	}
}

// A history snapshot is the same TOML command map as the configuration, and
// RevertTo reads one and writes it over the other. This is the property that
// lets a snapshot taken by an older build be restored by this one.
func TestHistorySnapshotUsesTheCommandMapFormat(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	cfg := newTestConfigAt(t, cfgPath)
	cfg.allCommandsByName.Add("greet", command.New("greet", []string{"hi"}, "echo hello"))

	if err := cfg.SaveCommands(); err != nil {
		t.Fatalf("SaveCommands: %v", err)
	}

	snapshotPath := filepath.Join(t.TempDir(), "00001.yaml")

	body, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf("read configuration: %v", err)
	}

	if err := os.WriteFile(snapshotPath, body, 0o644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}

	snapshot := command.NewMapData()
	if err := utils.ReadToml(snapshotPath, &snapshot); err != nil {
		t.Fatalf("a configuration file could not be read as a snapshot: %v", err)
	}

	got := snapshot.Get("greet")
	if got == nil {
		t.Fatal("the snapshot did not contain the command")
	}

	if !slices.Equal(got.AllAlias, []string{"hi"}) {
		t.Errorf("snapshot aliases = %v, want [hi]", got.AllAlias)
	}
}
