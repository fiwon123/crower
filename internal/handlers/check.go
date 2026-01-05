package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fiwon123/crower/internal/data/app"
)

func CheckNewVersion(currentVersion string, app *app.Data) error {
	resp, err := http.Get("https://api.github.com/repos/fiwon123/crower/releases/latest")
	if err != nil {
		return fmt.Errorf("error: %v \n", err)
	}
	defer resp.Body.Close()

	var data struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error: %v \n", err)
	}

	if data.TagName != currentVersion {
		fmt.Printf("Current Version: %s \n", currentVersion)
		fmt.Printf("New Version Found: %s \n", data.TagName)
		fmt.Println("Check: https://github.com/fiwon123/crower/releases/latest")
	} else {
		fmt.Printf("Current Version: %s \n", currentVersion)
		fmt.Println("already up-to-date")
	}

	return nil
}
