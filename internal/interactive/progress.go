package interactive

import (
	"fmt"
	"strings"
	"time"
)

// ProgressTracker tracks and displays progress for long-running operations
type ProgressTracker struct {
	steps       []ProgressStep
	currentStep int
	startTime   time.Time
}

// ProgressStep represents a single step in a process
type ProgressStep struct {
	Name        string
	Description string
	Status      StepStatus
	StartTime   time.Time
	EndTime     time.Time
}

// StepStatus represents the status of a progress step
type StepStatus int

const (
	StepPending StepStatus = iota
	StepInProgress
	StepCompleted
	StepFailed
)

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(stepNames []string) *ProgressTracker {
	steps := make([]ProgressStep, len(stepNames))
	for i, name := range stepNames {
		steps[i] = ProgressStep{
			Name:   name,
			Status: StepPending,
		}
	}
	
	return &ProgressTracker{
		steps:     steps,
		startTime: time.Now(),
	}
}

// Start starts the progress tracker and displays initial state
func (pt *ProgressTracker) Start() {
	fmt.Println("┌─────────────────────────────────────────────────┐")
	fmt.Println("│ Creating Project...                              │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│                                                 │")
	
	pt.displaySteps()
	fmt.Println("│                                                 │")
	fmt.Println("└─────────────────────────────────────────────────┘")
}

// StartStep starts a specific step
func (pt *ProgressTracker) StartStep(stepIndex int, description string) {
	if stepIndex < 0 || stepIndex >= len(pt.steps) {
		return
	}
	
	pt.currentStep = stepIndex
	pt.steps[stepIndex].Description = description
	pt.steps[stepIndex].Status = StepInProgress
	pt.steps[stepIndex].StartTime = time.Now()
	
	pt.updateDisplay()
}

// CompleteStep marks a step as completed
func (pt *ProgressTracker) CompleteStep(stepIndex int) {
	if stepIndex < 0 || stepIndex >= len(pt.steps) {
		return
	}
	
	pt.steps[stepIndex].Status = StepCompleted
	pt.steps[stepIndex].EndTime = time.Now()
	
	pt.updateDisplay()
}

// FailStep marks a step as failed
func (pt *ProgressTracker) FailStep(stepIndex int, err error) {
	if stepIndex < 0 || stepIndex >= len(pt.steps) {
		return
	}
	
	pt.steps[stepIndex].Status = StepFailed
	pt.steps[stepIndex].EndTime = time.Now()
	pt.steps[stepIndex].Description = err.Error()
	
	pt.updateDisplay()
}

// updateDisplay updates the progress display
func (pt *ProgressTracker) updateDisplay() {
	// Move cursor up to overwrite previous display
	fmt.Printf("\033[%dA", len(pt.steps)+4) // Move up by number of steps + borders
	
	fmt.Println("┌─────────────────────────────────────────────────┐")
	fmt.Println("│ Creating Project...                              │")
	fmt.Println("├─────────────────────────────────────────────────┤")
	fmt.Println("│                                                 │")
	
	pt.displaySteps()
	fmt.Println("│                                                 │")
	fmt.Println("└─────────────────────────────────────────────────┘")
}

// displaySteps displays all steps with their current status
func (pt *ProgressTracker) displaySteps() {
	for i, step := range pt.steps {
		statusIcon := pt.getStatusIcon(step.Status)
		stepNum := fmt.Sprintf("[%d/%d]", i+1, len(pt.steps))
		
		line := fmt.Sprintf("│ %s %s %-35s │", statusIcon, stepNum, truncateString(step.Name, 35))
		fmt.Println(line)
		
		// Show description for in-progress or failed steps
		if step.Status == StepInProgress && step.Description != "" {
			desc := fmt.Sprintf("│     %s", truncateString(step.Description, 41))
			fmt.Printf("%-49s │\n", desc)
		} else if step.Status == StepFailed && step.Description != "" {
			desc := fmt.Sprintf("│     Error: %s", truncateString(step.Description, 35))
			fmt.Printf("%-49s │\n", desc)
		}
	}
}

// getStatusIcon returns the appropriate icon for a step status
func (pt *ProgressTracker) getStatusIcon(status StepStatus) string {
	switch status {
	case StepPending:
		return "[ ]"
	case StepInProgress:
		return "[~]"
	case StepCompleted:
		return "[✓]"
	case StepFailed:
		return "[✗]"
	default:
		return "[ ]"
	}
}

// Finish completes the progress tracker
func (pt *ProgressTracker) Finish(success bool) {
	totalTime := time.Since(pt.startTime)
	
	if success {
		fmt.Println("┌─────────────────────────────────────────────────┐")
		fmt.Println("│ Project Created Successfully!                    │")
		fmt.Println("├─────────────────────────────────────────────────┤")
		fmt.Println("│                                                 │")
		fmt.Printf("│ Total time: %-36s │\n", totalTime.Round(time.Second).String())
		fmt.Println("│                                                 │")
		fmt.Println("│ Configuration saved to .leanmcp/config.json    │")
		fmt.Println("│                                                 │")
		fmt.Println("│ Next Steps:                                     │")
		fmt.Println("│ • Build: leanmcp-cli build                     │")
		fmt.Println("│ • Deploy: leanmcp-cli deploy                   │")
		fmt.Println("│ • Status: leanmcp-cli status                   │")
		fmt.Println("│                                                 │")
		fmt.Println("└─────────────────────────────────────────────────┘")
	} else {
		fmt.Println("┌─────────────────────────────────────────────────┐")
		fmt.Println("│ Project Creation Failed                          │")
		fmt.Println("├─────────────────────────────────────────────────┤")
		fmt.Println("│                                                 │")
		fmt.Printf("│ Total time: %-36s │\n", totalTime.Round(time.Second).String())
		fmt.Println("│                                                 │")
		fmt.Println("│ Please check the error messages above and       │")
		fmt.Println("│ try again.                                      │")
		fmt.Println("│                                                 │")
		fmt.Println("└─────────────────────────────────────────────────┘")
	}
}

// ShowProgress displays a simple progress bar
func ShowProgress(current, total int, description string) {
	if total == 0 {
		return
	}
	
	percentage := int((float64(current) / float64(total)) * 100)
	barLength := 30
	filledLength := int((float64(current) / float64(total)) * float64(barLength))
	
	bar := strings.Repeat("█", filledLength) + strings.Repeat("░", barLength-filledLength)
	
	fmt.Printf("\r[%s] %d%% - %s", bar, percentage, description)
	
	if current == total {
		fmt.Println() // New line when complete
	}
}
