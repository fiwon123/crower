package updatedom

import (
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/testutil"
)

// The prompt and the handler share findCommand, so the command shown for
// confirmation is the command that gets updated. A colliding key resolves name
// first, as every lookup in the project does.
func TestFindCommand(t *testing.T) {
	app, _, err := testutil.New(t)
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	byName := command.New("key", nil, "echo name")
	app.AllCommandsByName.Add(byName.Name, byName)

	byAlias := command.New("owner", []string{"key", "alias-only"}, "echo alias")
	app.AllCommandsByName.Add(byAlias.Name, byAlias)
	app.AllCommandsByAlias.Add("key", byAlias)
	app.AllCommandsByAlias.Add("alias-only", byAlias)

	tests := []struct {
		name string
		key  string
		want string
	}{
		{"a name", "owner", "owner"},
		{"an alias", "alias-only", "owner"},
		{"a collision resolves the name first", "key", "key"},
		{"missing", "nope", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findCommand(app, tt.key)
			if tt.want == "" {
				if got != nil {
					t.Errorf("findCommand(%q) = %+v, want nil", tt.key, got)
				}
				return
			}

			if got == nil || got.Name != tt.want {
				t.Errorf("findCommand(%q) = %+v, want %q", tt.key, got, tt.want)
			}
		})
	}
}
