package utils

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// A command's nested table (Env) must not be read as a second command: the
// order reader takes the first key of every table, so [deploy.Env] used to add
// "deploy" again and shift the command on every save.
func TestReadKeysTomlInOrderIgnoresNestedTables(t *testing.T) {
	path := filepath.Join(t.TempDir(), "crower.yaml")

	body := `[zebra]
  Name = "zebra"

[alpha]
  Name = "alpha"

  [alpha.Env]
    STAGE = "prod"

[beta]
  Name = "beta"
`

	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatalf("write: %v", err)
	}

	keys, err := ReadKeysTomlInOrder(path)
	if err != nil {
		t.Fatalf("ReadKeysTomlInOrder: %v", err)
	}

	if want := []string{"zebra", "alpha", "beta"}; !slices.Equal(keys, want) {
		t.Errorf("keys = %v, want %v", keys, want)
	}
}
