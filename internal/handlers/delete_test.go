package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestDelete(t *testing.T) {

	t.Run("Delete command using name", func(t *testing.T) {

		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		var mock = []struct {
			name string
		}{
			{"c1"},
			{"c2"},
			{"c3"},
			{"c4"},
			{"c5"},
		}

		for _, command := range mock {
			handlers.CreateCommand(command.name, nil, "exec", nil,
				app)
		}

		var tests = []struct {
			name string
			want bool
		}{
			{"c1", true},
			{"c1", false},
			{"c3", true},
			{"c4", true},
			{"c3", false},
			{"aa", false},
		}

		for _, test := range tests {
			_, got := handlers.DeleteCommand(test.name, app)

			assertDeleteTest(test.want, got, t)
		}
	})

	t.Run("Delete command using alias", func(t *testing.T) {

		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		var mock = []struct {
			name  string
			alias string
		}{
			{"c1", "a1"},
			{"c2", "a2"},
			{"c3", "a3"},
			{"c4", "a4"},
			{"c5", "a5"},
		}

		for _, command := range mock {
			handlers.CreateCommand(
				command.name, []string{command.alias}, "exec", nil,
				app)
		}

		var tests = []struct {
			alias string
			want  bool
		}{
			{"a1", true},
			{"a7", false},
			{"a2", true},
			{"a3", true},
			{"a1", false},
			{"a2", false},
		}

		for _, test := range tests {
			_, got := handlers.DeleteCommand(test.alias, app)

			assertDeleteTest(test.want, got, t)
		}
	})

}

func TestDeleteFile(t *testing.T) {
	t.Run("Delete File single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.DeleteFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Delete File using name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file")

		err = handlers.CreateFile(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.DeleteFile(newFolderPath, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func TestDeleteFolder(t *testing.T) {

	t.Run("Delete Folder single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.DeleteFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Delete Folder using name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.DeleteFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func assertDeleteTest(want bool, got bool, t *testing.T) {
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
