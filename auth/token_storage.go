package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yourusername/robinstock_go/models"
)

func getTokenFilePath(username string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	tokensDir := filepath.Join(homeDir, ".tokens")
	os.MkdirAll(tokensDir, 0700)
	return filepath.Join(tokensDir, fmt.Sprintf("robinhood_%s.json", username))
}

func tryLoadStoredToken(username string) (*models.Auth, bool) {
	filePath := getTokenFilePath(username)
	if filePath == "" {
		return nil, false
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, false
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var token models.Auth
	if err := json.Unmarshal(content, &token); err != nil {
		return nil, false
	}

	if token.IsExpired() {
		os.Remove(filePath)
		return nil, false
	}

	return &token, true
}

func saveToken(username string, token *models.Auth) error {
	filePath := getTokenFilePath(username)
	if filePath == "" {
		return fmt.Errorf("could not determine home directory")
	}

	if token.IssuedAt.IsZero() {
		token.IssuedAt = time.Now()
	}

	fileData, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, fileData, 0600)
}

func deleteToken(username string) {
	filePath := getTokenFilePath(username)
	if filePath != "" {
		os.Remove(filePath)
	}
}

