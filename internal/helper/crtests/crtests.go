package crtests

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/core"
	"github.com/fiwon123/crower/internal/data/app"
	"github.com/fiwon123/crower/pkg/utils"
)

func InitCrowerTests() (*app.Data, []string, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, nil, err
	}

	testPath := filepath.Join(homePath, "test", "crower")
	os.RemoveAll(testPath)
	err = utils.CreateFolderIfNotExists(testPath)
	if err != nil {
		return nil, nil, err
	}

	testPaths := []string{}
	//test folders
	for i := range 3 {
		newTestFolderBuilder := strings.Builder{}
		newTestFolderBuilder.WriteString("test")
		newTestFolderBuilder.WriteString(strconv.Itoa(i))

		folderPath := filepath.Join(testPath, newTestFolderBuilder.String())
		testPaths = append(testPaths, folderPath)
		err = utils.CreateFolderIfNotExists(folderPath)
		if err != nil {
			return nil, nil, err
		}
	}

	return core.InitApp(filepath.Join(testPath, "crower.yaml")), testPaths, nil

}
