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
	return LatestTagTimeout(30 * time.Second)
}

// LatestTagTimeout is LatestTag with a shorter deadline, for the periodic
// notice: a check that delays every command must not wait half a minute.
func LatestTagTimeout(timeout time.Duration) (string, error) {
	return latestTagFromWithin(LatestReleaseURL, timeout)
}

// latestTagFrom is LatestTag against an arbitrary endpoint, so a test can
// answer without the network.
func latestTagFrom(url string) (string, error) {
	return latestTagFromWithin(url, 30*time.Second)
}

func latestTagFromWithin(url string, timeout time.Duration) (string, error) {
	client := &http.Client{Timeout: timeout}

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
