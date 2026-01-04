package handlers_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/handlers"
	"github.com/fiwon123/crower/internal/helper/crtests"
)

func TestList(t *testing.T) {
	app, _, err := crtests.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	t.Run("List Folder", func(t *testing.T) {
		output, err := handlers.ListFolder(filepath.Dir(app.CfgFilePath), app)
		if err != nil {
			t.Errorf("list Folder error: %v , out: %s", err, string(output))
		}
	})

	t.Run("List History", func(t *testing.T) {
		err := handlers.ListHistory(app)
		if err != nil {
			t.Errorf("list History error: %v", err)
		}
	})

	t.Run("List Process", func(t *testing.T) {
		err := handlers.ListProcess([]string{}, app)
		if err != nil {
			t.Errorf("list Process error: %v", err)
		}
	})

	t.Run("List System Path", func(t *testing.T) {
		output, err := handlers.ListSysPath(app)
		if err != nil {
			t.Errorf("list System Path error: %v , out: %s", err, string(output))
		}
	})

	t.Run("List System Variables", func(t *testing.T) {
		output, err := handlers.ListSystem(app)
		if err != nil {
			t.Errorf("list System Variables error: %v , out: %s", err, string(output))
		}
	})
}
