package checkdom

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiwon123/crower/internal/app"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) CheckNewVersion(currentVersion string) (string, error) {
	resp, err := http.Get("https://api.github.com/repos/fiwon123/crower/releases/latest")
	if err != nil {
		return "", fmt.Errorf("error: %v \n", err)
	}
	defer resp.Body.Close()

	var data struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("error: %v \n", err)
	}

	if data.TagName == currentVersion {
		return currentVersion, nil
	}

	return data.TagName, nil
}
