package shell

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/fiwon123/crower/internal/errors"
)

// ValidateEnv checks the environment entries a command declares.
//
// It is the rule the runner and the write path share: a command that reaches
// the executor has already passed it, and a pack cannot be imported with a
// literal that the executor would refuse. References are checked for shape
// only; resolving ${VAR} is an execution-time question because the value comes
// from the process environment, not from the file.
func ValidateEnv(env map[string]string) error {
	for _, key := range slices.Sorted(maps.Keys(env)) {
		value := env[key]

		if key == "" || strings.ContainsAny(key, "=\x00") {
			return fmt.Errorf("%w: key %q", errors.ErrInvalidEnv, key)
		}

		if strings.ContainsRune(value, 0) {
			return fmt.Errorf("%w: %s", errors.ErrInvalidEnv, key)
		}

		if _, ok := EnvReference(value); ok {
			continue
		}

		if LooksSecret(key) {
			return fmt.Errorf("%w: %s; use ${VAR}", errors.ErrSecretEnvValue, key)
		}
	}

	return nil
}

// EnvReference reports whether the whole value is a "${VAR}" reference.
//
// Only the exact form is a reference: a value that merely contains "${VAR}"
// stays literal, so a command cannot smuggle a resolved secret into a value the
// secret check would otherwise see as text.
func EnvReference(value string) (string, bool) {
	if !strings.HasPrefix(value, "${") || !strings.HasSuffix(value, "}") {
		return "", false
	}

	name := value[2 : len(value)-1]
	if name == "" {
		return "", false
	}

	for i, r := range name {
		switch {
		case r == '_', r >= 'A' && r <= 'Z', r >= 'a' && r <= 'z':
		case i > 0 && r >= '0' && r <= '9':
		default:
			return "", false
		}
	}

	return name, true
}

// LooksSecret reports whether a key names a credential by convention.
//
// It is a key-name heuristic, not a secret scanner: it catches the common
// names and fails safe on the false-positive side.
func LooksSecret(key string) bool {
	upper := strings.ToUpper(key)

	for _, marker := range []string{"TOKEN", "SECRET", "PASSWORD", "PASSWD", "CREDENTIAL"} {
		if strings.Contains(upper, marker) {
			return true
		}
	}

	return upper == "KEY" || strings.HasSuffix(upper, "_KEY")
}
