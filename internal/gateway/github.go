package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LatestReleaseResponse struct {
	TagName string               `json:"tag_name"`
	Assets  []LatestReleaseAsset `json:"assets"`
}

type LatestReleaseAsset struct {
	Url string `json:"browser_download_url"`
}

const (
	baseUrl           = "https://api.github.com"
	latestReleasePath = "repos/hpcsc/kz/releases/latest"
)

type Github interface {
	LatestRelease() (*LatestReleaseResponse, error)
	Download(url string) ([]byte, error)
}

var _ Github = (*github)(nil)

type github struct {
	httpClient http.Client
}

func NewGithubGateway() Github {
	return &github{
		httpClient: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (g *github) LatestRelease() (*LatestReleaseResponse, error) {
	url := fmt.Sprintf("%s/%s", baseUrl, latestReleasePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GET request to %s: %v", url, err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make GET request to %s, status code %d: %v", url, resp.StatusCode, err)
	}

	defer resp.Body.Close()
	var latestResponse LatestReleaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&latestResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &latestResponse, nil
}

func (g *github) Download(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new GET request to %s: %v", url, err)
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to %s: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to make GET request to %s, status code %d: %v", url, resp.StatusCode, err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return body, nil
}
