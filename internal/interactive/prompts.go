package interactive

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ddod/leanmcp-cli/internal/filesystem"
)

// ProjectCreationFlow handles the interactive project creation process
type ProjectCreationFlow struct {
	Name        string
	Description string
	Path        string
	Files       []filesystem.FileInfo
	Stats       filesystem.FileStats
}

// CollectProjectInfo collects project information interactively
func (p *ProjectCreationFlow) CollectProjectInfo() error {
	fmt.Println("┌─────────────────────────────────────────────────┐")
	fmt.Println("│ Create New Project                               │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│                                                 │")
	
	// Get project name
	if p.Name == "" {
		name, err := p.promptProjectName()
		if err != nil {
			return err
		}
		p.Name = name
	}
	
	// Get project description (optional)
	if p.Description == "" {
		p.Description = p.promptProjectDescription()
	}
	
	// Get project path
	if p.Path == "" {
		path, err := p.promptProjectPath()
		if err != nil {
			return err
		}
		p.Path = path
	}
	
	// Scan files
	err := p.scanProjectFiles()
	if err != nil {
		return err
	}
	
	// Show confirmation
	return p.confirmCreation()
}

// promptProjectName prompts for project name
func (p *ProjectCreationFlow) promptProjectName() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("│ Project Name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		
		name := strings.TrimSpace(input)
		if name == "" {
			fmt.Println("│ Project name cannot be empty. Please try again. │")
			continue
		}
		
		if len(name) > 100 {
			fmt.Println("│ Project name too long (max 100 characters).     │")
			continue
		}
		
		return name, nil
	}
}

// promptProjectDescription prompts for project description
func (p *ProjectCreationFlow) promptProjectDescription() string {
	fmt.Println("│                                                 │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│ Project Description (Optional)                   │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│                                                 │")
	fmt.Print("│ Description: ")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	description := strings.TrimSpace(input)
	
	if description == "" {
		fmt.Println("│ (No description provided)                       │")
	}
	
	return description
}

// promptProjectPath prompts for project path
func (p *ProjectCreationFlow) promptProjectPath() (string, error) {
	fmt.Println("│                                                 │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│ Select Project Directory                         │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│                                                 │")
	
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	
	fmt.Printf("│ Current Directory: %-29s │\n", truncatePath(currentDir, 29))
	fmt.Println("│                                                 │")
	fmt.Println("│ Options:                                        │")
	fmt.Println("│   1. Use current directory                      │")
	fmt.Println("│   2. Enter custom path                          │")
	fmt.Println("│                                                 │")
	fmt.Print("│ Choice [1]: ")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(input)
	
	if choice == "" || choice == "1" {
		return currentDir, nil
	}
	
	if choice == "2" {
		return p.promptCustomPath()
	}
	
	fmt.Println("│ Invalid choice. Using current directory.        │")
	return currentDir, nil
}

// promptCustomPath prompts for a custom path
func (p *ProjectCreationFlow) promptCustomPath() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Print("│ Enter path: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		
		path := strings.TrimSpace(input)
		if path == "" {
			fmt.Println("│ Path cannot be empty. Please try again.         │")
			continue
		}
		
		// Expand ~ to home directory
		if strings.HasPrefix(path, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				fmt.Printf("│ Cannot expand ~: %s                             │\n", err.Error())
				continue
			}
			path = filepath.Join(home, path[2:])
		}
		
		// Convert to absolute path
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("│ Invalid path: %s                                │\n", err.Error())
			continue
		}
		
		// Validate directory
		err = filesystem.ValidateDirectory(absPath)
		if err != nil {
			fmt.Printf("│ %s                                              │\n", truncateString(err.Error(), 45))
			continue
		}
		
		return absPath, nil
	}
}

// scanProjectFiles scans the selected directory
func (p *ProjectCreationFlow) scanProjectFiles() error {
	fmt.Println("│                                                 │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│ Scanning Directory...                           │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	
	zipper := filesystem.NewProjectZipper(p.Path)
	files, stats, err := zipper.PreviewFiles(10) // Preview first 10 files
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}
	
	p.Files = files
	p.Stats = stats
	
	fmt.Println("│                                                 │")
	fmt.Printf("│ Found: %d files, %s                           │\n", 
		stats.TotalFiles, 
		truncateString(filesystem.GetHumanReadableSize(stats.TotalSize), 20))
	fmt.Println("│                                                 │")
	
	// Show preview of files
	if len(files) > 0 {
		fmt.Println("│ File Preview:                                   │")
		for i, file := range files {
			if i >= 5 { // Show max 5 files in preview
				break
			}
			fmt.Printf("│   %s                                            │\n", 
				truncateString(file.RelPath, 43))
		}
		
		if len(files) > 5 {
			fmt.Printf("│   ... and %d more files                        │\n", len(files)-5)
		}
		fmt.Println("│                                                 │")
	}
	
	return nil
}

// confirmCreation shows confirmation dialog
func (p *ProjectCreationFlow) confirmCreation() error {
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│ Confirm Project Creation                         │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│                                                 │")
	fmt.Printf("│ Name:        %-35s │\n", truncateString(p.Name, 35))
	
	if p.Description != "" {
		fmt.Printf("│ Description: %-35s │\n", truncateString(p.Description, 35))
	} else {
		fmt.Println("│ Description: (none)                             │")
	}
	
	fmt.Printf("│ Path:        %-35s │\n", truncateString(p.Path, 35))
	fmt.Printf("│ Files:       %d files (%s)                      │\n", 
		p.Stats.TotalFiles, 
		filesystem.GetHumanReadableSize(p.Stats.TotalSize))
	fmt.Println("│                                                 │")
	fmt.Println("│ Continue? [Y/n]: ")
	fmt.Print("│ > ")
	
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	confirmation := strings.ToLower(strings.TrimSpace(input))
	
	if confirmation == "" || confirmation == "y" || confirmation == "yes" {
		fmt.Println("└─────────────────────────────────────────────────┘")
		return nil
	}
	
	fmt.Println("│                                                 │")
	fmt.Println("│ Project creation cancelled.                     │")
	fmt.Println("└─────────────────────────────────────────────────┘")
	return fmt.Errorf("project creation cancelled by user")
}

// truncateString truncates a string to fit within a given length
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	
	if maxLen <= 3 {
		return s[:maxLen]
	}
	
	return s[:maxLen-3] + "..."
}

// truncatePath truncates a path to fit within a given length, keeping the end
func truncatePath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}
	
	if maxLen <= 3 {
		return path[:maxLen]
	}
	
	return "..." + path[len(path)-(maxLen-3):]
}
