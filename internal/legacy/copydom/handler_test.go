package copydom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/legacy/copydom"
	"github.com/fiwon123/crower/internal/legacy/createdom"
	"github.com/fiwon123/crower/internal/testutil"
)

func TestCopyFile(t *testing.T) {

	t.Run("Copy File single name", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		copyHandler := copydom.NewHandler(app)

		newFilePath := filepath.Join(app.ConfigFolder(), "new.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = copyHandler.CopyFile(newFilePath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Copy File name with space", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		copyHandler := copydom.NewHandler(app)

		newFilePath := filepath.Join(app.ConfigFolder(), "new file.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = copyHandler.CopyFile(newFilePath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

}

func TestCopyFolder(t *testing.T) {

	t.Run("Copy Folder single name", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		copyHandler := copydom.NewHandler(app)

		newFolderPath := filepath.Join(app.ConfigFolder(), "new")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = copyHandler.CopyFolder(newFolderPath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Copy Folder using name with space", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}
		createHandler := createdom.NewHandler(app)
		copyHandler := copydom.NewHandler(app)
		newFolderPath := filepath.Join(app.ConfigFolder(), "new folder")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = copyHandler.CopyFolder(newFolderPath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}
