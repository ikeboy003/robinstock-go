package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/yourusername/robinstock_go"
	"github.com/yourusername/robinstock_go/models"
	"github.com/yourusername/robinstock_go/urls"
)

// Login authenticates with Robinhood and returns auth credentials.
func Login(ctx context.Context, client *robinstock_go.Client, username, password, mfaCode string) (*models.Auth, error) {
	if token, ok := tryLoadStoredToken(username); ok {
		client.SetAuth(token)
		return token, nil
	}

	deviceToken, err := generateDeviceToken()
	if err != nil {
		return nil, fmt.Errorf("generate device token: %w", err)
	}

	payload := map[string]string{
		"username":       username,
		"password":       password,
		"mfa_code":       mfaCode,
		"device_token":   deviceToken,
		"client_id":      models.ClientID,
		"grant_type":     "password",
		"scope":          "internal",
		"expires_in":     "86400",
		"challenge_type": "email",
	}

	resp, err := client.Post(ctx, urls.LoginURL(), payload, false)
	if err != nil {
		return nil, err
	}

	// Check for Sheriff verification workflow FIRST (can come with 403 status)
	if verificationWorkflow, ok := resp.Data["verification_workflow"].(map[string]interface{}); ok {
		workflowID := robinstock_go.GetString(verificationWorkflow, "id")
		if workflowID != "" {
			log.Println("Sheriff verification required, starting workflow...")
			if err := handleSheriffVerification(ctx, client, deviceToken, workflowID); err != nil {
				return nil, fmt.Errorf("sheriff verification failed: %w", err)
			}

			log.Println("Retrying login after Sheriff verification...")
			resp, err = client.Post(ctx, urls.LoginURL(), payload, false)
			if err != nil {
				return nil, fmt.Errorf("login after verification failed: %w", err)
			}
		}
	}

	// Check for MFA requirement
	if mfaRequired, ok := resp.Data["mfa_required"].(bool); ok && mfaRequired {
		return nil, fmt.Errorf("MFA required but not provided")
	}

	// Check for challenge
	if challenge, ok := resp.Data["challenge"].(map[string]interface{}); ok {
		challengeID, _ := challenge["id"].(string)
		return nil, fmt.Errorf("challenge required: %s", challengeID)
	}

	// Now check for error status codes
	if resp.StatusCode >= 400 {
		if detail, ok := resp.Data["detail"].(string); ok {
			return nil, fmt.Errorf("login failed: %s", detail)
		}
		return nil, fmt.Errorf("login failed: status %d", resp.StatusCode)
	}

	auth := &models.Auth{
		AccessToken:  robinstock_go.GetString(resp.Data, "access_token"),
		RefreshToken: robinstock_go.GetString(resp.Data, "refresh_token"),
		TokenType:    robinstock_go.GetString(resp.Data, "token_type"),
		DeviceToken:  deviceToken,
		ExpiresIn:    robinstock_go.GetInt(resp.Data, "expires_in"),
	}

	if auth.AccessToken == "" {
		return nil, fmt.Errorf("no access token in response")
	}

	client.SetAuth(auth)

	if err := saveToken(username, auth); err != nil {
		return nil, fmt.Errorf("failed to save token: %w", err)
	}

	return auth, nil
}

// RefreshToken refreshes an existing authentication token.
func RefreshToken(ctx context.Context, client *robinstock_go.Client, refreshToken, deviceToken string) (*models.Auth, error) {
	payload := map[string]string{
		"client_id":     models.ClientID,
		"grant_type":    "refresh_token",
		"refresh_token": refreshToken,
		"device_token":  deviceToken,
		"scope":         "internal",
	}

	resp, err := client.Post(ctx, urls.LoginURL(), payload, false)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("refresh failed: status %d", resp.StatusCode)
	}

	auth := &models.Auth{
		AccessToken:  robinstock_go.GetString(resp.Data, "access_token"),
		RefreshToken: robinstock_go.GetString(resp.Data, "refresh_token"),
		TokenType:    robinstock_go.GetString(resp.Data, "token_type"),
		DeviceToken:  deviceToken,
		ExpiresIn:    robinstock_go.GetInt(resp.Data, "expires_in"),
	}

	client.SetAuth(auth)
	return auth, nil
}

// Logout clears authentication and deletes stored token.
func Logout(username string, client *robinstock_go.Client) {
	deleteToken(username)
	client.SetAuth(nil)
}

// generateDeviceToken generates a unique device identifier.
func generateDeviceToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(bytes[0:4]),
		hex.EncodeToString(bytes[4:6]),
		hex.EncodeToString(bytes[6:8]),
		hex.EncodeToString(bytes[8:10]),
		hex.EncodeToString(bytes[10:16]),
	), nil
}
