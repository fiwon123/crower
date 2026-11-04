package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/fiwon123/crower/internal/errors"
	"github.com/minio/selfupdate"
)

const defaultReleaseBase = "https://github.com/fiwon123/crower/releases/download"

// The download and the extracted member are read into memory; the caps stop a
// corrupt or hostile response from exhausting the machine.
const (
	maxDownloadBytes = 128 << 20
	maxBinaryBytes   = 512 << 20
	maxChecksumBytes = 1 << 20
)

// The release workflow accepts only this tag shape.
var versionPattern = regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+$`)

// applyUpdate is selfupdate.Apply behind a variable so a test can observe the
// bytes the upgrade would install without replacing the test binary.
var applyUpdate = selfupdate.Apply

// Updater downloads, verifies and applies a release archive.
type Updater struct {
	baseURL string
	client  *http.Client
	goos    string
	goarch  string
}

// New builds an Updater for the host.
func New() *Updater {
	return &Updater{
		baseURL: defaultReleaseBase,
		goos:    runtime.GOOS,
		goarch:  runtime.GOARCH,
		client: &http.Client{
			Timeout:       60 * time.Second,
			CheckRedirect: httpsOnlyRedirect,
		},
	}
}

// WithRelease points the handler at another release base and HTTP client. A
// test uses it with an httptest server; production keeps the default.
func (h *Updater) WithRelease(base string, client *http.Client) *Updater {
	h.baseURL = base

	if client != nil {
		h.client = client
	}

	return h
}

// AssetName names the release archive for a platform, matching
// .github/workflows/release.yml (and the Makefile): the platform and
// architecture belong in the filename, and the archive is a tar.gz everywhere
// except Windows.
func AssetName(goos string, goarch string, version string) (string, error) {
	if !versionPattern.MatchString(version) {
		return "", fmt.Errorf("invalid release version: %q", version)
	}

	if goarch != "amd64" && goarch != "arm64" {
		return "", fmt.Errorf("%w: no release build for %s/%s", errors.ErrUnsupportedOS, goos, goarch)
	}

	switch goos {
	case "windows":
		return fmt.Sprintf("crower_%s_%s_%s.zip", version, goos, goarch), nil
	case "linux", "darwin":
		return fmt.Sprintf("crower_%s_%s_%s.tar.gz", version, goos, goarch), nil
	default:
		return "", fmt.Errorf("%w: no release build for %s/%s", errors.ErrUnsupportedOS, goos, goarch)
	}
}

// Upgrade downloads the named release for the host, verifies it and applies it.
func (h *Updater) Upgrade(newVersion string) error {
	asset, err := AssetName(h.goos, h.goarch, newVersion)
	if err != nil {
		return err
	}

	binary, err := h.fetchBinary(newVersion, asset)
	if err != nil {
		return err
	}

	return applyUpdate(bytes.NewReader(binary), selfupdate.Options{})
}

// fetchBinary downloads the asset, verifies it against the release's
// checksums.txt, and returns the executable inside it. The status is checked
// before the body is read: a 404 is a response, not an error, and used to be
// handed to selfupdate as if it were a binary.
func (h *Updater) fetchBinary(version string, asset string) ([]byte, error) {
	url, err := h.releaseURL(version, asset)
	if err != nil {
		return nil, err
	}

	// Fetch the small checksum first: a release without one fails before the
	// archive is downloaded.
	expected, err := h.fetchChecksum(version, asset)
	if err != nil {
		return nil, err
	}

	body, err := h.download(url, maxDownloadBytes)
	if err != nil {
		return nil, err
	}

	actual := sha256.Sum256(body)
	if !bytes.Equal(actual[:], expected) {
		return nil, fmt.Errorf("%w for %s", errors.ErrChecksumMismatch, asset)
	}

	return binaryFromArchive(asset, body)
}

// releaseURL builds a release download URL, refusing anything but HTTPS.
func (h *Updater) releaseURL(version string, file string) (string, error) {
	if !strings.HasPrefix(h.baseURL, "https://") {
		return "", fmt.Errorf("only HTTPS allowed for security")
	}

	return h.baseURL + "/" + version + "/" + file, nil
}

func (h *Updater) download(url string, limit int64) ([]byte, error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: %s", url, resp.Status)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, limit+1))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", url, err)
	}

	if int64(len(body)) > limit {
		return nil, fmt.Errorf("read %s: larger than %d bytes", url, limit)
	}

	return body, nil
}

// fetchChecksum reads the SHA-256 of one asset from checksums.txt. A missing or
// malformed entry is an error: skipping the check would silently accept a
// tampered archive.
func (h *Updater) fetchChecksum(version string, asset string) ([]byte, error) {
	url, err := h.releaseURL(version, "checksums.txt")
	if err != nil {
		return nil, err
	}

	body, err := h.download(url, maxChecksumBytes)
	if err != nil {
		return nil, err
	}

	return parseChecksum(body, asset)
}

// parseChecksum reads one `sha256sum` line. The name may carry the "*" that
// sha256sum writes in binary mode.
func parseChecksum(body []byte, asset string) ([]byte, error) {
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}

		if strings.TrimPrefix(fields[1], "*") != asset {
			continue
		}

		sum, err := hex.DecodeString(fields[0])
		if err != nil || len(sum) != sha256.Size {
			return nil, fmt.Errorf("%w: checksums.txt has a malformed hash for %s", errors.ErrChecksumUnavailable, asset)
		}

		return sum, nil
	}

	return nil, fmt.Errorf("%w: checksums.txt has no entry for %s", errors.ErrChecksumUnavailable, asset)
}

// binaryFromArchive extracts the crower executable from a release archive.
func binaryFromArchive(asset string, body []byte) ([]byte, error) {
	if strings.HasSuffix(asset, ".zip") {
		return binaryFromZip(body)
	}

	if strings.HasSuffix(asset, ".tar.gz") {
		return binaryFromTarGz(body)
	}

	return nil, fmt.Errorf("unsupported archive %q", asset)
}

func binaryFromZip(body []byte) ([]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return nil, fmt.Errorf("read zip: %w", err)
	}

	for _, file := range reader.File {
		if filepath.Base(file.Name) != "crower.exe" {
			continue
		}

		opened, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("open %s: %w", file.Name, err)
		}
		defer opened.Close()

		return readBinary(opened, "crower.exe")
	}

	return nil, fmt.Errorf("zip carries no crower.exe")
}

func binaryFromTarGz(body []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("read gzip: %w", err)
	}
	defer gzipReader.Close()

	archive := tar.NewReader(gzipReader)

	for {
		header, err := archive.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("read tar: %w", err)
		}

		if header.Typeflag != tar.TypeReg || filepath.Base(header.Name) != "crower" {
			continue
		}

		return readBinary(archive, "crower")
	}

	return nil, fmt.Errorf("archive carries no crower")
}

// readBinary reads the executable and refuses an empty one: selfupdate would
// write a zero-byte file over the running binary.
func readBinary(reader io.Reader, name string) ([]byte, error) {
	binary, err := io.ReadAll(io.LimitReader(reader, maxBinaryBytes))
	if err != nil {
		return nil, err
	}

	if len(binary) == 0 {
		return nil, fmt.Errorf("%s is empty", name)
	}

	return binary, nil
}

// httpsOnlyRedirect refuses a redirect to plaintext and keeps the default
// ten-hop limit that a custom CheckRedirect would otherwise remove.
func httpsOnlyRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("stopped after %d redirects", len(via))
	}

	if req.URL.Scheme != "https" {
		return fmt.Errorf("refusing a redirect to %s", req.URL.Scheme)
	}

	return nil
}
