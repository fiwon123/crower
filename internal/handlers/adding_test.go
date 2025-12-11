package handlers_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data"
	"github.com/fiwon123/crower/internal/handlers"
)

func TestAdd(t *testing.T) {
	t.Run("Add a single command using only name", func(t *testing.T) {

		command := data.Command{
			Name: "c1",
		}

		want := true
		app := core.InitApp("")

		assertAddingTest(command, want, handlers.AddCommand(&command, app), t)
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
			command := data.Command{
				Name: test.name,
			}

			assertAddingTest(command, test.want, handlers.AddCommand(&command, app), t)
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
			command := data.Command{
				AllAlias: []string{test.alias},
				Name:     test.name,
			}

			assertAddingTest(command, test.want, handlers.AddCommand(&command, app), t)
		}

	})
}

func assertAddingTest(command data.Command, want bool, got bool, t *testing.T) {
	if got != want {
		t.Errorf("command %+v got %v, want %v", command, got, want)
	}
}
