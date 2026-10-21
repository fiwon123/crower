package helper

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fiwon123/crower/internal/app"
	"github.com/fiwon123/crower/pkg/utils"
)

// InitCrowerTests prepares an isolated environment for a test.
//
// Every call gets its own directory. The previous implementation built one
// fixed path, ~/crower/test, removed it with os.RemoveAll and recreated it.
// Every test package calls this helper, and `go test ./...` runs the packages
// in parallel, so one package deleted the fixtures another package was using
// -- the suite failed in opendom, deletedom and movedom, and passed when a
// failing package was run alone. The fixed path also put the fixtures in the
// user's home directory, one directory away from the real crower data, and
// os.RemoveAll deleted whatever was there.
//
// The directory now comes from os.MkdirTemp: it is unique per call, and it
// lives in the system temporary directory instead of the user's home.
// Everything the tests touch -- the config file, the history file, the history
// folder and the log -- is derived from that directory, so one test's fixtures
// are invisible to another.
//
// The temporary directory is not removed automatically. Removing it needs the
// *testing.T that owns the test, which this helper does not receive; passing
// one would pull the testing package into a package the binary links.
// Follow-up: move this helper into a test-only package and use t.TempDir().
func InitCrowerTests() (*app.Config, []string, error) {
	testPath, err := os.MkdirTemp("", "crower-test-")
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

	return app.InitApp(filepath.Join(testPath, "crower.yaml")), testPaths, nil

}
