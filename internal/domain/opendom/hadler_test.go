package opendom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/createdom"
	"github.com/fiwon123/crower/internal/domain/opendom"

	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestOpenFile(t *testing.T) {
	t.Run("Open File single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		openHandler := opendom.NewHandler(app)

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = openHandler.Open([]string{newFilePath})
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Open File name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		openHandler := opendom.NewHandler(app)

		newFilePath := filepath.Join(filepath.Dir(app.CfgFilePath), "new file.txt")

		err = createHandler.CreateFile(newFilePath)
		if err != nil {
			t.Fatalf("error before test create file: %v", err)
		}

		err = openHandler.Open([]string{newFilePath})
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

		createHandler := createdom.NewHandler(app)
		openHandler := opendom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = openHandler.Open([]string{newFolderPath})
		if err != nil {
			t.Fatalf("error: %v", err)
		}
	})

	t.Run("Open Folder name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		createHandler := createdom.NewHandler(app)
		openHandler := opendom.NewHandler(app)

		newFolderPath := filepath.Join(filepath.Dir(app.CfgFilePath), "new folder")

		err = createHandler.CreateFolder(newFolderPath)
		if err != nil {
			t.Fatalf("error before test create folder: %v", err)
		}

		err = openHandler.Open([]string{newFolderPath})
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

	openHandler := opendom.NewHandler(app)

	err = openHandler.OpenSystem()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

}
