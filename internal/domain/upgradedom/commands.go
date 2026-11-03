package upgradedom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
)

// upgradeAsset names the release asset to download.
//
// goos is a parameter so a table test can cover every platform the release
// ships (#129). The old switch had no default and no darwin arm, so on macOS
// the response stayed nil and the caller dereferenced it.
//
// The naming is the historical one (no architecture segment). Changing the
// scheme is a separate fix: the caller hands the stream straight to
// selfupdate.Apply, which expects a binary, while these assets are archives.
func upgradeAsset(goos string, version string) (string, error) {
	switch goos {
	case "windows":
		return "crower_" + version + "_windows.zip", nil
	case "linux", "darwin":
		return "crower_" + version + "_" + goos + ".tar.gz", nil
	default:
		return "", fmt.Errorf("%w: upgrade on %q", errors.ErrUnsupportedOS, goos)
	}
}
