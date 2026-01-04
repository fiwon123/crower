package crtests

import (
	"os"
	"path/filepath"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/utils"
)

func InitCrowerTests() (*app.Data, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	testPath := filepath.Join(homePath, "test", "crower")
	os.RemoveAll(testPath)
	err = utils.CreateFolderIfNotExists(testPath)
	if err != nil {
		return nil, err
	}

	return core.InitApp(filepath.Join(testPath, "crower.yaml")), nil

}
