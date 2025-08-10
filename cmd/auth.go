package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/ddod/leanmcp-cli/internal/auth"
	"github.com/ddod/leanmcp-cli/internal/api"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Commands for managing authentication with the LeanMCP API",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with API key",
	Long: `Authenticate with your LeanMCP API key.

Your API key should start with 'airtrain_' and will be stored securely
in your local configuration file (~/.leanmcp-cli/config.yaml).

Example:
  leanmcp-cli auth login --api-key airtrain_your_key_here`,
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			return fmt.Errorf("--api-key is required")
		}

		// Validate API key format
		if err := auth.ValidateAPIKeyFormat(apiKey); err != nil {
			return fmt.Errorf("invalid API key format: %v", err)
		}

		fmt.Println("🔐 Authenticating with API key...")

		// For now, just store the credentials without validation
		// (as requested by user - no actual API validation)
		if err := auth.StoreCredentials(apiKey, nil); err != nil {
			return fmt.Errorf("failed to store credentials: %v", err)
		}

		fmt.Printf("✅ %s\n", color.GreenString("Successfully stored API key!"))
		fmt.Printf("Your API key has been securely stored in ~/.leanmcp-cli/config.yaml\n")
		fmt.Printf("You can now use other commands to interact with the API.\n")
		
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove stored credentials",
	Long:  "Remove stored API key and authentication information from local configuration.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := auth.ClearCredentials(); err != nil {
			return fmt.Errorf("failed to clear credentials: %v", err)
		}

		fmt.Printf("✅ %s\n", color.GreenString("Successfully logged out!"))
		fmt.Printf("Your stored credentials have been removed.\n")
		
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current authentication status",
	Long:  "Display information about the currently stored API key and authentication status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.LoadCredentials()
		if err != nil {
			fmt.Printf("❌ %s\n", color.RedString("Not authenticated"))
			fmt.Printf("Run 'leanmcp-cli auth login --api-key <your-key>' to authenticate.\n")
			return nil
		}

		fmt.Printf("✅ %s\n", color.GreenString("Authenticated"))
		fmt.Printf("%s %s\n", color.CyanString("API Key:"), maskAPIKey(creds.APIKey))
		
		if creds.UserEmail != "" {
			fmt.Printf("%s %s\n", color.CyanString("Email:"), creds.UserEmail)
		}
		
		if len(creds.Scopes) > 0 {
			fmt.Printf("%s %v\n", color.CyanString("Scopes:"), creds.Scopes)
		}
		
		fmt.Printf("%s %s\n", color.CyanString("Stored:"), creds.StoredAt.Format("2006-01-02 15:04:05"))
		
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check API connection status",
	Long:  "Test the connection to the API server and validate the stored API key.",
	RunE: func(cmd *cobra.Command, args []string) error {
		creds, err := auth.LoadCredentials()
		if err != nil {
			fmt.Printf("❌ %s\n", color.RedString("Not authenticated"))
			fmt.Printf("Run 'leanmcp-cli auth login --api-key <your-key>' to authenticate.\n")
			return nil
		}

		fmt.Println("🔍 Testing API connection...")

		client := api.NewClient(creds.APIKey)
		err = client.TestConnection()
		if err != nil {
			fmt.Printf("❌ %s: %v\n", color.RedString("Connection failed"), err)
			return nil
		}

		fmt.Printf("✅ %s\n", color.GreenString("API connection successful!"))
		
		// Update last used timestamp
		_ = auth.UpdateLastUsed()
		
		return nil
	},
}

// maskAPIKey masks the API key for display purposes
func maskAPIKey(apiKey string) string {
	if len(apiKey) <= 12 {
		return "***"
	}
	
	prefix := apiKey[:8]  // Show first 8 characters
	suffix := apiKey[len(apiKey)-4:] // Show last 4 characters
	
	return fmt.Sprintf("%s...%s", prefix, suffix)
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(whoamiCmd)
	authCmd.AddCommand(statusCmd)

	// Login command flags
	loginCmd.Flags().String("api-key", "", "API key for authentication (required)")
	loginCmd.MarkFlagRequired("api-key")
}
