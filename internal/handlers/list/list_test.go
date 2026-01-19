package listhandlers_test

import (
	"path/filepath"
	"testing"

	listhandlers "github.com/fiwon123/crower/internal/handlers/list"
	crowertests "github.com/fiwon123/crower/internal/helper/tests"
)

func TestList(t *testing.T) {
	app, _, err := crowertests.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	t.Run("List Folder", func(t *testing.T) {
		output, err := listhandlers.ListFolder(filepath.Dir(app.CfgFilePath), app)
		if err != nil {
			t.Errorf("list Folder error: %v , out: %s", err, string(output))
		}
	})

	t.Run("List History", func(t *testing.T) {
		err := listhandlers.ListHistory(app)
		if err != nil {
			t.Errorf("list History error: %v", err)
		}
	})

	t.Run("List Process", func(t *testing.T) {
		err := listhandlers.ListProcess([]string{}, app)
		if err != nil {
			t.Errorf("list Process error: %v", err)
		}
	})

	t.Run("List System Path", func(t *testing.T) {
		output, err := listhandlers.ListSysPath(app)
		if err != nil {
			t.Errorf("list System Path error: %v , out: %s", err, string(output))
		}
	})

	t.Run("List System Variables", func(t *testing.T) {
		output, err := listhandlers.ListSystem(app)
		if err != nil {
			t.Errorf("list System Variables error: %v , out: %s", err, string(output))
		}
	})
}
