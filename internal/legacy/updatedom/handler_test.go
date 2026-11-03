package updatedom_test

import (
	"bytes"
	"io"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/legacy/createdom"
	"github.com/fiwon123/crower/internal/legacy/updatedom"
	"github.com/fiwon123/crower/internal/testutil"
)

func TestUpdate(t *testing.T) {

	t.Run("Update command name", func(t *testing.T) {

		app, _, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		updateHandler := updatedom.NewHandler(app)

		var mock = []struct {
			name  string
			alias string
			exec  string
		}{
			{"c1", "a1", ""},
			{"c2", "a2", ""},
			{"c3", "a3", ""},
			{"c4", "a4", ""},
			{"c5", "a5", ""},
		}

		for _, command := range mock {
			createHandler.CreateCommand(
				command.name, []string{command.alias}, "exec")
		}

		var tests = []struct {
			oldName string
			newName string
			want    bool
		}{
			{"c1", "c2", false},
			{"c7", "c1", false},
			{"c1", "c6", true},
			{"c2", "c1", true},
			{"c3", "c1", false},
			{"a3", "c6", false},
		}

		for _, test := range tests {
			newCommand := command.New(test.newName, []string{}, "")
			key := test.oldName
			_, _, err := updateHandler.UpdateCommand(key, newCommand.Name, newCommand.AllAlias, newCommand.Exec)
			got := err == nil
			assertUpdateTest(test.want, got, key, *newCommand, err, t)

		}

	})
}

func TestUpdateRejectsShadowing(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		newName  string
		newAlias []string
		want     string
	}{
		{"a name that is an alias", "other", "g2", nil, `already an alias of "holder"`},
		{"an alias that is a name", "other", "", []string{"holder"}, "already a command name"},
		{"an alias that is an alias", "other", "", []string{"g2"}, "alias already in use"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, _, err := testutil.New(t)
			if err != nil {
				t.Fatalf("error before test: %v", err)
			}

			createHandler := createdom.NewHandler(app)
			if _, err := createHandler.CreateCommand("holder", []string{"g2"}, "echo two"); err != nil {
				t.Fatalf("create holder: %v", err)
			}

			if _, err := createHandler.CreateCommand("other", nil, "echo other"); err != nil {
				t.Fatalf("create other: %v", err)
			}

			handler := updatedom.NewHandler(app)

			_, _, err = handler.UpdateCommand(tt.key, tt.newName, tt.newAlias, "echo other")
			if err == nil {
				t.Fatal("update accepted a shadowing name or alias")
			}

			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("error = %v, want %q", err, tt.want)
			}

			if byName := app.AllCommandsByName.Get("other"); byName == nil {
				t.Error("the rejected update removed the command")
			}
		})
	}
}

// A rename must move the command and drop the aliases it no longer has: the
// old name and the old alias have to disappear from their maps.
func TestUpdateRenamesAndReplacesAliases(t *testing.T) {
	app, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	if _, err := createdom.NewHandler(app).CreateCommand("holder", []string{"g2"}, "echo two"); err != nil {
		t.Fatalf("create holder: %v", err)
	}

	handler := updatedom.NewHandler(app)

	if _, _, err := handler.UpdateCommand("holder", "renamed", []string{"g3"}, "echo new"); err != nil {
		t.Fatalf("rename: %v", err)
	}

	if app.AllCommandsByName.Get("holder") != nil {
		t.Error("the old name is still registered")
	}

	if app.AllCommandsByAlias.Get("g2") != nil {
		t.Error("the replaced alias is still registered")
	}

	if byName := app.AllCommandsByName.Get("renamed"); byName == nil || byName.Exec != "echo new" {
		t.Errorf("renamed command = %+v, want exec echo new", byName)
	}

	if byAlias := app.AllCommandsByAlias.Get("g3"); byAlias == nil || byAlias.Name != "renamed" {
		t.Errorf("new alias = %+v, want renamed", byAlias)
	}
}

// The maps are keyed by name and alias, so a command meets itself. Re-sending
// its own name must not be reported as a conflict, and an omitted alias list
// keeps the aliases it already had.
func TestUpdateAllowsTheCommandsOwnNameAndAliases(t *testing.T) {
	app, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	if _, err := createdom.NewHandler(app).CreateCommand("holder", []string{"g2"}, "echo two"); err != nil {
		t.Fatalf("create holder: %v", err)
	}

	handler := updatedom.NewHandler(app)

	if _, _, err := handler.UpdateCommand("holder", "holder", nil, "echo changed"); err != nil {
		t.Fatalf("re-sending the command's own name failed: %v", err)
	}

	if byName := app.AllCommandsByName.Get("holder"); byName == nil || byName.Exec != "echo changed" {
		t.Errorf("command after update = %+v, want the new exec", byName)
	}

	if byAlias := app.AllCommandsByAlias.Get("g2"); byAlias == nil || byAlias.Name != "holder" {
		t.Errorf("alias after update = %+v, want holder", byAlias)
	}
}

// A config written before the shadow checks (or edited by hand) may already
// hold a collision. Changing an unrelated field must still work: only the
// values the update changes are validated.
func TestUpdateWorksOnAPreExistingCollision(t *testing.T) {
	t.Run("a name that is also an alias", func(t *testing.T) {
		app, _, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		other := command.New("other", []string{"shadow"}, "echo other")
		app.AllCommandsByName.Add(other.Name, other)
		app.AllCommandsByAlias.Add("shadow", other)

		shadow := command.New("shadow", nil, "echo shadow")
		app.AllCommandsByName.Add(shadow.Name, shadow)

		handler := updatedom.NewHandler(app)

		if _, _, err := handler.UpdateCommand("shadow", "", nil, "echo changed"); err != nil {
			t.Fatalf("updating an unrelated field on a collision failed: %v", err)
		}

		if byName := app.AllCommandsByName.Get("shadow"); byName == nil || byName.Exec != "echo changed" {
			t.Errorf("command after update = %+v, want the new exec", byName)
		}
	})

	t.Run("an alias that is also a name", func(t *testing.T) {
		app, _, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		holder := command.New("holder", []string{"g2"}, "echo two")
		app.AllCommandsByName.Add(holder.Name, holder)
		app.AllCommandsByAlias.Add("g2", holder)

		shadow := command.New("g2", nil, "echo shadow")
		app.AllCommandsByName.Add(shadow.Name, shadow)

		handler := updatedom.NewHandler(app)

		if _, _, err := handler.UpdateCommand("holder", "", nil, "echo changed"); err != nil {
			t.Fatalf("updating an unrelated field on a collision failed: %v", err)
		}

		if byName := app.AllCommandsByName.Get("holder"); byName == nil || byName.Exec != "echo changed" {
			t.Errorf("command after update = %+v, want the new exec", byName)
		}
	})
}

func assertUpdateTest(want bool, got bool, key string, newCommand command.Data, err error, t *testing.T) {
	if want != got {
		t.Errorf("error %v, key %v, command %v, got %v, want %v", err, key, newCommand, got, want)
	}
}

// A supplied alias list drops a claim another command owns, and the warning
// says where the alias stayed. The record stops claiming it.
func TestUpdateLeavesAnAliasAnotherCommandOwns(t *testing.T) {
	cfg, stderr := newConfigWithStderr(t)

	a := command.New("a", []string{"x"}, "echo a")
	cfg.AllCommandsByName.Add("a", a)
	cfg.AllCommandsByAlias.Add("x", a)

	b := command.New("b", []string{"x"}, "echo b")
	cfg.AllCommandsByName.Add("b", b)

	handler := updatedom.NewHandler(cfg)

	if _, _, err := handler.UpdateCommand("b", "", []string{"x"}, "echo changed"); err != nil {
		t.Fatalf("update: %v", err)
	}

	if owner := cfg.AllCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Errorf("alias x = %+v, want it to stay on a", owner)
	}

	if got := cfg.AllCommandsByName.Get("b"); got == nil || slices.Contains(got.AllAlias, "x") {
		t.Errorf("b = %+v, want it to stop claiming x", got)
	}

	if !strings.Contains(stderr.String(), `alias "x" belongs to "a"`) {
		t.Errorf("the dropped alias was not reported: %q", stderr.String())
	}
}

// An update of another field must not delete the dormant claim: the record is
// only normalised when the user replaces the alias list.
func TestUpdateKeepsADormantAliasOnAnUnrelatedField(t *testing.T) {
	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	a := command.New("a", []string{"x"}, "echo a")
	cfg.AllCommandsByName.Add("a", a)
	cfg.AllCommandsByAlias.Add("x", a)

	b := command.New("b", []string{"x"}, "echo b")
	cfg.AllCommandsByName.Add("b", b)

	handler := updatedom.NewHandler(cfg)

	if _, _, err := handler.UpdateCommand("b", "", nil, "echo changed"); err != nil {
		t.Fatalf("update: %v", err)
	}

	if owner := cfg.AllCommandsByAlias.Get("x"); owner == nil || owner.Name != "a" {
		t.Errorf("alias x = %+v, want it to stay on a", owner)
	}

	if got := cfg.AllCommandsByName.Get("b"); got == nil || !slices.Contains(got.AllAlias, "x") {
		t.Errorf("b = %+v, want it to keep the dormant claim", got)
	}
}

// The same for an alias that is another command's name: name-first lookup
// makes it unreachable, so a supplied list drops the claim.
func TestUpdateLeavesAnAliasThatIsAnotherName(t *testing.T) {
	cfg, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	b := command.New("b", nil, "echo b")
	cfg.AllCommandsByName.Add("b", b)

	a := command.New("a", []string{"b"}, "echo a")
	cfg.AllCommandsByName.Add("a", a)
	cfg.AllCommandsByAlias.Add("b", a)

	handler := updatedom.NewHandler(cfg)

	if _, _, err := handler.UpdateCommand("a", "", []string{"b"}, "echo changed"); err != nil {
		t.Fatalf("update: %v", err)
	}

	if got := cfg.AllCommandsByName.Get("a"); got == nil || slices.Contains(got.AllAlias, "b") {
		t.Errorf("a = %+v, want it to stop claiming b as an alias", got)
	}

	if cfg.AllCommandsByAlias.Get("b") != nil {
		t.Errorf("alias b = %+v, want it removed", cfg.AllCommandsByAlias.Get("b"))
	}
}

func newConfigWithStderr(t *testing.T) (*app.Config, *bytes.Buffer) {
	t.Helper()

	var stderr bytes.Buffer

	cfgPath := filepath.Join(t.TempDir(), "crower.yaml")
	cfg, err := app.InitApp(cfgPath, app.Options{
		Streams: app.Streams{In: strings.NewReader(""), Out: io.Discard, Err: &stderr},
	})
	if err != nil {
		t.Fatalf("InitApp: %v", err)
	}

	return cfg, &stderr
}

// A Steps command runs its steps instead of its exec. update has no flag for
// them, so an update that rebuilt the command without them replaced a
// multi-line workflow with a single line and reported success (#217).
func TestUpdateKeepsSteps(t *testing.T) {
	tests := []struct {
		name     string
		newName  string
		newAlias []string
		newExec  string
	}{
		{name: "an exec change", newExec: "echo new"},
		{name: "a rename", newName: "renamed"},
		{name: "an alias change", newAlias: []string{"m2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, _, err := testutil.New(t)
			if err != nil {
				t.Fatalf("error before test: %v", err)
			}

			steps := []string{"echo one", "echo two"}
			stored := &command.Data{
				Name:     "multi",
				AllAlias: []string{"m"},
				Exec:     "echo last",
				Steps:    steps,
			}

			cfg.AllCommandsByName.Add("multi", stored)
			cfg.AllCommandsByAlias.Add("m", stored)

			handler := updatedom.NewHandler(cfg)

			if _, _, err := handler.UpdateCommand("multi", tt.newName, tt.newAlias, tt.newExec); err != nil {
				t.Fatalf("UpdateCommand: %v", err)
			}

			key := "multi"
			if tt.newName != "" {
				key = tt.newName
			}

			updated := cfg.AllCommandsByName.Get(key)
			if updated == nil {
				t.Fatalf("command %q is missing after the update", key)
			}

			if !slices.Equal(updated.Steps, steps) {
				t.Errorf("Steps = %v, want %v", updated.Steps, steps)
			}
		})
	}
}
