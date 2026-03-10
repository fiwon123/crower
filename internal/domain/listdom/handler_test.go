package listdom_test

import (
	"path/filepath"
	"testing"

	"github.com/fiwon123/crower/internal/domain/listdom"
	"github.com/fiwon123/crower/internal/helper"
)

func TestList(t *testing.T) {
	app, _, err := helper.InitCrowerTests()
	if err != nil {
		t.Fatalf("error before test: %v", err)
	}

	handler := listdom.NewHandler(app)

	t.Run("List Folder", func(t *testing.T) {
		output, err := handler.ListFolder(filepath.Dir(app.CfgFilePath))
		if err != nil {
			t.Errorf("list Folder error: %v , out: %s", err, string(output))
		}
	})

	t.Run("List History", func(t *testing.T) {
		err := handler.ListHistory()
		if err != nil {
			t.Errorf("list History error: %v", err)
		}
	})

	t.Run("List Process", func(t *testing.T) {
		err := handler.ListProcess([]string{})
		if err != nil {
			t.Errorf("list Process error: %v", err)
		}
	})

	t.Run("List System Path", func(t *testing.T) {
		output, err := handler.ListSysPath()
		if err != nil {
			t.Errorf("list System Path error: %v , out: %s", err, string(output))
		}
	})

	t.Run("List System Variables", func(t *testing.T) {
		output, err := handler.ListSystem()
		if err != nil {
			t.Errorf("list System Variables error: %v , out: %s", err, string(output))
		}
	})
}
