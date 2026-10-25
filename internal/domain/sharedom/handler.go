package sharedom

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fiwon123/crower/internal/app"
	command "github.com/fiwon123/crower/internal/data/command"
	"github.com/fiwon123/crower/internal/helper"
	"github.com/pelletier/go-toml"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{app: app}
}

type RegistryIndex struct {
	Packs []RegistryPack `json:"packs"`
}

type RegistryPack struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

func (h *Handler) Export(name string, outPath string) error {
	if name == "" {
		// Export all
		data := h.app.AllCommandsByName
		out, err := toml.Marshal(data)
		if err != nil {
			return err
		}
		if outPath == "" {
			outPath = "crower-export.toml"
		}
		return os.WriteFile(outPath, out, 0644)
	}

	cmd := h.app.AllCommandsByName.Get(name)
	if cmd == nil {
		cmd = h.app.AllCommandsByAlias.Get(name)
	}
	if cmd == nil {
		return fmt.Errorf("command not found")
	}

	data := map[string]*command.Data{name: cmd}
	out, err := toml.Marshal(data)
	if err != nil {
		return err
	}
	if outPath == "" {
		outPath = fmt.Sprintf("%s.toml", name)
	}
	return os.WriteFile(outPath, out, 0644)
}

func (h *Handler) Import(inPath string, overwrite bool) error {
	data, err := os.ReadFile(inPath)
	if err != nil {
		return err
	}

	var commands map[string]command.Data
	if err := toml.Unmarshal(data, &commands); err != nil {
		return err
	}

	for k, v := range commands {
		existing := h.app.AllCommandsByName.Get(k)
		if existing != nil && !overwrite {
			// Ask to rename or skip
			if !helper.GetUserConfirmation(fmt.Sprintf("Command '%s' exists. Rename?", k), h.app) {
				continue
			}
			v.Name = k + "_imported"
		}
		h.app.AllCommandsByName.Add(k, &v)
		h.app.Save()
	}
	return nil
}

func (h *Handler) getRegistryURL() string {
	if h.app.CfgFilePath != "" {
		// Try to read registry URL from config if stored
	}
	return "https://raw.githubusercontent.com/fiwon123/crower-resources/main"
}

func (h *Handler) ListRegistry() error {
	url := h.getRegistryURL() + "/index.json"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to fetch registry: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var idx RegistryIndex
	if err := json.Unmarshal(body, &idx); err != nil {
		return err
	}

	for _, p := range idx.Packs {
		fmt.Printf("%s - %s (v%s)\n", p.Name, p.Description, p.Version)
	}
	return nil
}

func (h *Handler) Pull(name string) error {
	if !strings.HasPrefix(h.getRegistryURL(), "https://") {
		return fmt.Errorf("only HTTPS allowed for security")
	}
	url := h.getRegistryURL() + "/packs/" + name + ".toml"
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("pack not found: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	
	// Save to temp file and import
	tmp := filepath.Join(os.TempDir(), name+".toml")
	if err := os.WriteFile(tmp, body, 0644); err != nil {
		return err
	}
	defer os.Remove(tmp)

	fmt.Printf("Importing pack '%s' from %s\n", name, url)
	fmt.Printf("Contents:\n%s\n\n", string(body))
	if !helper.GetUserConfirmation("Import?", h.app) {
		return nil
	}

	return h.Import(tmp, false)
}
