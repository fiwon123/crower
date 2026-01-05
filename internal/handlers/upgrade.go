package handlers

import (
	"net/http"
	"runtime"

	"github.com/fiwon123/crower/internal/data/app"
	"github.com/minio/selfupdate"
)

func UpgradeApp(newVersion string, app *app.Data) error {
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
