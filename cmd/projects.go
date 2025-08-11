package cmd

import (
	"fmt"
	"strings"

	"github.com/ddod/leanmcp-cli/internal/api"
	"github.com/ddod/leanmcp-cli/internal/config"
	"github.com/ddod/leanmcp-cli/internal/display"
	"github.com/ddod/leanmcp-cli/internal/filesystem"
	"github.com/ddod/leanmcp-cli/internal/interactive"
	"github.com/ddod/leanmcp-cli/internal/auth"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage projects",
	Long:  "Commands for managing your LeanMCP projects",
}

var projectsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  "List all projects associated with your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			handleAuthError()
			return nil
		}

		fmt.Println("📋 Fetching projects...")

		projects, err := client.ListProjects()
		if err != nil {
			handleAPIError(err, "list projects")
			return nil // Return nil to prevent usage help from showing
		}

		fmt.Printf("\nFound %d project(s):\n\n", len(projects))
		display.ProjectsTable(projects)

		return nil
	},
}

var projectsShowCmd = &cobra.Command{
	Use:   "show <project-id>",
	Short: "Show project details",
	Long:  "Display detailed information about a specific project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		projectID := args[0]
		fmt.Printf("🔍 Fetching project %s...\n\n", projectID)

		project, err := client.GetProject(projectID)
		if err != nil {
			handleAPIError(err, "get project details")
			return nil
		}

		display.PrintProject(project)

		return nil
	},
}

var projectsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long: `Create a new project and upload files from a directory.

Supports both interactive and flag-based modes:
- Interactive mode: Run without flags to be prompted for all details
- Flag mode: Use --name, --description, and --path flags
- Partial flags: Provide some flags, be prompted for missing ones

The command will:
1. Create a project record in LeanMCP
2. Scan the specified directory (respecting .gitignore)
3. Create a zip archive of the project files
4. Upload the zip to S3
5. Update the project with the S3 location
6. Save local configuration in .leanmcp/config.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		// Get command line flags
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		projectPath, _ := cmd.Flags().GetString("path")

		// Create progress tracker
		steps := []string{
			"Collect project information",
			"Create project record",
			"Scan and zip files",
			"Upload to S3",
			"Update project record",
			"Save local configuration",
		}
		
		progress := interactive.NewProgressTracker(steps)
		progress.Start()

		// Step 1: Collect project information
		progress.StartStep(0, "Gathering project details...")
		
		flow := &interactive.ProjectCreationFlow{
			Name:        name,
			Description: description,
			Path:        projectPath,
		}

		err = flow.CollectProjectInfo()
		if err != nil {
			progress.FailStep(0, err)
			progress.Finish(false)
			return err
		}
		
		progress.CompleteStep(0)

		// Step 2: Create project record
		progress.StartStep(1, fmt.Sprintf("Creating project '%s'...", flow.Name))
		
		createReq := api.CreateProjectRequest{
			Name:        flow.Name,
			Description: flow.Description,
		}

		project, err := client.CreateProject(createReq)
		if err != nil {
			progress.FailStep(1, err)
			progress.Finish(false)
			return fmt.Errorf("failed to create project: %w", err)
		}
		
		progress.CompleteStep(1)

		// Step 3: Scan and zip files
		progress.StartStep(2, fmt.Sprintf("Zipping %d files...", flow.Stats.TotalFiles))
		
		zipper := filesystem.NewProjectZipper(flow.Path)
		zipResult, err := zipper.CreateZip()
		if err != nil {
			progress.FailStep(2, err)
			progress.Finish(false)
			return fmt.Errorf("failed to create zip: %w", err)
		}

		// Validate zip size
		err = filesystem.ValidateZipSize(zipResult.Data)
		if err != nil {
			progress.FailStep(2, err)
			progress.Finish(false)
			return fmt.Errorf("zip validation failed: %w", err)
		}
		
		progress.CompleteStep(2)

		// Step 4: Upload to S3
		progress.StartStep(3, "Getting upload URL...")
		
		uploadResp, err := client.GetUploadURL(project.ID, "project.zip", int64(len(zipResult.Data)))
		if err != nil {
			progress.FailStep(3, err)
			progress.Finish(false)
			return fmt.Errorf("failed to get upload URL: %w", err)
		}

		progress.StartStep(3, "Uploading to S3...")
		err = client.UploadToS3(uploadResp.URL, zipResult.Data)
		if err != nil {
			progress.FailStep(3, err)
			progress.Finish(false)
			return fmt.Errorf("failed to upload to S3: %w", err)
		}
		
		progress.CompleteStep(3)

		// Step 5: Update project record
		progress.StartStep(4, "Updating project record...")
		
		updatedProject, err := client.UpdateS3Location(project.ID, uploadResp.S3Location)
		if err != nil {
			progress.FailStep(4, err)
			progress.Finish(false)
			return fmt.Errorf("failed to update S3 location: %w", err)
		}
		
		progress.CompleteStep(4)

		// Step 6: Save local configuration
		progress.StartStep(5, "Saving local configuration...")
		
		err = config.SaveProjectConfig(flow.Path, updatedProject)
		if err != nil {
			progress.FailStep(5, err)
			progress.Finish(false)
			return fmt.Errorf("failed to save local config: %w", err)
		}
		
		progress.CompleteStep(5)

		// Finish successfully
		progress.Finish(true)

		// Show project summary
		fmt.Println("\n" + color.GreenString("Project Details:"))
		display.PrintProject(updatedProject)

		return nil
	},
}

var projectsDeleteCmd = &cobra.Command{
	Use:   "delete <project-id>",
	Short: "Delete a project",
	Long:  "Delete a project by its ID. This action cannot be undone.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		projectID := args[0]
		
		// Check if force flag is provided
		force, _ := cmd.Flags().GetBool("force")
		
		if !force {
			fmt.Printf("⚠️  %s\n", color.YellowString("WARNING: This will permanently delete the project and all associated data."))
			fmt.Printf("Use --force to confirm deletion.\n")
			return nil
		}

		fmt.Printf("🗑️  Deleting project %s...\n", projectID)

		if err := client.DeleteProject(projectID); err != nil {
			return fmt.Errorf("failed to delete project: %v", err)
		}

		fmt.Printf("✅ %s\n", color.GreenString("Project deleted successfully!"))

		return nil
	},
}

var projectsBuildsCmd = &cobra.Command{
	Use:   "builds <project-id>",
	Short: "List project builds",
	Long:  "List all builds for a specific project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		projectID := args[0]
		fmt.Printf("🔨 Fetching builds for project %s...\n", projectID)

		builds, err := client.GetProjectBuilds(projectID)
		if err != nil {
			return fmt.Errorf("failed to get project builds: %v", err)
		}

		fmt.Printf("\nFound %d build(s):\n\n", len(builds))
		display.BuildsTable(builds)

		return nil
	},
}

var projectsBuildCmd = &cobra.Command{
	Use:   "build <project-id>",
	Short: "Start a new build",
	Long:  "Start a new build for the specified project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		projectID := args[0]
		fmt.Printf("🔨 Starting build for project %s...\n", projectID)

		build, err := client.StartBuild(projectID)
		if err != nil {
			return fmt.Errorf("failed to start build: %v", err)
		}

		fmt.Printf("✅ %s\n", color.GreenString("Build started successfully!"))
		fmt.Printf("%s %s\n", color.CyanString("Build ID:"), build.ID)
		fmt.Printf("%s %s\n", color.CyanString("Status:"), build.Status)
		fmt.Printf("%s %s\n", color.CyanString("Created:"), build.CreatedAt.Format("2006-01-02 15:04:05"))

		return nil
	},
}

// getAuthenticatedClient creates an authenticated API client
func getAuthenticatedClient() (*api.Client, error) {
	creds, err := auth.LoadCredentials()
	if err != nil {
		return nil, fmt.Errorf("not authenticated. Run 'leanmcp-cli auth login --api-key <your-key>' first")
	}

	return api.NewClient(creds.APIKey), nil
}

// handleAuthError handles authentication errors with user-friendly messages
func handleAuthError() {
	fmt.Printf("❌ %s\n", color.RedString("Not authenticated"))
	fmt.Printf("Please run: %s\n", color.CyanString("leanmcp-cli auth login --api-key <your-key>"))
}

// handleAPIError provides user-friendly error messages for common API errors
func handleAPIError(err error, action string) {
	if strings.Contains(err.Error(), "status 401") {
		fmt.Printf("❌ %s\n", color.RedString("Authentication failed"))
		fmt.Printf("Your API key is invalid or has expired.\n")
		fmt.Printf("Please run: %s\n", color.CyanString("leanmcp-cli auth login --api-key <your-key>"))
		return
	}
	if strings.Contains(err.Error(), "status 403") {
		fmt.Printf("❌ %s\n", color.RedString("Access denied"))
		fmt.Printf("Your API key doesn't have permission to %s.\n", action)
		return
	}
	if strings.Contains(err.Error(), "status 404") {
		fmt.Printf("❌ %s\n", color.RedString("Not found"))
		fmt.Printf("The requested resource was not found.\n")
		return
	}
	if strings.Contains(err.Error(), "connection failed") {
		fmt.Printf("❌ %s\n", color.RedString("Connection failed"))
		fmt.Printf("Unable to connect to the API. Please check your internet connection.\n")
		return
	}
	// Generic error for other cases
	fmt.Printf("❌ %s: %v\n", color.RedString("Error"), err)
}

func init() {
	rootCmd.AddCommand(projectsCmd)
	projectsCmd.AddCommand(projectsListCmd)
	projectsCmd.AddCommand(projectsShowCmd)
	projectsCmd.AddCommand(projectsCreateCmd)
	projectsCmd.AddCommand(projectsDeleteCmd)
	projectsCmd.AddCommand(projectsBuildsCmd)
	projectsCmd.AddCommand(projectsBuildCmd)

	// Create command flags
	projectsCreateCmd.Flags().StringP("name", "n", "", "Project name")
	projectsCreateCmd.Flags().StringP("description", "d", "", "Project description")
	projectsCreateCmd.Flags().StringP("path", "p", "", "Path to project directory (defaults to current directory)")
	// Note: name is no longer required - interactive mode will prompt if missing

	// Delete command flags
	projectsDeleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")
}
