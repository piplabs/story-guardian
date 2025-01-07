package httpclient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultRequestTimeout = 60 * time.Second // Define a reasonable HTTP request timeout

	ContentTypeHeader   = "Content-Type"
	AuthorizationHeader = "Authorization"
	ContentTypeJSON     = "application/json"
)

// Client is a wrapper around http.Client to enforce best practices, like timeout and context usage.
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new instance of Client with a default timeout.
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout, // Enforce a global timeout for all requests
		},
	}
}

// DefaultClient creates a new instance of Client with a default request timeout.
func DefaultClient() *Client {
	return NewClient(defaultRequestTimeout)
}

// Do send an HTTP request and returns an HTTP response, handling context-related cancellation or deadline exceeded errors.
// It automatically handles requests with a given `context.Context`.
func (c *Client) Do(ctx context.Context, method, url string, body io.Reader, header map[string]string) (*http.Response, error) {
	// Create an HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	// Set default Content-Type header if not provided
	if body != nil && req.Header.Get(ContentTypeHeader) == "" {
		req.Header.Set(ContentTypeHeader, ContentTypeJSON)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Handle specific context-related errors
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("request canceled due to context cancellation: %w", err)
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("request failed due to context deadline exceeded: %w", err)
		}
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	// Check for unexpected response statuses (example: return an error for 5xx responses, etc.)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return resp, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	return resp, nil
}
