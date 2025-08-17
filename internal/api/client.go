package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

// DeployAndStream starts an end-to-end deployment with streaming progress updates
func (c *Client) DeployAndStream(request *DeployStreamRequest, updateHandler func(*StreamUpdate) error) error {
	// Create HTTP request for streaming endpoint
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api-key/end-to-end/deploy-stream", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers for streaming
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	// Create client with no timeout for streaming
	streamClient := &http.Client{
		Timeout: 0, // No timeout for streaming
	}

	resp, err := streamClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to start deployment stream: %w", err)
	}
	defer resp.Body.Close()

	// Check for success status codes (200 OK or 201 Created)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("deployment failed (status %d): %s", resp.StatusCode, string(body))
	}

	// Process Server-Sent Events stream
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// Parse Server-Sent Events format
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Skip empty data lines
			if data == "" || data == "[DONE]" {
				continue
			}

			// Parse JSON data
			var update StreamUpdate
			if err := json.Unmarshal([]byte(data), &update); err != nil {
				// If JSON parsing fails, treat as a simple message
				update = StreamUpdate{
					Type:    "log",
					Message: data,
				}
			}

			// Call the update handler
			if err := updateHandler(&update); err != nil {
				return err
			}

			// Stop processing if deployment is complete or failed
			if update.Type == "complete" || update.Type == "error" ||
				update.CurrentStep == "COMPLETED" || update.CurrentStep == "FAILED" ||
				update.BuildStatus == "failed" {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
