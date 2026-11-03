package sharedom

import "testing"

// A pack name is a slug, not a path: it is interpolated into the URL and used
// to derive a file name, so a separator or ".." must be rejected.
func TestValidPackName(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"good-pack", false},
		{"docker-setup_v2", false},
		{"v1.0", false},
		{"", true},
		{".", true},
		{"..", true},
		{"../escape", true},
		{"a/b", true},
		{`a\b`, true},
		{"/etc/passwd", true},
		{"good?x", true},
		{"a b", true},
		{"%2e%2e", true},
		{"café", true},
		{"a／b", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validPackName(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("validPackName(%q) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
