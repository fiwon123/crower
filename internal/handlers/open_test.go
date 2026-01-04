package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestOpenFile(t *testing.T) {
	t.Run("Open File single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.Open([]string{newFilePath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Open File name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = handlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = handlers.Open([]string{newFilePath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func TestOpenFolder(t *testing.T) {
	t.Run("Open c single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.Open([]string{newFolderPath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Open Folder name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = handlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = handlers.Open([]string{newFolderPath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func TestOpenSystem(t *testing.T) {
	app, _, err := crtests.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	err = handlers.OpenSystem(app)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

}
