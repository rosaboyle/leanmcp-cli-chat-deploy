package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/ddod/leanmcp-cli/internal/api"
	"github.com/ddod/leanmcp-cli/internal/display"
)

var chatsCmd = &cobra.Command{
	Use:   "chats",
	Short: "Manage chats",
	Long:  "Commands for managing your LeanMCP chat conversations",
}

var chatsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all chats",
	Long:  "List all chat conversations associated with your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		fmt.Println("💬 Fetching chats...")

		chats, err := client.ListChats()
		if err != nil {
			return fmt.Errorf("failed to list chats: %v", err)
		}

		fmt.Printf("\nFound %d chat(s):\n\n", len(chats))
		display.ChatsTable(chats)

		return nil
	},
}

var chatsShowCmd = &cobra.Command{
	Use:   "show <chat-id>",
	Short: "Show chat details",
	Long:  "Display detailed information about a specific chat conversation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		chatID := args[0]
		fmt.Printf("🔍 Fetching chat %s...\n\n", chatID)

		chat, err := client.GetChat(chatID)
		if err != nil {
			return fmt.Errorf("failed to get chat: %v", err)
		}

		display.PrintChat(chat)

		return nil
	},
}

var chatsHistoryCmd = &cobra.Command{
	Use:   "history <chat-id>",
	Short: "Show chat message history",
	Long:  "Display the complete message history for a specific chat conversation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		chatID := args[0]
		limit, _ := cmd.Flags().GetInt("limit")

		fmt.Printf("📜 Fetching chat history for %s...\n\n", chatID)

		messages, err := client.GetChatHistory(chatID)
		if err != nil {
			return fmt.Errorf("failed to get chat history: %v", err)
		}

		// Apply limit if specified
		if limit > 0 && len(messages) > limit {
			messages = messages[len(messages)-limit:] // Show last N messages
		}

		fmt.Printf("Showing %d message(s):\n\n", len(messages))
		display.PrintChatHistory(messages)

		return nil
	},
}

var chatsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new chat",
	Long:  "Create a new chat conversation with the specified title",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		model, _ := cmd.Flags().GetString("model")

		if title == "" {
			return fmt.Errorf("--title is required")
		}

		fmt.Printf("💬 Creating chat '%s'...\n", title)

		req := api.CreateChatRequest{
			Title:     title,
			ModelUsed: model,
		}

		chat, err := client.CreateChat(req)
		if err != nil {
			return fmt.Errorf("failed to create chat: %v", err)
		}

		fmt.Printf("✅ %s\n\n", color.GreenString("Chat created successfully!"))
		display.PrintChat(chat)

		return nil
	},
}

var chatsDeleteCmd = &cobra.Command{
	Use:   "delete <chat-id>",
	Short: "Delete a chat",
	Long:  "Delete a chat conversation by its ID. This action cannot be undone.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := getAuthenticatedClient()
		if err != nil {
			return err
		}

		chatID := args[0]
		
		// Check if force flag is provided
		force, _ := cmd.Flags().GetBool("force")
		
		if !force {
			fmt.Printf("⚠️  %s\n", color.YellowString("WARNING: This will permanently delete the chat and all messages."))
			fmt.Printf("Use --force to confirm deletion.\n")
			return nil
		}

		fmt.Printf("🗑️  Deleting chat %s...\n", chatID)

		if err := client.DeleteChat(chatID); err != nil {
			return fmt.Errorf("failed to delete chat: %v", err)
		}

		fmt.Printf("✅ %s\n", color.GreenString("Chat deleted successfully!"))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(chatsCmd)
	chatsCmd.AddCommand(chatsListCmd)
	chatsCmd.AddCommand(chatsShowCmd)
	chatsCmd.AddCommand(chatsHistoryCmd)
	chatsCmd.AddCommand(chatsCreateCmd)
	chatsCmd.AddCommand(chatsDeleteCmd)

	// History command flags
	chatsHistoryCmd.Flags().Int("limit", 0, "Limit number of messages to show (0 = all)")

	// Create command flags
	chatsCreateCmd.Flags().String("title", "", "Chat title (required)")
	chatsCreateCmd.Flags().String("model", "", "Model to use for the chat")
	chatsCreateCmd.MarkFlagRequired("title")

	// Delete command flags
	chatsDeleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")
}
