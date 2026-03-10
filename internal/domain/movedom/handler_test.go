package movedom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/movedom"
	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestMoveFile(t *testing.T) {

	t.Run("Move File single name", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		moveHandler := movedom.NewHandler(app)

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = moveHandler.MoveFile(newFilePath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move File name with space", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		moveHandler := movedom.NewHandler(app)

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = moveHandler.MoveFile(newFilePath, testPaths[0])
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

		createHandler := createdom.NewHandler(app)
		moveHandler := movedom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = moveHandler.MoveFolder(newFolderPath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move Folder using name with space", func(t *testing.T) {
		app, testPaths, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		moveHandler := movedom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = moveHandler.MoveFolder(newFolderPath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}
