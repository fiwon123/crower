package openhandlers_test

import (
	"path/filepath"
	"testing"

	createhandlers "github.com/fiwon123/crower/internal/handlers/create"
	openhandlers "github.com/fiwon123/crower/internal/handlers/open"
	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestOpenFile(t *testing.T) {
	t.Run("Open File single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = createhandlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = openhandlers.Open([]string{newFilePath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Open File name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = createhandlers.CreateFile(newFilePath, app)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = openhandlers.Open([]string{newFilePath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func TestOpenFolder(t *testing.T) {
	t.Run("Open c single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createhandlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = openhandlers.Open([]string{newFolderPath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Open Folder name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createhandlers.CreateFolder(newFolderPath, app)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = openhandlers.Open([]string{newFolderPath}, app)
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})
}

func TestOpenSystem(t *testing.T) {
	app, _, err := testshelper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	err = openhandlers.OpenSystem(app)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

}
