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
// The naming is the historical one (no architecture segment), and the caller
// hands the stream straight to selfupdate.Apply, which writes it over the
// running binary. The release assets are archives, so an arm for a new platform
// cannot be added safely until the download is unpacked first (#143). That is
// why darwin -- a platform the release ships -- still returns an error here
// instead of an asset name: the old code panicked on it, and naming an archive
// for it would have turned a crash into a corrupted installation.
func upgradeAsset(goos string, version string) (string, error) {
	switch goos {
	case "windows":
		return "crower_" + version + "_windows.zip", nil
	case "linux":
		return "crower_" + version + "_linux.tar.gz", nil
	default:
		return "", fmt.Errorf("%w: upgrade on %q", errors.ErrUnsupportedOS, goos)
	}
}
