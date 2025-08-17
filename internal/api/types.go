package api

import "time"

// APIKeyInfo represents information about an API key
type APIKeyInfo struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Scopes    []string  `json:"scopes"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// Project represents a project
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Status        string    `json:"status"`
	Framework     string    `json:"framework,omitempty"`
	RepositoryURL string    `json:"repositoryUrl,omitempty"`
	S3Location    string    `json:"s3Location,omitempty"`
	FirebaseUID   string    `json:"firebaseUid,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	UserID        string    `json:"userId"`
}

// Chat represents a chat conversation
type Chat struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary,omitempty"`
	ModelUsed    string    `json:"modelUsed,omitempty"`
	MessageCount int       `json:"messageCount"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	UserID       string    `json:"userId"`
}

// ChatMessage represents a message in a chat
type ChatMessage struct {
	ID           string    `json:"id"`
	ChatID       string    `json:"chatId"`
	Role         string    `json:"role"` // "user" or "assistant"
	Content      string    `json:"content"`
	MessageIndex int       `json:"messageIndex"`
	CreatedAt    time.Time `json:"createdAt"`
	UserID       string    `json:"userId"`
}

// Deployment represents a deployment
type Deployment struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Status    string    `json:"status"`
	URL       string    `json:"url,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Build represents a project build
type Build struct {
	ID        string    `json:"id"`
	ProjectID string    `json:"projectId"`
	Status    string    `json:"status"`
	BuildLog  string    `json:"buildLog,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UploadURLRequest represents a request to get an upload URL
type UploadURLRequest struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
}

// UploadURLResponse represents the response from upload URL request
type UploadURLResponse struct {
	URL        string `json:"url"`
	S3Location string `json:"s3Location"`
}

// UpdateS3LocationRequest represents a request to update S3 location
type UpdateS3LocationRequest struct {
	S3Location string `json:"s3Location"`
}

// CreateChatRequest represents a request to create a chat
type CreateChatRequest struct {
	Title     string `json:"title"`
	ModelUsed string `json:"modelUsed,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// DeployStreamRequest represents a request for end-to-end deployment streaming
type DeployStreamRequest struct {
	ProjectID     string   `json:"projectId"`
	ContainerPort int      `json:"containerPort,omitempty"`
	SecretIDs     []string `json:"secretIds,omitempty"`
}

// StreamUpdate represents a single update from the deployment stream
type StreamUpdate struct {
	Type                    string `json:"type"`                    // "progress", "log", "error", "complete"
	DeploymentID            string `json:"deploymentId,omitempty"`
	CurrentStep             string `json:"currentStep,omitempty"`
	BuildID                 string `json:"buildId,omitempty"`
	BuildStatus             string `json:"buildStatus,omitempty"`
	Progress                float64 `json:"progress,omitempty"`
	EstimatedTimeRemaining  int     `json:"estimatedTimeRemaining,omitempty"`
	Message                 string  `json:"message,omitempty"`
	Error                   string  `json:"error,omitempty"`
	DeploymentURL           string  `json:"deploymentUrl,omitempty"`
}

// EndToEndDeploymentResponse represents the final deployment response
type EndToEndDeploymentResponse struct {
	DeploymentID  string `json:"deploymentId"`
	DeploymentURL string `json:"deploymentUrl"`
	BuildID       string `json:"buildId"`
	BuildStatus   string `json:"buildStatus"`
	Error         string `json:"error,omitempty"`
}

// BuildLogsResponse represents build logs from the server
type BuildLogsResponse struct {
	BuildID string   `json:"buildId"`
	Logs    []string `json:"logs"`
	Error   string   `json:"error,omitempty"`
}
