package movehandlers_test

import (
	"path/filepath"
	"testing"

	createhandlers "github.com/fiwon123/crower/internal/handlers/create"
	movehandlers "github.com/fiwon123/crower/internal/handlers/move"
	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestMoveFile(t *testing.T) {

	t.Run("Move File single name", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = createhandlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = movehandlers.MoveFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move File name with space", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = createhandlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = movehandlers.MoveFile(newFilePath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

}

func TestMoveFolder(t *testing.T) {

	t.Run("Move Folder single name", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createhandlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = movehandlers.MoveFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move Folder using name with space", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createhandlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = movehandlers.MoveFolder(newFolderPath, testPaths[0], app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}
