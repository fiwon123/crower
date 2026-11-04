// Package testutil builds an isolated crower installation for a test.
package testutil

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/fiwon123/crower/internal/session"

	"github.com/fiwon123/crower/pkg/utils"
)

// New returns a configuration whose state lives in a directory the test owns,
// and three empty folders inside it.
//
// The directory comes from t.TempDir: it is unique per test, and the testing
// package removes it when the test ends. The helper this replaces used
// os.MkdirTemp and never cleaned up, because it had no *testing.T to hand it
// to -- every run of the suite left its fixtures behind in the system temp
// directory.
//
// Passing t also means a failure here is reported against the test that owns
// the fixture, not against a line inside the helper.
//
// It returns an error rather than calling t.Fatal so the signature stays the
// one the call sites already use; the value is always nil in practice.
func New(t *testing.T) (*session.Session, []string, error) {
	t.Helper()

	root := t.TempDir()

	paths := make([]string, 0, 3)

	for i := range 3 {
		folder := filepath.Join(root, "test"+strconv.Itoa(i))
		if err := utils.CreateFolderIfNotExists(folder); err != nil {
			return nil, nil, err
		}

		paths = append(paths, folder)
	}

	cfg, err := session.Load(filepath.Join(root, "crower.yaml"), session.Options{})
	if err != nil {
		return nil, nil, err
	}

	return cfg, paths, nil
}
