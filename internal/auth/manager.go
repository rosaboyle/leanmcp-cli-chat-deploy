package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ddod/leanmcp-cli/internal/config"
)

// Credentials holds the authentication information
type Credentials struct {
	APIKey      string    `yaml:"api_key"`
	UserEmail   string    `yaml:"user_email,omitempty"`
	Scopes      []string  `yaml:"scopes,omitempty"`
	StoredAt    time.Time `yaml:"stored_at"`
	LastUsed    time.Time `yaml:"last_used,omitempty"`
}

// UserInfo represents user information returned from API
type UserInfo struct {
	Email  string   `json:"email"`
	Scopes []string `json:"scopes"`
}

// ValidateAPIKeyFormat checks if the API key has the correct format
func ValidateAPIKeyFormat(apiKey string) error {
	if apiKey == "" {
		return errors.New("API key cannot be empty")
	}
	
	if !strings.HasPrefix(apiKey, "airtrain_") {
		return errors.New("API key must start with 'airtrain_'")
	}
	
	if len(apiKey) < 20 {
		return errors.New("API key appears to be too short")
	}
	
	return nil
}

// StoreCredentials encrypts and stores the API key and user info
func StoreCredentials(apiKey string, userInfo *UserInfo) error {
	if err := ValidateAPIKeyFormat(apiKey); err != nil {
		return err
	}

	// Create credentials
	creds := &Credentials{
		APIKey:   apiKey,
		StoredAt: time.Now(),
	}
	
	if userInfo != nil {
		creds.UserEmail = userInfo.Email
		creds.Scopes = userInfo.Scopes
	}

	// Encrypt and store
	encryptedKey, err := encryptAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt API key: %v", err)
	}

	// Store in config
	config.SetString("api_key", encryptedKey)
	config.SetString("user_email", creds.UserEmail)
	config.SetString("stored_at", creds.StoredAt.Format(time.RFC3339))
	
	// Save scopes as comma-separated string
	if len(creds.Scopes) > 0 {
		config.SetString("scopes", strings.Join(creds.Scopes, ","))
	}

	return config.SaveConfig()
}

// LoadCredentials loads and decrypts the stored credentials
func LoadCredentials() (*Credentials, error) {
	encryptedKey := config.GetString("api_key")
	if encryptedKey == "" {
		return nil, errors.New("no stored credentials found")
	}

	// Decrypt API key
	apiKey, err := decryptAPIKey(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %v", err)
	}

	// Load other data
	creds := &Credentials{
		APIKey:    apiKey,
		UserEmail: config.GetString("user_email"),
	}

	// Parse stored_at
	if storedAtStr := config.GetString("stored_at"); storedAtStr != "" {
		if t, err := time.Parse(time.RFC3339, storedAtStr); err == nil {
			creds.StoredAt = t
		}
	}

	// Parse scopes
	if scopesStr := config.GetString("scopes"); scopesStr != "" {
		creds.Scopes = strings.Split(scopesStr, ",")
	}

	return creds, nil
}

// ClearCredentials removes stored credentials
func ClearCredentials() error {
	config.SetString("api_key", "")
	config.SetString("user_email", "")
	config.SetString("stored_at", "")
	config.SetString("scopes", "")
	return config.SaveConfig()
}

// UpdateLastUsed updates the last used timestamp
func UpdateLastUsed() error {
	config.SetString("last_used", time.Now().Format(time.RFC3339))
	return config.SaveConfig()
}

// encryptAPIKey encrypts the API key using AES
func encryptAPIKey(apiKey string) (string, error) {
	// Use a simple key derivation from machine info
	key := deriveKey()
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(apiKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptAPIKey decrypts the API key
func decryptAPIKey(encryptedKey string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return "", err
	}

	key := deriveKey()
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// deriveKey creates a key from machine-specific information
func deriveKey() []byte {
	hostname, _ := os.Hostname()
	data := fmt.Sprintf("leanmcp-cli-%s", hostname)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}
