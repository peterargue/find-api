package flow

import (
	"context"
	"net/http"
	"net/url"
)

// Client is an interface for making HTTP requests to the API
type Client interface {
	DoRequest(ctx context.Context, method, path string, query url.Values) (*http.Response, error)
	DecodeResponse(resp *http.Response, v any) error
}

// Service handles operations for the Flow API endpoints
type Service struct {
	client Client
}

// NewService creates a new Flow API service
func NewService(client Client) *Service {
	return &Service{client: client}
}
