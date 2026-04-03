// Package update provides functionality for checking and downloading updates.
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	// LatestReleaseURL is the GitHub API URL for the latest release.
	LatestReleaseURL = "https://api.github.com/repos/strings77wzq/claude-code-Go/releases/latest"
)

// ReleaseInfo represents the parsed GitHub release information.
type ReleaseInfo struct {
	TagName     string
	DownloadURL string
	NeedsUpdate bool
}

// CheckLatest checks if there's a newer version available.
// It returns the latest version, download URL, whether an update is needed, and any error.
func CheckLatest(currentVersion string) (latestVersion string, downloadURL string, needsUpdate bool, err error) {
	// Make HTTP request to GitHub API
	resp, err := http.Get(LatestReleaseURL)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var releaseData map[string]any
	if err := json.Unmarshal(body, &releaseData); err != nil {
		return "", "", false, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract tag_name
	tagName, ok := releaseData["tag_name"].(string)
	if !ok || tagName == "" {
		return "", "", false, fmt.Errorf("tag_name not found in response")
	}

	// Normalize version string (remove 'v' prefix for comparison)
	latestVer := strings.TrimPrefix(tagName, "v")
	currentVer := strings.TrimPrefix(currentVersion, "v")

	// Compare versions (simple string comparison for semver)
	needsUpdate = compareVersions(latestVer, currentVer) > 0

	// Find download URL from assets
	downloadURL = findDownloadURL(releaseData)

	return tagName, downloadURL, needsUpdate, nil
}

// findDownloadURL finds the first valid download URL from assets.
func findDownloadURL(releaseData map[string]any) string {
	assets, ok := releaseData["assets"].([]any)
	if !ok || len(assets) == 0 {
		return ""
	}

	for _, asset := range assets {
		assetMap, ok := asset.(map[string]any)
		if !ok {
			continue
		}

		url, ok := assetMap["browser_download_url"].(string)
		if ok && url != "" {
			return url
		}
	}

	return ""
}

// compareVersions compares two semver version strings.
// Returns: 1 if v1 > v2, 0 if v1 == v2, -1 if v1 < v2.
func compareVersions(v1, v2 string) int {
	// Simple string comparison for semver
	// This works for most cases where versions have same format
	if v1 > v2 {
		return 1
	}
	if v1 < v2 {
		return -1
	}
	return 0
}
