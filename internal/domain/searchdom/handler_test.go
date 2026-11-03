package searchdom_test

import (
	"testing"

	"github.com/fiwon123/crower/internal/domain/searchdom"
	"github.com/fiwon123/crower/internal/helper"
)

func TestSearchBrowser(t *testing.T) {
	t.Run("Search Browser single name", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
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
		app, _, err := helper.InitCrowerTests()
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
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFile(app.ConfigFolder(), "test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search File name with space", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFile(app.ConfigFolder(), "test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFolder(t *testing.T) {
	t.Run("Search Folder single name", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFolder(app.ConfigFolder(), "test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder name with space", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFolder(app.ConfigFolder(), "test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}

func TestSearchFileAndFolder(t *testing.T) {
	t.Run("Search Folder and File single name", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFileAndFolder(app.ConfigFolder(), "test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})

	t.Run("Search Folder anf File name with space", func(t *testing.T) {
		app, _, err := helper.InitCrowerTests()
		if err != nil {
			t.Fatalf("error before test: %v", err)
		}

		handler := searchdom.NewHandler(app)

		output, err := handler.SearchFileAndFolder(app.ConfigFolder(), "test test")
		if err != nil {
			t.Fatalf("error: %v output: %s", err, string(output))
		}
	})
}
