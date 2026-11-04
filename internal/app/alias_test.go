package app

import (
	"bytes"
	"github.com/fiwon123/crower/internal/domain/timeline"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/storage/command"
)

// Ranging a Go map made the same config run a different command on different
// starts. The build is sorted now: the first name wins.
func TestGetAliasMapIsDeterministic(t *testing.T) {
	commands := command.NewMapData()
	commands.Add("b", command.New("b", []string{"x"}, "echo b"))
	commands.Add("a", command.New("a", []string{"x"}, "echo a"))

	for range 50 {
		owner := getAliasMap(commands).Get("x")
		if owner == nil || owner.Name != "a" {
			t.Fatalf("alias x resolves to %+v, want a", owner)
		}
	}
}

func TestCommandCollisions(t *testing.T) {
	tests := []struct {
		name     string
		commands command.MapData
		want     []string
	}{
		{
			"no collision",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", []string{"x"}, "echo a"))
				m.Add("b", command.New("b", []string{"y"}, "echo b"))
				return m
			}(),
			nil,
		},
		{
			"two commands claim the same alias",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", []string{"x"}, "echo a"))
				m.Add("b", command.New("b", []string{"x"}, "echo b"))
				return m
			}(),
			[]string{`alias "x" is claimed by "a" and "b"; using "a"`},
		},
		{
			"an alias is another command's name",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", nil, "echo a"))
				m.Add("b", command.New("b", []string{"a"}, "echo b"))
				return m
			}(),
			[]string{`alias "a" of "b" is also a command name; the name wins`},
		},
		{
			"a self alias is not a collision",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", []string{"a"}, "echo a"))
				return m
			}(),
			nil,
		},
		{
			"a contested alias that is also a name",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", []string{"x"}, "echo a"))
				m.Add("x", command.New("x", []string{"x"}, "echo x"))
				return m
			}(),
			[]string{
				`alias "x" is claimed by "a" and "x"`,
				`alias "x" of "a" is also a command name; the name wins`,
			},
		},
		{
			"a duplicate alias inside one command warns once",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", []string{"x"}, "echo a"))
				m.Add("b", command.New("b", []string{"x", "x"}, "echo b"))
				return m
			}(),
			[]string{`alias "x" is claimed by "a" and "b"; using "a"`},
		},
		{
			"three claimants",
			func() command.MapData {
				m := command.NewMapData()
				m.Add("a", command.New("a", []string{"x"}, "echo a"))
				m.Add("b", command.New("b", []string{"x"}, "echo b"))
				m.Add("c", command.New("c", []string{"x"}, "echo c"))
				return m
			}(),
			[]string{
				`alias "x" is claimed by "a" and "b"; using "a"`,
				`alias "x" is claimed by "a" and "c"; using "a"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := commandCollisions(tt.commands)
			if !slices.Equal(got, tt.want) {
				t.Errorf("commandCollisions() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Loading a colliding config reports it once and resolves the alias the same
// way every run.
func TestInitAppWarnsAboutACollision(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	contents := `[a]
Name = "a"
AllAlias = ["x"]
Exec = "echo a"

[b]
Name = "b"
AllAlias = ["x"]
Exec = "echo b"
`

	if err := os.WriteFile(cfgPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var stderr bytes.Buffer

	cfg, err := InitApp(cfgPath, Options{
		Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: &stderr},
	})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	if owner := cfg.AllCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Errorf("alias x resolves to %+v, want a", owner)
	}

	if !strings.Contains(stderr.String(), `alias "x" is claimed by "a" and "b"`) {
		t.Errorf("the collision was not reported: %q", stderr.String())
	}
}

func TestInitAppWarnsAboutAnAliasThatIsAName(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	contents := `[a]
Name = "a"
Exec = "echo a"

[b]
Name = "b"
AllAlias = ["a"]
Exec = "echo b"
`

	if err := os.WriteFile(cfgPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var stderr bytes.Buffer

	_, err := InitApp(cfgPath, Options{
		Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: &stderr},
	})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	if !strings.Contains(stderr.String(), `alias "a" of "b" is also a command name`) {
		t.Errorf("the shadowed alias was not reported: %q", stderr.String())
	}
}

// A clean config must not write to the diagnostics stream: `list --json`
// promises an empty stderr, and this warning runs for every command.
func TestInitAppWritesNoWarningForACleanConfig(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	contents := `[a]
Name = "a"
AllAlias = ["x"]
Exec = "echo a"

[b]
Name = "b"
AllAlias = ["y"]
Exec = "echo b"
`

	if err := os.WriteFile(cfgPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	var stderr bytes.Buffer

	if _, err := InitApp(cfgPath, Options{
		Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: &stderr},
	}); err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	if stderr.Len() != 0 {
		t.Errorf("a clean config wrote to stderr: %q", stderr.String())
	}
}

// A snapshot with a collision restores the same command every run, and a name
// wins over an alias.
func TestFindCommandPrefersNameAndIsDeterministic(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")

	cfg, err := InitApp(cfgPath, Options{
		Streams: Streams{In: strings.NewReader(""), Out: io.Discard, Err: io.Discard},
	})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	snapshot := filepath.Join(cfg.ConfigFolder(), "history", "00001.yaml")
	contents := `[a]
Name = "a"
AllAlias = ["x"]
Exec = "echo a"

[b]
Name = "b"
AllAlias = ["x", "a"]
Exec = "echo b"
`

	if err := os.WriteFile(snapshot, []byte(contents), 0o644); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}

	content := timeline.Content{File: "00001.yaml"}

	for range 20 {
		got, err := cfg.FindCommand("a", content)
		if err != nil {
			t.Fatalf("FindCommand(a): %v", err)
		}

		if got.Name != "a" {
			t.Fatalf("FindCommand(a) = %q, want the command named a", got.Name)
		}

		got, err = cfg.FindCommand("x", content)
		if err != nil {
			t.Fatalf("FindCommand(x): %v", err)
		}

		if got.Name != "a" {
			t.Fatalf("FindCommand(x) = %q, want the first name with that alias (a)", got.Name)
		}
	}
}
