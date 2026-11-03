// Package github reads release metadata from the project's repository.
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// LatestReleaseURL is the endpoint that names the newest published release.
const LatestReleaseURL = "https://api.github.com/repos/fiwon123/crower/releases/latest"

// LatestTag returns the tag of the newest published release.
func LatestTag() (string, error) {
	return latestTagFrom(LatestReleaseURL)
}

// latestTagFrom is LatestTag against an arbitrary endpoint, so a test can
// answer without the network.
func latestTagFrom(url string) (string, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("check for a new release: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("check for a new release: %s", response.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(response.Body).Decode(&release); err != nil {
		return "", fmt.Errorf("read the latest release: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("read the latest release: no tag in the response")
	}

	return release.TagName, nil
}
