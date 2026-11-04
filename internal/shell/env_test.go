package shell

import (
	"errors"
	"testing"

	crerrors "github.com/fiwon123/crower/internal/errors"
)

// The rule the runner and the write path share: a literal under a
// secret-looking key is refused, a ${VAR} reference is the way to carry one,
// and an ordinary literal is fine.
func TestValidateEnv(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		wantErr error
	}{
		{"empty", nil, nil},
		{"ordinary literal", map[string]string{"STAGE": "prod"}, nil},
		{"reference", map[string]string{"TOKEN": "${TOKEN}"}, nil},
		{"token literal", map[string]string{"API_TOKEN": "literal"}, crerrors.ErrSecretEnvValue},
		{"password literal", map[string]string{"DB_PASSWORD": "literal"}, crerrors.ErrSecretEnvValue},
		{"key suffix", map[string]string{"API_KEY": "literal"}, crerrors.ErrSecretEnvValue},
		{"bare key", map[string]string{"KEY": "literal"}, crerrors.ErrSecretEnvValue},
		{"ordinary key with key in it", map[string]string{"KEYBOARD_LAYOUT": "us"}, nil},
		{"empty key", map[string]string{"": "x"}, crerrors.ErrInvalidEnv},
		{"equals in key", map[string]string{"A=B": "x"}, crerrors.ErrInvalidEnv},
		{"nul in value", map[string]string{"A": "x\x00"}, crerrors.ErrInvalidEnv},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnv(tt.env)

			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("ValidateEnv(%v) = %v, want nil", tt.env, err)
				}

				return
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ValidateEnv(%v) = %v, want %v", tt.env, err, tt.wantErr)
			}
		})
	}
}

// Only the exact "${VAR}" form is a reference.
func TestEnvReference(t *testing.T) {
	tests := []struct {
		value string
		name  string
		ok    bool
	}{
		{"${TOKEN}", "TOKEN", true},
		{"${_A1}", "_A1", true},
		{"${}", "", false},
		{"${1BAD}", "", false},
		{"${A-B}", "", false},
		{"prefix${TOKEN}", "", false},
		{"${TOKEN}suffix", "", false},
		{"plain", "", false},
	}

	for _, tt := range tests {
		name, ok := EnvReference(tt.value)
		if ok != tt.ok || name != tt.name {
			t.Errorf("EnvReference(%q) = %q, %v; want %q, %v", tt.value, name, ok, tt.name, tt.ok)
		}
	}
}

func TestLooksSecret(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"API_TOKEN", true},
		{"CLIENT_SECRET", true},
		{"DB_PASSWORD", true},
		{"PASSWD", true},
		{"AWS_CREDENTIAL", true},
		{"API_KEY", true},
		{"KEY", true},
		{"KEYBOARD_LAYOUT", false},
		{"STAGE", false},
		{"PATH", false},
	}

	for _, tt := range tests {
		if got := LooksSecret(tt.key); got != tt.want {
			t.Errorf("LooksSecret(%q) = %v, want %v", tt.key, got, tt.want)
		}
	}
}
