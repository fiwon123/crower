package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	crowererrors "github.com/fiwon123/crower/internal/errors"
	"github.com/minio/selfupdate"
)

func TestAssetName(t *testing.T) {
	tests := []struct {
		name    string
		goos    string
		goarch  string
		version string
		want    string
		wantErr bool
	}{
		{"windows amd64", "windows", "amd64", "v1.2.3", "crower_v1.2.3_windows_amd64.zip", false},
		{"windows arm64", "windows", "arm64", "v1.2.3", "crower_v1.2.3_windows_arm64.zip", false},
		{"linux amd64", "linux", "amd64", "v1.2.3", "crower_v1.2.3_linux_amd64.tar.gz", false},
		{"linux arm64", "linux", "arm64", "v1.2.3", "crower_v1.2.3_linux_arm64.tar.gz", false},
		{"darwin amd64", "darwin", "amd64", "v1.2.3", "crower_v1.2.3_darwin_amd64.tar.gz", false},
		{"darwin arm64", "darwin", "arm64", "v1.2.3", "crower_v1.2.3_darwin_arm64.tar.gz", false},
		{"unsupported platform", "plan9", "amd64", "v1.2.3", "", true},
		{"unsupported architecture", "linux", "386", "v1.2.3", "", true},
		{"empty version", "linux", "amd64", "", "", true},
		{"version without a patch", "linux", "amd64", "v1.2", "", true},
		{"version with a traversal", "linux", "amd64", "v1.2.3/../x", "", true},
		{"version with a query", "linux", "amd64", "v1.2.3?x", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AssetName(tt.goos, tt.goarch, tt.version)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("AssetName(%q, %q, %q) returned no error", tt.goos, tt.goarch, tt.version)
				}

				if tt.goos == "plan9" || tt.goarch == "386" {
					if !errors.Is(err, crowererrors.ErrUnsupportedOS) {
						t.Errorf("error = %v, want ErrUnsupportedOS", err)
					}
				}

				return
			}

			if err != nil {
				t.Fatalf("AssetName error = %v", err)
			}

			if got != tt.want {
				t.Errorf("asset = %q, want %q", got, tt.want)
			}
		})
	}
}

func tarGzArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var archive bytes.Buffer
	gz := gzip.NewWriter(&archive)
	tarball := tar.NewWriter(gz)

	for name, contents := range files {
		header := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(contents))}
		if err := tarball.WriteHeader(header); err != nil {
			t.Fatalf("write header: %v", err)
		}

		if _, err := tarball.Write([]byte(contents)); err != nil {
			t.Fatalf("write payload: %v", err)
		}
	}

	if err := tarball.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}

	if err := gz.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}

	return archive.Bytes()
}

func zipArchive(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var archive bytes.Buffer
	writer := zip.NewWriter(&archive)

	for name, contents := range files {
		file, err := writer.Create(name)
		if err != nil {
			t.Fatalf("create %s: %v", name, err)
		}

		if _, err := file.Write([]byte(contents)); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}

	return archive.Bytes()
}

// checksumLine is the shape `sha256sum *` writes for one file.
func checksumLine(body []byte, asset string) string {
	sum := sha256.Sum256(body)

	return fmt.Sprintf("%x  %s\n", sum, asset)
}

func TestParseChecksum(t *testing.T) {
	hash := strings.Repeat("a", 64)

	asset := "crower_v1_linux_amd64.tar.gz"

	tests := []struct {
		name    string
		body    string
		asset   string
		want    string
		wantErr bool
	}{
		{"text mode", hash + "  " + asset + "\n", asset, hash, false},
		{"binary mode", hash + " *" + asset + "\n", asset, hash, false},
		{"CRLF", hash + "  " + asset + "\r\n", asset, hash, false},
		{"tabs and extra spaces", hash + "\t\t" + asset + "  \n", asset, hash, false},
		{"no trailing newline", hash + "  " + asset, asset, hash, false},
		{"another asset first", strings.Repeat("b", 64) + "  other.tar.gz\n" + hash + "  " + asset + "\n", asset, hash, false},
		{"a superstring is not the asset", strings.Repeat("b", 64) + "  x" + asset + ".bak\n", asset, "", true},
		{"missing entry", hash + "  other.tar.gz\n", asset, "", true},
		{"malformed hash", "not-a-hash  " + asset + "\n", asset, "", true},
		{"empty", "", asset, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseChecksum([]byte(tt.body), tt.asset)

			if tt.wantErr {
				if !errors.Is(err, crowererrors.ErrChecksumUnavailable) {
					t.Errorf("error = %v, want ErrChecksumUnavailable", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("parseChecksum() error = %v", err)
			}

			if got == nil || fmt.Sprintf("%x", got) != tt.want {
				t.Errorf("parseChecksum() = %x, want %s", got, tt.want)
			}
		})
	}
}

func TestBinaryFromArchive(t *testing.T) {
	binary := "#!/bin/sh\necho crower\n"

	t.Run("tar.gz", func(t *testing.T) {
		body := tarGzArchive(t, map[string]string{
			"README.md": "docs",
			"crower":    binary,
		})

		got, err := binaryFromArchive("crower_v1_linux_amd64.tar.gz", body)
		if err != nil {
			t.Fatalf("binaryFromArchive: %v", err)
		}

		if string(got) != binary {
			t.Errorf("binary = %q, want %q", got, binary)
		}
	})

	t.Run("zip", func(t *testing.T) {
		body := zipArchive(t, map[string]string{
			"LICENSE":    "docs",
			"crower.exe": binary,
		})

		got, err := binaryFromArchive("crower_v1_windows_amd64.zip", body)
		if err != nil {
			t.Fatalf("binaryFromArchive: %v", err)
		}

		if string(got) != binary {
			t.Errorf("binary = %q, want %q", got, binary)
		}
	})

	t.Run("empty binary", func(t *testing.T) {
		body := tarGzArchive(t, map[string]string{"crower": ""})

		if _, err := binaryFromArchive("crower_v1_linux_amd64.tar.gz", body); err == nil {
			t.Fatal("an empty crower returned no error")
		}
	})

	t.Run("missing binary", func(t *testing.T) {
		body := tarGzArchive(t, map[string]string{"README.md": "docs"})

		if _, err := binaryFromArchive("crower_v1_linux_amd64.tar.gz", body); err == nil {
			t.Fatal("an archive without crower returned no error")
		}
	})

	t.Run("corrupt archive", func(t *testing.T) {
		if _, err := binaryFromArchive("crower_v1_linux_amd64.tar.gz", []byte("not a tarball")); err == nil {
			t.Fatal("a corrupt archive returned no error")
		}
	})

	t.Run("unknown archive", func(t *testing.T) {
		if _, err := binaryFromArchive("crower_v1_linux_amd64.rpm", nil); err == nil {
			t.Fatal("an unknown archive returned no error")
		}
	})
}

func TestFetchBinaryChecksTheStatus(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.NotFound(w, nil)
	}))
	defer server.Close()

	handler := New().WithRelease(server.URL, server.Client())

	_, err := handler.fetchBinary("v1.2.3", "crower_v1.2.3_linux_amd64.tar.gz")
	if err == nil {
		t.Fatal("a 404 returned no error")
	}

	if !strings.Contains(err.Error(), "404") || !strings.Contains(err.Error(), server.URL) {
		t.Errorf("error = %v, want the URL and the status", err)
	}
}

func TestFetchBinaryReturnsTheExecutable(t *testing.T) {
	binary := "binary"
	body := tarGzArchive(t, map[string]string{"crower": binary})

	asset := "crower_v1.2.3_linux_amd64.tar.gz"

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1.2.3/" + asset:
			w.Write(body)
		case "/v1.2.3/checksums.txt":
			io.WriteString(w, checksumLine(body, asset))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	handler := New().WithRelease(server.URL, server.Client())

	got, err := handler.fetchBinary("v1.2.3", "crower_v1.2.3_linux_amd64.tar.gz")
	if err != nil {
		t.Fatalf("fetchBinary: %v", err)
	}

	if string(got) != binary {
		t.Errorf("binary = %q, want %q", got, binary)
	}
}

func TestHTTPSOnlyRedirect(t *testing.T) {
	plain, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	if httpsOnlyRedirect(plain, nil) == nil {
		t.Error("a redirect to http was accepted")
	}

	secure, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatalf("request: %v", err)
	}

	if err := httpsOnlyRedirect(secure, nil); err != nil {
		t.Errorf("an https redirect was refused: %v", err)
	}

	if err := httpsOnlyRedirect(secure, make([]*http.Request, 10)); err == nil {
		t.Error("eleven redirects were accepted")
	}
}

// The upgrade must hand the extracted executable, not the archive, to
// selfupdate. The applier is replaced so the test binary survives.
func TestUpgradeAppAppliesTheExtractedBinary(t *testing.T) {
	asset, err := AssetName(runtime.GOOS, runtime.GOARCH, "v1.2.3")
	if err != nil {
		t.Skipf("no upgrade asset for this platform: %v", err)
	}

	const binary = "the crower binary"

	var body []byte
	if strings.HasSuffix(asset, ".zip") {
		body = zipArchive(t, map[string]string{"crower.exe": binary})
	} else {
		body = tarGzArchive(t, map[string]string{"crower": binary})
	}

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1.2.3/" + asset:
			w.Write(body)
		case "/v1.2.3/checksums.txt":
			io.WriteString(w, checksumLine(body, asset))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	oldApply := applyUpdate
	defer func() { applyUpdate = oldApply }()

	var applied []byte

	applyUpdate = func(reader io.Reader, _ selfupdate.Options) error {
		body, err := io.ReadAll(reader)
		applied = body
		return err
	}

	handler := New().WithRelease(server.URL, server.Client())

	if err := handler.Upgrade("v1.2.3"); err != nil {
		t.Fatalf("UpgradeApp: %v", err)
	}

	if string(applied) != binary {
		t.Errorf("applied %q, want the extracted binary", applied)
	}
}

// A failed download must not reach the applier.
func TestUpgradeAppDoesNotApplyAFailedDownload(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer server.Close()

	oldApply := applyUpdate
	defer func() { applyUpdate = oldApply }()

	applied := false
	applyUpdate = func(io.Reader, selfupdate.Options) error {
		applied = true
		return nil
	}

	handler := New().WithRelease(server.URL, server.Client())

	if err := handler.Upgrade("v1.2.3"); err == nil {
		t.Fatal("a 404 returned no error")
	}

	if applied {
		t.Error("the applier ran after a failed download")
	}
}

func TestFetchBinaryRefusesAPlaintextBase(t *testing.T) {
	handler := New().WithRelease("http://example.com", nil)

	_, err := handler.fetchBinary("v1.2.3", "crower_v1.2.3_linux_amd64.tar.gz")
	if err == nil || !strings.Contains(err.Error(), "only HTTPS") {
		t.Errorf("error = %v, want the HTTPS refusal", err)
	}
}

func TestFetchBinaryRejectsATamperedArchive(t *testing.T) {
	asset := "crower_v1.2.3_linux_amd64.tar.gz"
	body := tarGzArchive(t, map[string]string{"crower": "binary"})

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1.2.3/" + asset:
			w.Write(body)
		case "/v1.2.3/checksums.txt":
			io.WriteString(w, strings.Repeat("0", 64)+"  "+asset+"\n")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	handler := New().WithRelease(server.URL, server.Client())

	_, err := handler.fetchBinary("v1.2.3", asset)
	if !errors.Is(err, crowererrors.ErrChecksumMismatch) {
		t.Errorf("error = %v, want ErrChecksumMismatch", err)
	}
}

func TestFetchBinaryRequiresAChecksumEntry(t *testing.T) {
	asset := "crower_v1.2.3_linux_amd64.tar.gz"
	body := tarGzArchive(t, map[string]string{"crower": "binary"})

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1.2.3/" + asset:
			w.Write(body)
		case "/v1.2.3/checksums.txt":
			io.WriteString(w, strings.Repeat("0", 64)+"  another.tar.gz\n")
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	handler := New().WithRelease(server.URL, server.Client())

	_, err := handler.fetchBinary("v1.2.3", asset)
	if !errors.Is(err, crowererrors.ErrChecksumUnavailable) {
		t.Errorf("error = %v, want ErrChecksumUnavailable", err)
	}
}

// A release without checksums.txt must fail before the archive is fetched.
func TestFetchBinaryRequiresChecksums(t *testing.T) {
	asset := "crower_v1.2.3_linux_amd64.tar.gz"
	body := tarGzArchive(t, map[string]string{"crower": "binary"})
	archiveFetched := false

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1.2.3/" + asset:
			archiveFetched = true
			w.Write(body)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	handler := New().WithRelease(server.URL, server.Client())

	_, err := handler.fetchBinary("v1.2.3", asset)
	if err == nil {
		t.Fatal("a release without checksums.txt returned no error")
	}

	if archiveFetched {
		t.Error("the archive was downloaded before the checksum was known to exist")
	}
}
