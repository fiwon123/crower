package extractdom

import (
	"fmt"

	"github.com/fiwon123/crower/internal/errors"
)

// extractCommand builds the platform's invocation for one archive.
//
// Both the platform and the extension are parameters so a table test can cover
// every combination the release can meet: runtime.GOOS cannot be overridden at
// run time, and `.zip` had a linux arm but no darwin one, so macOS reported
// "no extractor" for a format `unzip` handles (#129).
func extractCommand(goos string, ext string, filePath string, outDir string) (string, error) {
	switch goos {
	case "windows", "linux", "darwin":
	default:
		return "", fmt.Errorf("%w: extract on %q", errors.ErrUnsupportedOS, goos)
	}

	switch ext {
	case "tar":
		return fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir), nil
	case "gz":
		return fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir), nil
	case "tgz":
		return fmt.Sprintf(`tar -xzf '%s' -C '%s'`, filePath, outDir), nil
	case "bz2":
		return fmt.Sprintf(`tar -xjf '%s' -C '%s'`, filePath, outDir), nil
	case "xz":
		return fmt.Sprintf(`tar -xJf '%s' -C '%s'`, filePath, outDir), nil
	case "zip":
		switch goos {
		case "windows":
			// Windows' tar reads zip archives; unzip is not part of the OS.
			return fmt.Sprintf(`tar -xf '%s' -C '%s'`, filePath, outDir), nil
		default:
			// linux and darwin both ship unzip.
			return fmt.Sprintf(`unzip '%s' -d  '%s'`, filePath, outDir), nil
		}
	case "7z", "rar":
		return fmt.Sprintf(`7z x '%s' '-o%s'`, filePath, outDir), nil
	default:
		return "", fmt.Errorf("no extractor for .%s archives", ext)
	}
}
