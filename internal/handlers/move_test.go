package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestMoveFile(t *testing.T) {

	t.Run("Move File single name", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.MoveFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move File name with space", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.MoveFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

}

func TestMoveFolder(t *testing.T) {

	t.Run("Move Folder single name", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.MoveFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move Folder using name with space", func(t *testing.T) {
		app, testPaths, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.MoveFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}
