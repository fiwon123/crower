package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fiwon123/crower/internal/application/sharecommands"
	"github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/domain/catalog"
	"github.com/fiwon123/crower/internal/errors"
	"github.com/pelletier/go-toml"
)

const defaultRegistryURL = "https://raw.githubusercontent.com/fiwon123/crower-resources/main"

// Registry reads the shared command registry.
type Registry struct {
	baseURL string
	client  *http.Client
}

// NewRegistry builds a Registry for the published registry.
func NewRegistry() *Registry {
	return &Registry{
		baseURL: defaultRegistryURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
			// A registry may not redirect the download to plaintext.
			CheckRedirect: func(req *http.Request, _ []*http.Request) error {
				if req.URL.Scheme != "https" {
					return fmt.Errorf("refusing a redirect to %s", req.URL.Scheme)
				}

				return nil
			},
		},
	}
}

// WithBase points the registry at another base and HTTP client. A test uses it
// with an httptest server; production keeps the default.
func (r *Registry) WithBase(url string, client *http.Client) *Registry {
	r.baseURL = url

	if client != nil {
		r.client = client
	}

	return r
}

// Index returns the packs the registry publishes.
func (r *Registry) Index() ([]sharecommands.Pack, error) {
	base, err := r.base()
	if err != nil {
		return nil, err
	}

	body, status, err := r.fetch(base + "/index.json")
	if err != nil {
		return nil, err
	}

	if status != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch registry: %d", status)
	}

	var index struct {
		Packs []sharecommands.Pack `json:"packs"`
	}

	if err := json.Unmarshal(body, &index); err != nil {
		return nil, err
	}

	return index.Packs, nil
}

// Pack fetches one pack, returning its commands and the URL they came from.
func (r *Registry) Pack(name string) ([]catalog.Command, string, error) {
	if err := ValidPackName(name); err != nil {
		return nil, "", err
	}

	base, err := r.base()
	if err != nil {
		return nil, "", err
	}

	url := base + "/packs/" + name + ".toml"

	body, status, err := r.fetch(url)
	if err != nil {
		return nil, "", err
	}

	if status != http.StatusOK {
		return nil, "", fmt.Errorf("pack not found: %d", status)
	}

	var commands map[string]command.Data
	if err := toml.Unmarshal(body, &commands); err != nil {
		return nil, "", fmt.Errorf("read pack %q: %w", name, err)
	}

	if len(commands) == 0 {
		return nil, "", fmt.Errorf("pack %q is empty", name)
	}

	pack := make([]catalog.Command, 0, len(commands))

	for key, value := range commands {
		pack = append(pack, catalog.Command{
			Name:    key,
			Aliases: value.AllAlias,
			Exec:    value.Exec,
			Steps:   value.Steps,
		})
	}

	return pack, url, nil
}

// ValidPackName rejects anything that is a path rather than a pack slug.
//
// The name is interpolated into the pack URL and used to derive a file, so a
// separator or ".." must never reach either.
func ValidPackName(name string) error {
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("%w: %q", errors.ErrInvalidPackName, name)
	}

	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
		case r == '-', r == '_', r == '.':
		default:
			return fmt.Errorf("%w: %q", errors.ErrInvalidPackName, name)
		}
	}

	return nil
}

// base is the configured registry, refusing anything but HTTPS.
func (r *Registry) base() (string, error) {
	if !strings.HasPrefix(r.baseURL, "https://") {
		return "", fmt.Errorf("only HTTPS allowed for security")
	}

	return r.baseURL, nil
}

// fetch downloads one registry file and returns it with the status, which the
// callers phrase differently (a missing index is not a missing pack).
func (r *Registry) fetch(url string) ([]byte, int, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}
