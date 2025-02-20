package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mpc/internal/model"
	"net/http"
)

type GoogleOAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// GetOAuthTokens exchanges the authorization code for OAuth tokens
func (c *GoogleOAuthClient) GetOAuthTokens(code string) (*model.GoogleTokenResponse, error) {
	// Create form data
	formData := map[string]string{
		"code":          code,
		"client_id":     c.ClientID,
		"client_secret": c.ClientSecret,
		"redirect_uri":  c.RedirectURI,
		"grant_type":    "authorization_code",
	}

	// Convert form data to JSON
	jsonData, err := json.Marshal(formData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal form data: %v", err)
	}

	// Make POST request to Google OAuth token endpoint
	resp, err := http.Post("https://oauth2.googleapis.com/token", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to make POST request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var tokenResponse model.GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token response: %v", err)
	}

	return &tokenResponse, nil
}

// GetUserInfo retrieves user information using the access token
func (c *GoogleOAuthClient) GetUserInfo(accessToken string) (*model.GoogleUserInfo, error) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %v", err)
	}

	// Add Authorization header
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var userInfo model.GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return &userInfo, nil
}
