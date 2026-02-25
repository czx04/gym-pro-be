package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"gym-pro-2026-ptit/internal/config"
	"gym-pro-2026-ptit/internal/domain/user"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

// OAuthProvider defines interface for OAuth providers
type OAuthProvider interface {
	GetAuthURL(state string) string
	GetUserInfo(ctx context.Context, code string) (*user.OAuthUserInfo, error)
}

// GoogleOAuthProvider implements OAuth for Google
type GoogleOAuthProvider struct {
	config *oauth2.Config
}

// NewGoogleOAuthProvider creates a new Google OAuth provider
func NewGoogleOAuthProvider(cfg *config.OAuthProviderConfig) *GoogleOAuthProvider {
	// TODO: Initialize Google OAuth2 config
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthProvider{config: config}
}

// GetAuthURL returns the Google OAuth authorization URL
func (p *GoogleOAuthProvider) GetAuthURL(state string) string {
	// TODO: Generate auth URL with state parameter
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// GetUserInfo retrieves user information from Google
func (p *GoogleOAuthProvider) GetUserInfo(ctx context.Context, code string) (*user.OAuthUserInfo, error) {
	// TODO: 1. Exchange code for token
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// TODO: 2. Get user info from Google API
	client := p.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// TODO: 3. Parse response
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// TODO: 4. Convert to domain model
	return &user.OAuthUserInfo{
		Provider:  "google",
		ID:        googleUser.ID,
		Email:     googleUser.Email,
		Name:      googleUser.Name,
		AvatarURL: &googleUser.Picture,
	}, nil
}

// FacebookOAuthProvider implements OAuth for Facebook
type FacebookOAuthProvider struct {
	config *oauth2.Config
}

// NewFacebookOAuthProvider creates a new Facebook OAuth provider
func NewFacebookOAuthProvider(cfg *config.OAuthProviderConfig) *FacebookOAuthProvider {
	// TODO: Initialize Facebook OAuth2 config
	config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"email",
			"public_profile",
		},
		Endpoint: facebook.Endpoint,
	}

	return &FacebookOAuthProvider{config: config}
}

// GetAuthURL returns the Facebook OAuth authorization URL
func (p *FacebookOAuthProvider) GetAuthURL(state string) string {
	// TODO: Generate auth URL with state parameter
	return p.config.AuthCodeURL(state)
}

// GetUserInfo retrieves user information from Facebook
func (p *FacebookOAuthProvider) GetUserInfo(ctx context.Context, code string) (*user.OAuthUserInfo, error) {
	// TODO: 1. Exchange code for token
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	// TODO: 2. Get user info from Facebook Graph API
	url := fmt.Sprintf("https://graph.facebook.com/me?fields=id,name,email,picture&access_token=%s", token.AccessToken)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// TODO: 3. Parse response
	var fbUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.Unmarshal(body, &fbUser); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	// TODO: 4. Convert to domain model
	avatarURL := fbUser.Picture.Data.URL
	return &user.OAuthUserInfo{
		Provider:  "facebook",
		ID:        fbUser.ID,
		Email:     fbUser.Email,
		Name:      fbUser.Name,
		AvatarURL: &avatarURL,
	}, nil
}

// TODO: Implement OAuth use cases:
// - OAuthLoginUseCase: Handle OAuth login flow
//   1. Check if user exists by oauth_provider + oauth_id
//   2. If exists: generate tokens and return
//   3. If not exists:
//      a. Check if user with same email exists
//      b. If yes: link OAuth account to existing user
//      c. If no: create new user with OAuth info
//   4. Generate and return tokens
