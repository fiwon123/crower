package createdom_test

import (
	"path/filepath"
	"strings"
	"testing"

	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/helper"
)

func TestCreate(t *testing.T) {

	t.Run("Create a single command using only name", func(t *testing.T) {

		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := createdom.NewHandler(app)

		command := command.New("c1", nil, "")

		want := true

		_, error := handler.CreateCommand(command.Name, nil, "exec")
		got := error == nil
		assertCreateTest(command, want, got, error, t)
	})

	t.Run("Create multiple commands using only name", func(t *testing.T) {

		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := createdom.NewHandler(app)

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

			_, err := handler.CreateCommand(command.Name, nil, "exec")
			got := err == nil
			assertCreateTest(command, test.want, got, err, t)
		}

	})

	t.Run("Create multiple commands using only name, alias", func(t *testing.T) {

		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := createdom.NewHandler(app)

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
			_, err := handler.CreateCommand(command.Name, command.AllAlias, "exec")
			got := err == nil
			assertCreateTest(command, test.want, got, err, t)
		}

	})
}

// A process create without a process is a usage error, not a panic: both the
// missing-argument form and the form where the only word was taken as the name
// used to index an empty slice.
func TestCreateProcessRequiresAProcess(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	handler := createdom.NewHandler(app)

	tests := []struct {
		name        string
		commandName string
		args        []string
	}{
		{"no arguments", "name", nil},
		{"the only word is taken as the name", "", []string{"only-name"}},
		{"an empty process argument", "name", []string{""}},
		{"a whitespace process argument", "name", []string{"   "}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.CreateProcess(tt.commandName, tt.args)
			if err == nil {
				t.Fatal("CreateProcess with no process returned no error")
			}

			if !strings.Contains(err.Error(), "process name or pid") {
				t.Errorf("error = %v, want it to name the missing process", err)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {

	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}
	handler := createdom.NewHandler(app)

	t.Run("Create file using single name", func(t *testing.T) {
		newFilePath := filepath.Join(app.ConfigFolder(), "new.txt")
		err := handler.CreateFile(newFilePath)

		if err != nil {
			t.Errorf("error create file on %s: %v", newFilePath, err)
		}
	})

	t.Run("Create file using name with space", func(t *testing.T) {
		newFilePath := filepath.Join(app.ConfigFolder(), "new file.txt")
		err := handler.CreateFile(newFilePath)

		if err != nil {
			t.Errorf("error create file on %s: %v", newFilePath, err)
		}
	})
}

func TestCreateFolder(t *testing.T) {

	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}
	handler := createdom.NewHandler(app)

	t.Run("Create folder using single name", func(t *testing.T) {
		newFolderPath := filepath.Join(app.ConfigFolder(), "new")
		err := handler.CreateFile(newFolderPath)

		if err != nil {
			t.Errorf("error create folder on %s: %v", newFolderPath, err)
		}
	})

	t.Run("Create folder using name with space", func(t *testing.T) {
		newFolderPath := filepath.Join(app.ConfigFolder(), "new folder")
		err := handler.CreateFile(newFolderPath)

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
