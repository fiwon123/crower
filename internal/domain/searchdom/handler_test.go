package searchdom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/searchdom"
	testshelper "github.com/fiwon123/crower/internal/helper/tests"
)

func TestSearchBrowser(t *testing.T) {
	t.Run("Search Browser single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchBrowser("test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Browser name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchBrowser("test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFile(t *testing.T) {
	t.Run("Search File single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFile(filepath.Dir(app.CfgFilePath), "test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search File name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFile(filepath.Dir(app.CfgFilePath), "test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFolder(t *testing.T) {
	t.Run("Search Folder single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFolder(filepath.Dir(app.CfgFilePath), "test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFolder(filepath.Dir(app.CfgFilePath), "test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFileAndFolder(t *testing.T) {
	t.Run("Search Folder and File single name", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFileAndFolder(filepath.Dir(app.CfgFilePath), "test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder anf File name with space", func(t *testing.T) {
		app, _, err := testshelper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFileAndFolder(filepath.Dir(app.CfgFilePath), "test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}
