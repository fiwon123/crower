package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/helper/crtests"

	"github.com/fiwon123/crower/internal/handlers"
)

func TestCreate(t *testing.T) {

	t.Run("Create a single command using only name", func(t *testing.T) {

		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		command := command.New("c1", nil, "")

		want := true

		_, error := handlers.CreateCommand(command.Name, nil, "exec", nil, app)
		got := error == nil
		assertCreateTest(command, want, got, error, t)
	})

	t.Run("Create multiple commands using only name", func(t *testing.T) {

		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

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

		for _, test := range tests {
			command := command.New(test.name, nil, "")

			_, err := handlers.CreateCommand(command.Name, nil, "exec", nil, app)
			got := err == nil
			assertCreateTest(command, test.want, got, err, t)
		}

	})

	t.Run("Create multiple commands using only name, alias", func(t *testing.T) {

		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

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

		for _, test := range tests {
			command := command.New(test.name, []string{test.alias}, "")
			_, err := handlers.CreateCommand(command.Name, command.AllAlias, "exec", nil, app)
			got := err == nil
			assertCreateTest(command, test.want, got, err, t)
		}

	})
}

func TestCreateFile(t *testing.T) {

	app, _, err := crtests.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	t.Run("Create file using single name", func(t *testing.T) {
		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")
		err := handlers.CreateFile(newFilePath, app)

		if err != nil {
			t.Errorf("error create file on %s: %v", newFilePath, err)
		}
	})

	t.Run("Create file using name with space", func(t *testing.T) {
		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")
		err := handlers.CreateFile(newFilePath, app)

		if err != nil {
			t.Errorf("error create file on %s: %v", newFilePath, err)
		}
	})
}

func TestCreateFolder(t *testing.T) {

	app, _, err := crtests.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	t.Run("Create folder using single name", func(t *testing.T) {
		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")
		err := handlers.CreateFile(newFolderPath, app)

		if err != nil {
			t.Errorf("error create folder on %s: %v", newFolderPath, err)
		}
	})

	t.Run("Create folder using name with space", func(t *testing.T) {
		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")
		err := handlers.CreateFile(newFolderPath, app)

		if err != nil {
			t.Errorf("error create folder on %s: %v", newFolderPath, err)
		}
	})

}

func assertCreateTest(command *command.Data, want bool, got bool, err error, t *testing.T) {
	if got != want {
		t.Errorf("error %v, command %+v got %v, want %v", err, *command, got, want)
	}
}
