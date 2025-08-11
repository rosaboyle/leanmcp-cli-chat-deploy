package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("leanmcp-cli version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
