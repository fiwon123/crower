package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestCopyFile(t *testing.T) {

	t.Run("Copy File single name", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.CopyFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}
	})

	t.Run("Copy File name with space", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.CopyFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}
	})

}

func TestCopyFolder(t *testing.T) {

	t.Run("Copy Folder single name", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.CopyFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}
	})

	t.Run("Copy Folder using name with space", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.CopyFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}
	})
}
