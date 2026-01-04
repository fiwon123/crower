package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestSearchBrowser(t *testing.T) {
	t.Run("Search Browser single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchBrowser("test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Browser name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchBrowser("test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFile(t *testing.T) {
	t.Run("Search File single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchFile(filepath.Dir(app.CfgFilePath), "test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search File name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchFile(filepath.Dir(app.CfgFilePath), "test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFolder(t *testing.T) {
	t.Run("Search Folder single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchFolder(filepath.Dir(app.CfgFilePath), "test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchFolder(filepath.Dir(app.CfgFilePath), "test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFileAndFolder(t *testing.T) {
	t.Run("Search Folder and File single name", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchFileAndFolder(filepath.Dir(app.CfgFilePath), "test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder anf File name with space", func(t *testing.T) {
		app, _, err := crtests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := handlers.SearchFileAndFolder(filepath.Dir(app.CfgFilePath), "test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}
