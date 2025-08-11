package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/ddod/leanmcp-cli/internal/config"
)

var (
	cfgFile string
	verbose bool
	
	// Version information
	Version = "1.0.1"
)

var rootCmd = &cobra.Command{
	Use:     "mcli",
	Version: Version,
	Short:   "CLI for interacting with LeanMCP APIs",
	Long: `A command-line interface for managing your LeanMCP projects, 
chats, deployments, and API keys.

MCLI provides a simple way to interact with LeanMCP services from the command line.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		if err := config.Initialize(cfgFile); err != nil && verbose {
			fmt.Printf("Warning: Could not initialize config: %v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.leanmcp-cli/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"verbose output")

	// Add alias command
	rootCmd.AddCommand(&cobra.Command{
		Use:    "lcli",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Use 'leanmcp-cli' or create an alias: alias lcli='leanmcp-cli'")
		},
	})
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home + "/.leanmcp-cli")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
