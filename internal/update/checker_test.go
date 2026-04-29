package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// roundTripperFunc adapts a function to the http.RoundTripper interface.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestCheckLatest_UpdateNeeded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v1.1.0",
			"assets": []any{
				map[string]any{"browser_download_url": "https://example.com/release.tar.gz"},
			},
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	origTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == LatestReleaseURL {
			testReq, _ := http.NewRequest(r.Method, ts.URL, r.Body)
			return origTransport.RoundTrip(testReq)
		}
		return origTransport.RoundTrip(r)
	})
	defer func() { http.DefaultTransport = origTransport }()

	latest, url, needsUpdate, err := CheckLatest("v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if latest != "v1.1.0" {
		t.Errorf("expected latest 'v1.1.0', got %q", latest)
	}
	if url != "https://example.com/release.tar.gz" {
		t.Errorf("expected download URL 'https://example.com/release.tar.gz', got %q", url)
	}
	if !needsUpdate {
		t.Error("expected needsUpdate to be true")
	}
}

func TestCheckLatest_NoUpdateNeeded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]any{
			"tag_name": "v1.0.0",
			"assets":   []any{},
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	origTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == LatestReleaseURL {
			testReq, _ := http.NewRequest(r.Method, ts.URL, r.Body)
			return origTransport.RoundTrip(testReq)
		}
		return origTransport.RoundTrip(r)
	})
	defer func() { http.DefaultTransport = origTransport }()

	latest, url, needsUpdate, err := CheckLatest("v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if latest != "v1.0.0" {
		t.Errorf("expected latest 'v1.0.0', got %q", latest)
	}
	if url != "" {
		t.Errorf("expected empty download URL, got %q", url)
	}
	if needsUpdate {
		t.Error("expected needsUpdate to be false when versions match")
	}
}

func TestCheckLatest_NetworkError(t *testing.T) {
	origTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("simulated network error")
	})
	defer func() { http.DefaultTransport = origTransport }()

	_, _, _, err := CheckLatest("v1.0.0")
	if err == nil {
		t.Error("expected error for network failure, got nil")
	}
}

func TestCheckLatest_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	origTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == LatestReleaseURL {
			testReq, _ := http.NewRequest(r.Method, ts.URL, r.Body)
			return origTransport.RoundTrip(testReq)
		}
		return origTransport.RoundTrip(r)
	})
	defer func() { http.DefaultTransport = origTransport }()

	_, _, _, err := CheckLatest("v1.0.0")
	if err == nil {
		t.Error("expected error for non-OK status, got nil")
	}
}

func TestCheckLatest_MissingTagName(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string]any{
			"message": "not found",
		}); err != nil {
			t.Fatal(err)
		}
	}))
	defer ts.Close()

	origTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.String() == LatestReleaseURL {
			testReq, _ := http.NewRequest(r.Method, ts.URL, r.Body)
			return origTransport.RoundTrip(testReq)
		}
		return origTransport.RoundTrip(r)
	})
	defer func() { http.DefaultTransport = origTransport }()

	_, _, _, err := CheckLatest("v1.0.0")
	if err == nil {
		t.Error("expected error for missing tag_name, got nil")
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1, v2 string
		want   int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.1.0", "1.0.0", 1},
		{"1.0.0", "1.1.0", -1},
		{"2.0.0", "1.9.9", 1},
		{"0.9.0", "1.0.0", -1},
	}

	for _, tc := range tests {
		got := compareVersions(tc.v1, tc.v2)
		if got != tc.want {
			t.Errorf("compareVersions(%q, %q) = %d, want %d", tc.v1, tc.v2, got, tc.want)
		}
	}
}

func TestFindDownloadURL(t *testing.T) {
	// No assets key
	url := findDownloadURL(map[string]any{})
	if url != "" {
		t.Errorf("expected empty URL when no assets, got %q", url)
	}

	// Empty assets array
	url = findDownloadURL(map[string]any{
		"assets": []any{},
	})
	if url != "" {
		t.Errorf("expected empty URL for empty assets, got %q", url)
	}

	// With valid asset
	url = findDownloadURL(map[string]any{
		"assets": []any{
			map[string]any{"browser_download_url": "https://example.com/release.tar.gz"},
		},
	})
	if url != "https://example.com/release.tar.gz" {
		t.Errorf("expected 'https://example.com/release.tar.gz', got %q", url)
	}

	// Asset without download URL
	url = findDownloadURL(map[string]any{
		"assets": []any{
			map[string]any{"name": "readme.txt"},
		},
	})
	if url != "" {
		t.Errorf("expected empty URL for asset without browser_download_url, got %q", url)
	}

	// Skip invalid asset, use valid one
	url = findDownloadURL(map[string]any{
		"assets": []any{
			"not a map",
			map[string]any{"browser_download_url": "https://example.com/valid.tar.gz"},
		},
	})
	if url != "https://example.com/valid.tar.gz" {
		t.Errorf("expected 'https://example.com/valid.tar.gz', got %q", url)
	}
}
