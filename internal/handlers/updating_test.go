package handlers_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func TestUpdate(t *testing.T) {
	t.Run("Update command name", func(t *testing.T) {
		app := data.NewApp("", data.NewCommandsMap(), data.NewCommandsMap())

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
			handlers.AddCommand(
				data.NewCommand(command.name, []string{command.alias}, ""),
				app)
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
			newCommand := data.NewCommand(test.newName, []string{}, "")
			key := test.oldName
			got := handlers.UpdateCommand(key, newCommand, app)
			assertUpdatingTest(test.want, got, key, *newCommand, t)

		}

	})
}

func assertUpdatingTest(want bool, got bool, key string, newCommand data.Command, t *testing.T) {
	if want != got {
		t.Errorf("command %v,  key %v,  got %v, want %v", newCommand, key, got, want)
	}
}
