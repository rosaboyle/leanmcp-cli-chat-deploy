package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the API client
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(apiKey string) *Client {
	// Hardcoded base URL - never changes
	baseURL := "https://join-us.cracked-devs.link"

	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// makeRequest makes an HTTP request with authentication
func (c *Client) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

// TestConnection tests the API connection (simple endpoint)
func (c *Client) TestConnection() error {
	resp, err := c.makeRequest("GET", "/health", nil)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetAPIKeyInfo gets information about the current API key
func (c *Client) GetAPIKeyInfo() (*APIKeyInfo, error) {
	resp, err := c.makeRequest("GET", "/api/projects/api-key/info", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API key info failed (status %d): %s", resp.StatusCode, string(body))
	}

	var apiKeyInfo APIKeyInfo
	if err := json.NewDecoder(resp.Body).Decode(&apiKeyInfo); err != nil {
		return nil, err
	}

	return &apiKeyInfo, nil
}
