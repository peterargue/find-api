package auth

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

// Client is an interface for making HTTP requests to the API
type Client interface {
	DoRequestWithBasicAuth(ctx context.Context, method, path string, query url.Values, username, password string) (*http.Response, error)
	DecodeResponse(resp *http.Response, v any) error
}

// Service handles Auth API operations
type Service struct {
	client   Client
	username string
	password string
}

// NewService creates a new Auth API service
func NewService(client Client, username, password string) *Service {
	return &Service{
		client:   client,
		username: username,
		password: password,
	}
}

// TokenResponse represents the response from the JWT generation endpoint
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
	Exp          int64  `json:"exp"`
	Iat          int64  `json:"iat"`
}

// GenerateToken generates a new JWT token using Basic Auth
// expiry: Duration for the token validity (e.g., 10*time.Minute, 1*time.Hour, max 168*time.Hour)
func (s *Service) GenerateToken(ctx context.Context, expiry time.Duration) (*TokenResponse, error) {
	query := url.Values{}
	query.Set("expiry", expiry.String())

	resp, err := s.client.DoRequestWithBasicAuth(ctx, http.MethodPost, "/auth/v1/generate", query, s.username, s.password)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := s.client.DecodeResponse(resp, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}
