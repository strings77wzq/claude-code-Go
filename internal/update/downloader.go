package update

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadAndUpdate(downloadURL string, binaryPath string) error {
	if downloadURL == "" {
		return fmt.Errorf("download URL is empty")
	}
	if binaryPath == "" {
		return fmt.Errorf("binary path is empty")
	}

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected download status code: %d", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(binaryPath), "go-code-update-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
	}()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write temp binary: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp binary: %w", err)
	}

	if err := os.Chmod(tmpPath, 0o755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	if err := os.Rename(tmpPath, binaryPath); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	return nil
}
