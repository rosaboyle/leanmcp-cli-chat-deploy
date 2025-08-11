package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	
	"github.com/ddod/leanmcp-cli/internal/filesystem"
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

// GetUploadURL gets a pre-signed URL for uploading project files
func (c *Client) GetUploadURL(projectID, fileName string, fileSize int64) (*UploadURLResponse, error) {
	req := UploadURLRequest{
		FileName: fileName,
		FileType: "application/zip",
		FileSize: fileSize,
	}
	
	resp, err := c.makeRequest("POST", fmt.Sprintf("/api/projects/%s/upload-url", projectID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get upload URL failed (status %d): %s", resp.StatusCode, string(body))
	}
	
	var uploadResp UploadURLResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, err
	}
	
	return &uploadResp, nil
}

// UploadToS3 uploads data to S3 using a pre-signed URL
func (c *Client) UploadToS3(presignedURL string, data []byte) error {
	req, err := http.NewRequest("PUT", presignedURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/zip")
	req.ContentLength = int64(len(data))
	
	client := &http.Client{
		Timeout: 10 * time.Minute, // Allow up to 10 minutes for large uploads
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("S3 upload failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("S3 upload failed (status %d): %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// UpdateS3Location updates the project with the S3 location after upload
func (c *Client) UpdateS3Location(projectID, s3Location string) (*Project, error) {
	req := UpdateS3LocationRequest{
		S3Location: s3Location,
	}
	
	resp, err := c.makeRequest("POST", fmt.Sprintf("/api/projects/%s/s3-location", projectID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("update S3 location failed (status %d): %s", resp.StatusCode, string(body))
	}
	
	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}
	
	return &project, nil
}

// CreateProjectWithUpload creates a project and uploads files in one operation
func (c *Client) CreateProjectWithUpload(name, description, projectPath string) (*Project, error) {
	// Step 1: Create project record
	createReq := CreateProjectRequest{
		Name:        name,
		Description: description,
	}
	
	project, err := c.CreateProject(createReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}
	
	// Step 2: Create zip file
	zipper := filesystem.NewProjectZipper(projectPath)
	zipResult, err := zipper.CreateZip()
	if err != nil {
		return nil, fmt.Errorf("failed to create zip: %w", err)
	}
	
	// Step 3: Validate zip size
	err = filesystem.ValidateZipSize(zipResult.Data)
	if err != nil {
		return nil, fmt.Errorf("zip validation failed: %w", err)
	}
	
	// Step 4: Get upload URL
	uploadResp, err := c.GetUploadURL(project.ID, "project.zip", int64(len(zipResult.Data)))
	if err != nil {
		return nil, fmt.Errorf("failed to get upload URL: %w", err)
	}
	
	// Step 5: Upload to S3
	err = c.UploadToS3(uploadResp.URL, zipResult.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}
	
	// Step 6: Update project with S3 location
	updatedProject, err := c.UpdateS3Location(project.ID, uploadResp.S3Location)
	if err != nil {
		return nil, fmt.Errorf("failed to update S3 location: %w", err)
	}
	
	return updatedProject, nil
}
