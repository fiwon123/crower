package copyhandlers_test

import (
	"path/filepath"
	"testing"

	copyhandlers "github.com/fiwon123/crower/internal/handlers/copy"
	createhandlers "github.com/fiwon123/crower/internal/handlers/create"
	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestCopyFile(t *testing.T) {

	t.Run("Copy File single name", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = createhandlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = copyhandlers.CopyFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Copy File name with space", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = createhandlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = copyhandlers.CopyFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

}

func TestCopyFolder(t *testing.T) {

	t.Run("Copy Folder single name", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createhandlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = copyhandlers.CopyFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Copy Folder using name with space", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createhandlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = copyhandlers.CopyFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}
