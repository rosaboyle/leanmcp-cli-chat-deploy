package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ddod/leanmcp-cli/internal/api"
)

// ProjectConfig represents the local project configuration
type ProjectConfig struct {
	Project ProjectInfo `json:"project"`
	CLI     CLIInfo     `json:"cli"`
}

// ProjectInfo contains project details from the API
type ProjectInfo struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Framework     string `json:"framework"`
	Status        string `json:"status"`
	S3Location    string `json:"s3Location"`
	RepositoryURL string `json:"repositoryUrl"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// CLIInfo contains local CLI metadata
type CLIInfo struct {
	Version     string `json:"version"`
	LastSync    string `json:"lastSync"`
	ProjectPath string `json:"projectPath"`
}

// SaveProjectConfig saves project configuration to .leanmcp/config.json
func SaveProjectConfig(projectPath string, project *api.Project) error {
	leanmcpDir := filepath.Join(projectPath, ".leanmcp")
	
	// Create .leanmcp directory
	err := os.MkdirAll(leanmcpDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create .leanmcp directory: %w", err)
	}
	
	// Create config
	config := ProjectConfig{
		Project: ProjectInfo{
			ID:            project.ID,
			Name:          project.Name,
			Description:   project.Description,
			Framework:     project.Framework,
			Status:        project.Status,
			S3Location:    project.S3Location,
			RepositoryURL: project.RepositoryURL,
			CreatedAt:     project.CreatedAt.Format(time.RFC3339),
			UpdatedAt:     project.UpdatedAt.Format(time.RFC3339),
		},
		CLI: CLIInfo{
			Version:     "1.0.0",
			LastSync:    time.Now().Format(time.RFC3339),
			ProjectPath: projectPath,
		},
	}
	
	// Save to file
	configPath := filepath.Join(leanmcpDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}

// LoadProjectConfig loads project configuration from .leanmcp/config.json
func LoadProjectConfig(projectPath string) (*ProjectConfig, error) {
	configPath := filepath.Join(projectPath, ".leanmcp", "config.json")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("project not initialized in this directory. Run 'leanmcp-cli project create' first")
		}
		return nil, fmt.Errorf("failed to read project config: %w", err)
	}
	
	var config ProjectConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("invalid project configuration: %w", err)
	}
	
	return &config, nil
}

// HasProjectConfig checks if .leanmcp/config.json exists in the current directory
func HasProjectConfig(projectPath string) bool {
	configPath := filepath.Join(projectPath, ".leanmcp", "config.json")
	_, err := os.Stat(configPath)
	return err == nil
}

// UpdateProjectConfig updates specific fields in the project config
func UpdateProjectConfig(projectPath string, updates ProjectInfo) error {
	config, err := LoadProjectConfig(projectPath)
	if err != nil {
		return err
	}
	
	// Update fields if provided
	if updates.ID != "" {
		config.Project.ID = updates.ID
	}
	if updates.Name != "" {
		config.Project.Name = updates.Name
	}
	if updates.Description != "" {
		config.Project.Description = updates.Description
	}
	if updates.Framework != "" {
		config.Project.Framework = updates.Framework
	}
	if updates.Status != "" {
		config.Project.Status = updates.Status
	}
	if updates.S3Location != "" {
		config.Project.S3Location = updates.S3Location
	}
	if updates.RepositoryURL != "" {
		config.Project.RepositoryURL = updates.RepositoryURL
	}
	
	// Update last sync time
	config.CLI.LastSync = time.Now().Format(time.RFC3339)
	
	// Save updated config
	configPath := filepath.Join(projectPath, ".leanmcp", "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	return os.WriteFile(configPath, data, 0644)
}

// GetCurrentProjectID gets the project ID from the current directory's config
func GetCurrentProjectID() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	
	config, err := LoadProjectConfig(pwd)
	if err != nil {
		return "", err
	}
	
	return config.Project.ID, nil
}
