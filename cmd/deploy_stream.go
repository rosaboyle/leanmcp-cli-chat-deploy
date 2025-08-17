package cmd

import (
	"fmt"
	"strings"

	"github.com/ddod/leanmcp-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	projectID     string
	containerPort int
	secretIDs     []string
)

// deployStreamCmd represents the deploy-stream command
var deployStreamCmd = &cobra.Command{
	Use:   "deploy-stream",
	Short: "Deploy a project end-to-end with real-time streaming updates",
	Long: `Deploy a project through the complete pipeline with real-time progress updates.
This command triggers the full deployment process: build → containerize → deploy → return URL.

Examples:
  # Basic deployment
  leanmcp deploy-stream --project-id proj_1234567890abcdef

  # Advanced deployment with custom port and secrets
  leanmcp deploy-stream --project-id proj_1234567890abcdef --port 3000 --secrets secret1,secret2`,
	RunE: runDeployStream,
}

func runDeployStream(cmd *cobra.Command, args []string) error {
	// Validate required parameters
	if projectID == "" {
		return fmt.Errorf("project-id is required")
	}

	// Get authenticated client
	client, err := getAuthenticatedClient()
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Prepare deployment request
	request := &api.DeployStreamRequest{
		ProjectID:     projectID,
		ContainerPort: containerPort,
		SecretIDs:     secretIDs,
	}

	fmt.Printf("Starting end-to-end deployment for project: %s\n", projectID)
	if containerPort > 0 {
		fmt.Printf("Container port: %d\n", containerPort)
	}
	if len(secretIDs) > 0 {
		fmt.Printf("Secrets: %v\n", secretIDs)
	}
	fmt.Println("Connecting to deployment stream...")

	// Check if verbose mode is enabled
	verbose, _ := cmd.Flags().GetBool("verbose")
	if verbose {
		fmt.Printf("Debug: Request payload: %+v\n", request)
	}
	fmt.Println()

	// Start streaming deployment
	err = client.DeployAndStream(request, func(update *api.StreamUpdate) error {
		return handleStreamUpdate(update, verbose)
	})
	if err != nil {
		return fmt.Errorf("deployment failed: %w", err)
	}

	return nil
}

// handleStreamUpdate processes each streaming update from the server
func handleStreamUpdate(update *api.StreamUpdate, verbose bool) error {
	if verbose {
		fmt.Printf("Debug: Raw update: %+v\n", update)
	}
	// Check for failure conditions first
	if update.BuildStatus == "failed" || update.CurrentStep == "FAILED" {
		errorMsg := getErrorMessage(update)
		fmt.Printf("\n Build failed: %s\n", errorMsg)

		// Note: Build logs API endpoint not available (returns 404)
		if update.BuildID != "" {
			fmt.Printf("Build ID: %s\n", update.BuildID)
			fmt.Printf("Contact support with this Build ID for detailed logs\n")
		}

		return fmt.Errorf("build failed: %s", errorMsg)
	}

	// Handle based on currentStep from actual API response
	switch update.CurrentStep {
	case "BUILDING":
		// Show building progress
		progressBar := createProgressBar(int(update.Progress))
		fmt.Printf("\r%s %s (%.1f%%) - %s",
			getStepIcon(update.CurrentStep),
			update.CurrentStep,
			update.Progress,
			progressBar)

		if update.EstimatedTimeRemaining > 0 {
			fmt.Printf(" - ETA: %ds", update.EstimatedTimeRemaining)
		}

	case "DEPLOYING":
		progressBar := createProgressBar(int(update.Progress))
		fmt.Printf("\r%s %s (%.1f%%) - %s",
			getStepIcon(update.CurrentStep),
			update.CurrentStep,
			update.Progress,
			progressBar)

		if update.EstimatedTimeRemaining > 0 {
			fmt.Printf(" - ETA: %ds", update.EstimatedTimeRemaining)
		}

	case "COMPLETED", "COMPLETE":
		fmt.Printf("\n Deployment completed successfully!\n")
		if update.DeploymentURL != "" {
			fmt.Printf("Your application is live at: %s\n", update.DeploymentURL)
		}
		if update.DeploymentID != "" {
			fmt.Printf("Deployment ID: %s\n", update.DeploymentID)
		}

	default:
		// Handle any other steps or show progress
		if update.Progress > 0 {
			progressBar := createProgressBar(int(update.Progress))
			fmt.Printf("\r%s %s (%.1f%%) - %s",
				getStepIcon(update.CurrentStep),
				update.CurrentStep,
				update.Progress,
				progressBar)

			if update.EstimatedTimeRemaining > 0 {
				fmt.Printf(" - ETA: %ds", update.EstimatedTimeRemaining)
			}
		} else if update.Message != "" {
			fmt.Printf("%s\n", update.Message)
		}
	}

	return nil
}

// getErrorMessage extracts a meaningful error message from the update
func getErrorMessage(update *api.StreamUpdate) string {
	if update.Error != "" {
		return update.Error
	}
	if update.Message != "" {
		return update.Message
	}
	if update.BuildStatus == "failed" {
		return fmt.Sprintf("Build failed (BuildID: %s). Check server logs for details.", update.BuildID)
	}
	return "Unknown error occurred"
}

// createProgressBar creates a visual progress bar
func createProgressBar(progress int) string {
	const barLength = 20
	filled := (progress * barLength) / 100

	bar := "["
	for i := 0; i < barLength; i++ {
		if i < filled {
			bar += "="
		} else if i == filled && progress < 100 {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"

	return bar
}

// getStepIcon returns an appropriate icon for each deployment step
func getStepIcon(step string) string {
	switch strings.ToUpper(step) {
	case "BUILDING", "BUILD":
		return "BUILDING"
	case "CONTAINERIZING", "CONTAINER":
		return "CONTAINERIZING"
	case "DEPLOYING", "DEPLOY":
		return "DEPLOYING"
	case "CONFIGURING", "CONFIG":
		return "CONFIGURING"
	case "TESTING", "TEST":
		return "TESTING"
	case "COMPLETED", "COMPLETE":
		return "COMPLETED"
	default:
		return "unknown"
	}
}

func init() {
	rootCmd.AddCommand(deployStreamCmd)

	// Required flags
	deployStreamCmd.Flags().StringVarP(&projectID, "project-id", "p", "", "Project ID to deploy (required)")
	deployStreamCmd.MarkFlagRequired("project-id")

	// Optional flags
	deployStreamCmd.Flags().IntVar(&containerPort, "port", 0, "Container port (defaults to 3001)")
	deployStreamCmd.Flags().StringSliceVar(&secretIDs, "secrets", []string{}, "Comma-separated list of secret IDs to inject")
}
