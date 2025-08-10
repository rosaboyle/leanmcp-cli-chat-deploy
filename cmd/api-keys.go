package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var apiKeysCmd = &cobra.Command{
	Use:   "api-keys",
	Short: "Manage API keys",
	Long:  "Commands for managing your LeanMCP API keys",
}

var apiKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List API keys",
	Long:  "List all API keys associated with your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		fmt.Println("🔑 Fetching API key information...")

		// Try to get info about the current API key
		keyInfo, err := client.GetAPIKeyInfo()
		if err != nil {
			return fmt.Errorf("failed to get API key info: %v", err)
		}

		fmt.Printf("✅ %s\n\n", color.GreenString("Current API Key Information:"))
		fmt.Printf("%s %s\n", color.CyanString("ID:"), keyInfo.ID)
		fmt.Printf("%s %s\n", color.CyanString("Name:"), keyInfo.Name)
		fmt.Printf("%s %v\n", color.CyanString("Scopes:"), keyInfo.Scopes)
		fmt.Printf("%s %t\n", color.CyanString("Active:"), keyInfo.IsActive)
		fmt.Printf("%s %s\n", color.CyanString("Created:"), keyInfo.CreatedAt.Format("2006-01-02 15:04:05"))
		
		if keyInfo.ExpiresAt != nil {
			fmt.Printf("%s %s\n", color.CyanString("Expires:"), keyInfo.ExpiresAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%s %s\n", color.CyanString("Expires:"), "Never")
		}

		return nil
	},
}

var apiKeysInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show current API key info",
	Long:  "Display detailed information about the currently authenticated API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		fmt.Println("🔍 Getting API key information...")

		keyInfo, err := client.GetAPIKeyInfo()
		if err != nil {
			return fmt.Errorf("failed to get API key info: %v", err)
		}

		fmt.Printf("✅ %s\n\n", color.GreenString("API Key Details:"))
		fmt.Printf("%s %s\n", color.CyanString("ID:"), keyInfo.ID)
		fmt.Printf("%s %s\n", color.CyanString("Name:"), keyInfo.Name)
		
		// Display scopes with colors
		fmt.Printf("%s ", color.CyanString("Scopes:"))
		for i, scope := range keyInfo.Scopes {
			if i > 0 {
				fmt.Print(", ")
			}
			scopeColor := color.GreenString
			switch scope {
			case "ADMIN":
				scopeColor = color.RedString
			case "BUILD_AND_DEPLOY":
				scopeColor = color.YellowString
			}
			fmt.Print(scopeColor(scope))
		}
		fmt.Println()
		
		status := "❌ Inactive"
		if keyInfo.IsActive {
			status = "✅ Active"
		}
		fmt.Printf("%s %s\n", color.CyanString("Status:"), status)
		
		fmt.Printf("%s %s\n", color.CyanString("Created:"), keyInfo.CreatedAt.Format("2006-01-02 15:04:05"))
		
		if keyInfo.ExpiresAt != nil {
			fmt.Printf("%s %s\n", color.CyanString("Expires:"), keyInfo.ExpiresAt.Format("2006-01-02 15:04:05"))
		} else {
			fmt.Printf("%s %s\n", color.CyanString("Expires:"), color.GreenString("Never"))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(apiKeysCmd)
	apiKeysCmd.AddCommand(apiKeysListCmd)
	apiKeysCmd.AddCommand(apiKeysInfoCmd)
}
