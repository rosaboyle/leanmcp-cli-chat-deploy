package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Initialize sets up the configuration
func Initialize(cfgFile string) error {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		configDir := filepath.Join(home, ".leanmcp-cli")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return err
		}

		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv()
	
	// Set defaults
	viper.SetDefault("base_url", "https://api.leanmcp.ai")

	// Try to read config, but don't fail if it doesn't exist
	_ = viper.ReadInConfig()
	
	return nil
}

// GetString gets a config value as string
func GetString(key string) string {
	return viper.GetString(key)
}

// SetString sets a config value
func SetString(key, value string) {
	viper.Set(key, value)
}

// SaveConfig saves the current config to file
func SaveConfig() error {
	// Try to write config first
	err := viper.WriteConfig()
	if err != nil {
		// If config doesn't exist, use SafeWriteConfig to create it
		return viper.SafeWriteConfig()
	}
	return nil
}

// WriteConfigAs saves config to a specific file
func WriteConfigAs(filename string) error {
	return viper.WriteConfigAs(filename)
}
