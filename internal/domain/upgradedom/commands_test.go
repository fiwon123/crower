package upgradedom

import (
	"errors"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
)

func TestUpgradeAsset(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		want    string
		wantErr bool
	}{
		{"windows", "windows", "crower_v1.2.3_windows.zip", false},
		{"linux", "linux", "crower_v1.2.3_linux.tar.gz", false},
		{"darwin", "darwin", "crower_v1.2.3_darwin.tar.gz", false},
		{"unsupported platform", "plan9", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := upgradeAsset(tt.goos, "v1.2.3")

			if tt.wantErr {
				if !errors.Is(err, crowererrors.ErrUnsupportedOS) {
					t.Fatalf("error = %v, want ErrUnsupportedOS", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("error = %v", err)
			}

			if got != tt.want {
				t.Errorf("asset = %q, want %q", got, tt.want)
			}
		})
	}
}
