package display

import (
	"fmt"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/ddod/leanmcp-cli/internal/api"
)

// ProjectsTable displays projects in a table format
func ProjectsTable(projects []api.Project) {
	if len(projects) == 0 {
		fmt.Println("No projects found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Framework", "Status", "Created", "Updated"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, project := range projects {
		status := colorizeStatus(project.Status)
		framework := project.Framework
		if framework == "" {
			framework = "-"
		}
		table.Append([]string{
			project.ID[:8] + "...", // Truncate ID for readability
			project.Name,
			framework,
			status,
			project.CreatedAt.Format("2006-01-02 15:04"),
			project.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}

	table.Render()
}

// ChatsTable displays chats in a table format
func ChatsTable(chats []api.Chat) {
	if len(chats) == 0 {
		fmt.Println("No chats found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Title", "Messages", "Model", "Created"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, chat := range chats {
		title := chat.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}
		
		model := chat.ModelUsed
		if model == "" {
			model = "N/A"
		}

		table.Append([]string{
			chat.ID[:8] + "...",
			title,
			strconv.Itoa(chat.MessageCount),
			model,
			chat.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	table.Render()
}

// BuildsTable displays builds in a table format
func BuildsTable(builds []api.Build) {
	if len(builds) == 0 {
		fmt.Println("No builds found.")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Status", "Created", "Updated"})
	table.SetBorder(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, build := range builds {
		status := colorizeStatus(build.Status)
		table.Append([]string{
			build.ID[:8] + "...",
			status,
			build.CreatedAt.Format("2006-01-02 15:04"),
			build.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}

	table.Render()
}

// colorizeStatus adds color to status strings
func colorizeStatus(status string) string {
	switch status {
	case "active", "running", "success", "completed":
		return color.GreenString(status)
	case "pending", "building", "deploying":
		return color.YellowString(status)
	case "failed", "error", "inactive":
		return color.RedString(status)
	default:
		return status
	}
}

// PrintProject displays detailed project information
func PrintProject(project *api.Project) {
	fmt.Printf("%s %s\n", color.CyanString("Project:"), color.WhiteString(project.Name))
	fmt.Printf("%s %s\n", color.CyanString("ID:"), project.ID)
	fmt.Printf("%s %s\n", color.CyanString("Status:"), colorizeStatus(project.Status))
	if project.Description != "" {
		fmt.Printf("%s %s\n", color.CyanString("Description:"), project.Description)
	}
	if project.Framework != "" {
		fmt.Printf("%s %s\n", color.CyanString("Framework:"), project.Framework)
	}
	if project.RepositoryURL != "" {
		fmt.Printf("%s %s\n", color.CyanString("Repository:"), project.RepositoryURL)
	}
	if project.S3Location != "" {
		fmt.Printf("%s %s\n", color.CyanString("S3 Location:"), project.S3Location)
	}
	fmt.Printf("%s %s\n", color.CyanString("User ID:"), project.UserID)
	fmt.Printf("%s %s\n", color.CyanString("Created:"), project.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s %s\n", color.CyanString("Updated:"), project.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// PrintChat displays detailed chat information
func PrintChat(chat *api.Chat) {
	fmt.Printf("%s %s\n", color.CyanString("Chat:"), color.WhiteString(chat.Title))
	fmt.Printf("%s %s\n", color.CyanString("ID:"), chat.ID)
	fmt.Printf("%s %d\n", color.CyanString("Messages:"), chat.MessageCount)
	if chat.ModelUsed != "" {
		fmt.Printf("%s %s\n", color.CyanString("Model:"), chat.ModelUsed)
	}
	if chat.Summary != "" {
		fmt.Printf("%s %s\n", color.CyanString("Summary:"), chat.Summary)
	}
	fmt.Printf("%s %s\n", color.CyanString("Created:"), chat.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s %s\n", color.CyanString("Updated:"), chat.UpdatedAt.Format("2006-01-02 15:04:05"))
}

// PrintChatHistory displays chat messages
func PrintChatHistory(messages []api.ChatMessage) {
	if len(messages) == 0 {
		fmt.Println("No messages found.")
		return
	}

	for i, msg := range messages {
		if i > 0 {
			fmt.Println()
		}

		roleColor := color.BlueString
		if msg.Role == "assistant" {
			roleColor = color.GreenString
		}

		fmt.Printf("%s [%s] %s:\n", 
			roleColor(fmt.Sprintf("#%d", msg.MessageIndex)),
			msg.CreatedAt.Format("15:04:05"),
			roleColor(msg.Role))
		
		// Print content with indentation
		content := msg.Content
		if len(content) > 500 {
			content = content[:497] + "..."
		}
		fmt.Printf("  %s\n", content)
	}
}
