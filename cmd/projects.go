package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/ddod/leanmcp-cli/internal/auth"
	"github.com/ddod/leanmcp-cli/internal/api"
	"github.com/ddod/leanmcp-cli/internal/display"
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
	Long:  "Create a new project with the specified name and optional description",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")

		if name == "" {
			return fmt.Errorf("--name is required")
		}

		fmt.Printf("🚀 Creating project '%s'...\n", name)

		req := api.CreateProjectRequest{
			Name:        name,
			Description: description,
		}

		project, err := client.CreateProject(req)
		if err != nil {
			return fmt.Errorf("failed to create project: %v", err)
		}

		fmt.Printf("✅ %s\n\n", color.GreenString("Project created successfully!"))
		display.PrintProject(project)

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
	projectsCreateCmd.Flags().String("name", "", "Project name (required)")
	projectsCreateCmd.Flags().String("description", "", "Project description")
	projectsCreateCmd.MarkFlagRequired("name")

	// Delete command flags
	projectsDeleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")
}
