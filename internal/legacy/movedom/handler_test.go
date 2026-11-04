package movedom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/system"
	"github.com/fiwon123/crower/internal/legacy/movedom"
	"github.com/fiwon123/crower/internal/testutil"
)

func TestMoveFile(t *testing.T) {

	t.Run("Move File single name", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		creator := system.NewCreator()
		moveHandler := movedom.NewHandler(app)

		newFilePath := filepath.Join(app.ConfigFolder(), "new.txt")

		err = creator.File(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = moveHandler.MoveFile(newFilePath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move File name with space", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		creator := system.NewCreator()
		moveHandler := movedom.NewHandler(app)

		newFilePath := filepath.Join(app.ConfigFolder(), "new file.txt")

		err = creator.File(newFilePath)
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
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		creator := system.NewCreator()
		moveHandler := movedom.NewHandler(app)

		newFolderPath := filepath.Join(app.ConfigFolder(), "new")

		err = creator.Folder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = moveHandler.MoveFolder(newFolderPath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Move Folder using name with space", func(t *testing.T) {
		app, testPaths, err := testutil.New(t)
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		creator := system.NewCreator()
		moveHandler := movedom.NewHandler(app)

		newFolderPath := filepath.Join(app.ConfigFolder(), "new folder")

		err = creator.Folder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = moveHandler.MoveFolder(newFolderPath, testPaths[0])
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}
