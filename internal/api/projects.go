package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListProjects gets all projects for the authenticated user
func (c *Client) ListProjects() ([]Project, error) {
	resp, err := c.makeRequest("GET", "/api/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list projects failed (status %d): %s", resp.StatusCode, string(body))
	}

	var projects []Project
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProject gets a specific project by ID
func (c *Client) GetProject(projectID string) (*Project, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/projects/%s", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("project not found")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get project failed (status %d): %s", resp.StatusCode, string(body))
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(req CreateProjectRequest) (*Project, error) {
	resp, err := c.makeRequest("POST", "/api/projects", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("create project failed (status %d): %s", resp.StatusCode, string(body))
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

// DeleteProject deletes a project
func (c *Client) DeleteProject(projectID string) error {
	resp, err := c.makeRequest("DELETE", fmt.Sprintf("/api/projects/%s", projectID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("project not found")
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete project failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetProjectBuilds gets all builds for a project
func (c *Client) GetProjectBuilds(projectID string) ([]Build, error) {
	resp, err := c.makeRequest("GET", fmt.Sprintf("/api/projects/%s/builds", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get project builds failed (status %d): %s", resp.StatusCode, string(body))
	}

	var builds []Build
	if err := json.NewDecoder(resp.Body).Decode(&builds); err != nil {
		return nil, err
	}

	return builds, nil
}

// StartBuild starts a new build for a project
func (c *Client) StartBuild(projectID string) (*Build, error) {
	resp, err := c.makeRequest("POST", fmt.Sprintf("/api/projects/%s/build", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("start build failed (status %d): %s", resp.StatusCode, string(body))
	}

	var build Build
	if err := json.NewDecoder(resp.Body).Decode(&build); err != nil {
		return nil, err
	}

	return &build, nil
}
