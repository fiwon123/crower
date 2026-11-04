package github_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/fiwon123/crower/internal/adapters/github"
	crowererrors "github.com/fiwon123/crower/internal/errors"
)

// registryServer serves an index and one pack.
func registryServer(t *testing.T, index string, pack string) *httptest.Server {
	t.Helper()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/index.json":
			io.WriteString(w, index)
		case "/packs/tools.toml":
			io.WriteString(w, pack)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)

	return server
}

func TestRegistryIndex(t *testing.T) {
	server := registryServer(t, `{"packs":[{"name":"tools","description":"Useful","version":"1.0.0","tags":["cli"]}]}`, "")

	packs, err := github.NewRegistry().WithBase(server.URL, server.Client()).Index()
	if err != nil {
		t.Fatalf("Index: %v", err)
	}

	if len(packs) != 1 || packs[0].Name != "tools" || packs[0].Version != "1.0.0" {
		t.Errorf("packs = %+v, want the index entry", packs)
	}
}

func TestRegistryPack(t *testing.T) {
	server := registryServer(t, "", `[greet]
Name = "greet"
Exec = "echo hi"
AllAlias = ["hi"]
`)

	pack, url, err := github.NewRegistry().WithBase(server.URL, server.Client()).Pack("tools")
	if err != nil {
		t.Fatalf("Pack: %v", err)
	}

	if !strings.HasSuffix(url, "/packs/tools.toml") {
		t.Errorf("url = %q, want the pack URL", url)
	}

	if len(pack) != 1 || pack[0].Name != "greet" || !slices.Equal(pack[0].Aliases, []string{"hi"}) {
		t.Errorf("pack = %+v, want greet with its alias", pack)
	}
}

func TestRegistryRejectsAnInvalidPackName(t *testing.T) {
	for _, name := range []string{"", ".", "..", "../etc", "a/b", "a b"} {
		t.Run(fmt.Sprintf("%q", name), func(t *testing.T) {
			_, _, err := github.NewRegistry().Pack(name)
			if !errors.Is(err, crowererrors.ErrInvalidPackName) {
				t.Fatalf("error = %v, want ErrInvalidPackName", err)
			}
		})
	}
}

func TestRegistryRefusesAPlaintextBase(t *testing.T) {
	registry := github.NewRegistry().WithBase("http://example.com", nil)

	if _, err := registry.Index(); err == nil || !strings.Contains(err.Error(), "only HTTPS") {
		t.Fatalf("error = %v, want the HTTPS refusal", err)
	}
}

func TestRegistryReportsMissingFiles(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	t.Cleanup(server.Close)

	registry := github.NewRegistry().WithBase(server.URL, server.Client())

	if _, err := registry.Index(); err == nil || !strings.Contains(err.Error(), "failed to fetch registry: 404") {
		t.Errorf("index error = %v, want the 404 message", err)
	}

	if _, _, err := registry.Pack("missing"); err == nil || !strings.Contains(err.Error(), "pack not found: 404") {
		t.Errorf("pack error = %v, want the 404 message", err)
	}
}

func TestRegistryRefusesAnEmptyPack(t *testing.T) {
	empty := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "")
	}))
	t.Cleanup(empty.Close)

	_, _, err := github.NewRegistry().WithBase(empty.URL, empty.Client()).Pack("tools")
	if err == nil || !strings.Contains(err.Error(), `pack "tools" is empty`) {
		t.Fatalf("error = %v, want the empty-pack message", err)
	}
}
