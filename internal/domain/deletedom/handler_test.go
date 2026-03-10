package deletedom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/deletedom"
	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestDelete(t *testing.T) {

	t.Run("Delete command using name", func(t *testing.T) {

		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		deleteHandler := deletedom.NewHandler(app)

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
			createHandler.CreateCommand(command.name, nil, "exec")
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
			_, got := deleteHandler.DeleteCommand(test.name)

			assertDeleteTest(test.want, got, t)
		}
	})

	t.Run("Delete command using alias", func(t *testing.T) {

		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		deleteHandler := deletedom.NewHandler(app)

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
			createHandler.CreateCommand(
				command.name, []string{command.alias}, "exec")
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
			_, got := deleteHandler.DeleteCommand(test.alias)

			assertDeleteTest(test.want, got, t)
		}
	})

}

func TestDeleteFile(t *testing.T) {
	t.Run("Delete File single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		deleteHandler := deletedom.NewHandler(app)

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = deleteHandler.DeleteFile(newFilePath)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Delete File using name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		deleteHandler := deletedom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file")

		err = createHandler.CreateFile(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = deleteHandler.DeleteFile(newFolderPath)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func TestDeleteFolder(t *testing.T) {

	t.Run("Delete Folder single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}
		createHandler := createdom.NewHandler(app)
		deleteHandler := deletedom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = deleteHandler.DeleteFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Delete Folder using name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		deleteHandler := deletedom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = deleteHandler.DeleteFolder(newFolderPath)
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
