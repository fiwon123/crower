package searchhandlers_test

import (
	"path/filepath"
	"testing"

	searchhandlers "github.com/fiwon123/crower/internal/handlers/search"
	crowertests "github.com/fiwon123/crower/internal/helper/tests"
)

func TestSearchBrowser(t *testing.T) {
	t.Run("Search Browser single name", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchBrowser("test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Browser name with space", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchBrowser("test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFile(t *testing.T) {
	t.Run("Search File single name", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchFile(filepath.Dir(app.CfgFilePath), "test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search File name with space", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchFile(filepath.Dir(app.CfgFilePath), "test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFolder(t *testing.T) {
	t.Run("Search Folder single name", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchFolder(filepath.Dir(app.CfgFilePath), "test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder name with space", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchFolder(filepath.Dir(app.CfgFilePath), "test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFileAndFolder(t *testing.T) {
	t.Run("Search Folder and File single name", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchFileAndFolder(filepath.Dir(app.CfgFilePath), "test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder anf File name with space", func(t *testing.T) {
		app, _, err := crowertests.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		output, err := searchhandlers.SearchFileAndFolder(filepath.Dir(app.CfgFilePath), "test test", app)
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}
