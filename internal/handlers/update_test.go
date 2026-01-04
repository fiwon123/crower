package handlers_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestUpdate(t *testing.T) {

	t.Run("Update command name", func(t *testing.T) {

		app, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

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
			handlers.CreateCommand(
				command.name, []string{command.alias}, "exec", nil,
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
			newCommand := command.New(test.newName, []string{}, "")
			key := test.oldName
			_, _, err := handlers.UpdateCommand(key, newCommand.Name, newCommand.AllAlias, newCommand.Exec, app)
			got := err == nil
			assertUpdateTest(test.want, got, key, *newCommand, err, t)

		}

	})
}

func assertUpdateTest(want bool, got bool, key string, newCommand command.Data, err error, t *testing.T) {
	if want != got {
		t.Errorf("error %v, key %v, command %v, got %v, want %v", err, key, newCommand, got, want)
	}
}
