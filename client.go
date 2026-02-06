package findapi

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/peterargue/find-api/auth"
	"github.com/peterargue/find-api/flow"
	"github.com/peterargue/find-api/simple"
)

const (
	FindApiURL = "https://api.find.xyz"
)

// Client is the main client for interacting with the FindLabs API
type Client struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string

	// JWT token management
	tokenMu     sync.RWMutex
	accessToken string
	tokenExpiry time.Time

	// Services
	Simple *simple.Service
	Auth   *auth.Service
	Flow   *flow.Service
}

// ClientOption is a functional option for configuring the Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL for the API
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new FindLabs API client
func NewClient(username, password string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:  FindApiURL,
		username: username,
		password: password,
	}

	// Apply options
	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.Simple = simple.NewService(c)
	c.Auth = auth.NewService(c, username, password)
	c.Flow = flow.NewService(c)

	return c
}

// DoRequest performs an HTTP request with automatic authentication and rate limiting handling
// This method is exported to allow service packages to make requests
func (c *Client) DoRequest(ctx context.Context, method, path string, query url.Values) (*http.Response, error) {
	return c.doRequest(ctx, method, path, query, nil)
}

// DoRequestWithBasicAuth performs an HTTP request with Basic Auth (used by auth service)
// This method is exported to allow the auth service to make requests without JWT
func (c *Client) DoRequestWithBasicAuth(ctx context.Context, method, path string, query url.Values, username, password string) (*http.Response, error) {
	// Build URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Basic Auth header
	authStr := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))
	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	return resp, nil
}

// doRequest performs an HTTP request with automatic authentication and rate limiting handling
func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body io.Reader) (*http.Response, error) {
	// Build URL
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if query != nil {
		u.RawQuery = query.Encode()
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication token (skip for auth endpoints)
	if path != "/auth/v1/generate" {
		token, err := c.getValidToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get valid token: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Execute request with retry logic for rate limiting
	var resp *http.Response
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		resp, err = c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		// Handle rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := c.getRetryAfter(resp)
			if i < maxRetries-1 {
				resp.Body.Close()
				select {
				case <-time.After(retryAfter):
					continue
				case <-ctx.Done():
					return nil, ctx.Err()
				}
			}
			// Last retry exhausted
			defer resp.Body.Close()
			return nil, &RateLimitError{RetryAfter: retryAfter}
		}

		// Success or non-rate-limit error
		break
	}

	return resp, nil
}

// getValidToken returns a valid JWT token, refreshing if necessary
func (c *Client) getValidToken(ctx context.Context) (string, error) {
	c.tokenMu.RLock()
	token := c.accessToken
	expiry := c.tokenExpiry
	c.tokenMu.RUnlock()

	// Check if token is still valid (with 1 minute buffer)
	if token != "" && time.Now().Add(time.Minute).Before(expiry) {
		return token, nil
	}

	// Need to refresh token
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()

	// Double-check after acquiring write lock
	if c.accessToken != "" && time.Now().Add(time.Minute).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	// Generate new token
	tokenResp, err := c.Auth.GenerateToken(ctx, 10*time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	c.accessToken = tokenResp.AccessToken
	c.tokenExpiry = time.Unix(tokenResp.Exp, 0)

	return c.accessToken, nil
}

// getRetryAfter extracts the retry-after duration from response headers
func (c *Client) getRetryAfter(resp *http.Response) time.Duration {
	retryAfter := resp.Header.Get("Retry-After")
	if retryAfter == "" {
		// Default to 1 second if no header present
		return time.Second
	}

	// Try parsing as seconds
	if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
		return seconds
	}

	// Try parsing as HTTP date
	if t, err := http.ParseTime(retryAfter); err == nil {
		return time.Until(t)
	}

	// Default fallback
	return time.Second
}

// DecodeResponse decodes a JSON response into the provided interface
// This method is exported to allow service packages to decode responses
func (c *Client) DecodeResponse(resp *http.Response, v any) error {
	return c.decodeResponse(resp, v)
}

// decodeResponse decodes a JSON response into the provided interface
func (c *Client) decodeResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
