package handlers_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/command"

	"github.com/fiwon123/crower/internal/handlers"
)

func TestAdd(t *testing.T) {
	t.Run("Add a single command using only name", func(t *testing.T) {

		command := command.New("c1", nil, "")

		want := true
		app := core.InitApp("")

		_, error := handlers.AddCommand(command.Name, nil, "exec", nil, app)
		got := error == nil
		assertAddTest(command, want, got, error, t)
	})

	t.Run("Add multiple commands using only name", func(t *testing.T) {

		var tests = []struct {
			name string
			want bool
		}{
			{"1", true},
			{"1", false},
			{"2", true},
			{"3", true},
			{"4", true},
			{"5", true},
			{"4", false},
			{"2", false},
		}

		app := core.InitApp("")

		for _, test := range tests {
			command := command.New(test.name, nil, "")

			_, err := handlers.AddCommand(command.Name, nil, "exec", nil, app)
			got := err == nil
			assertAddTest(command, test.want, got, err, t)
		}

	})

	t.Run("Add multiple commands using only name, alias", func(t *testing.T) {

		var tests = []struct {
			name  string
			alias string
			want  bool
		}{
			{"1", "2", true},
			{"1", "3", false},
			{"2", "1", false},
			{"3", "5", true},
			{"4", "6", true},
			{"5", "2", false},
			{"6", "2", false},
			{"7", "1", false},
		}

		app := core.InitApp("")

		for _, test := range tests {
			command := command.New(test.name, []string{test.alias}, "")
			_, err := handlers.AddCommand(command.Name, command.AllAlias, "exec", nil, app)
			got := err == nil
			assertAddTest(command, test.want, got, err, t)
		}

	})
}

func assertAddTest(command *command.Data, want bool, got bool, err error, t *testing.T) {
	if got != want {
		t.Errorf("error %v, command %+v got %v, want %v", err, *command, got, want)
	}
}
