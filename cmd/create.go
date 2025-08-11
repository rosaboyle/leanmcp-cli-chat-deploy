package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the top-level create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long: `Create a new project and upload files from a directory.

This is a convenience alias for 'leanmcp-cli projects create'.

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
	RunE: projectsCreateCmd.RunE, // Reuse the exact same function
}

func init() {
	// Add to root command
	rootCmd.AddCommand(createCmd)
	
	// Copy all flags from the projects create command
	createCmd.Flags().StringP("name", "n", "", "Project name")
	createCmd.Flags().StringP("description", "d", "", "Project description")
	createCmd.Flags().StringP("path", "p", "", "Path to project directory (defaults to current directory)")
}
