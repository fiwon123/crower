package upgradedom

import (
	"net/http"
	"runtime"

	"github.com/fiwon123/crower/internal/app"
	"github.com/minio/selfupdate"
)

type Handler struct {
	app *app.Config
}

func NewHandler(app *app.Config) *Handler {
	return &Handler{
		app: app,
	}
}

func (h *Handler) UpgradeApp(newVersion string) error {
	var resp *http.Response
	var err error
	switch runtime.GOOS {
	case "windows":
		resp, err = http.Get(
			"https://github.com/fiwon123/crower/releases/latest/download/crower_" + newVersion + "_windows.zip",
		)
	case "linux":
		resp, err = http.Get(
			"https://github.com/fiwon123/crower/releases/latest/download/crower_" + newVersion + "_linux.tar.gz",
		)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return selfupdate.Apply(resp.Body, selfupdate.Options{})

}
