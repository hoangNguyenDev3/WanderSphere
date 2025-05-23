package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

// APIClient handles HTTP requests to the WanderSphere API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
	CookieJar  http.CookieJar
}

// NewAPIClient creates a new API client with cookie support
func NewAPIClient() (*APIClient, error) {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:19003/api/v1"
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		CookieJar: jar,
	}, nil
}

// APIResponse represents a generic API response
type APIResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// Request makes an HTTP request with optional JSON body
func (c *APIClient) Request(method, path string, body interface{}) (*APIResponse, error) {
	var bodyReader io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	fullURL := c.BaseURL + path
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

// GET makes a GET request
func (c *APIClient) GET(path string) (*APIResponse, error) {
	return c.Request("GET", path, nil)
}

// POST makes a POST request with JSON body
func (c *APIClient) POST(path string, body interface{}) (*APIResponse, error) {
	return c.Request("POST", path, body)
}

// PUT makes a PUT request with JSON body
func (c *APIClient) PUT(path string, body interface{}) (*APIResponse, error) {
	return c.Request("PUT", path, body)
}

// DELETE makes a DELETE request
func (c *APIClient) DELETE(path string) (*APIResponse, error) {
	return c.Request("DELETE", path, nil)
}

// ParseJSON parses the response body as JSON
func (r *APIResponse) ParseJSON(target interface{}) error {
	if err := json.Unmarshal(r.Body, target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}
	return nil
}

// GetStringBody returns the response body as a string
func (r *APIResponse) GetStringBody() string {
	return string(r.Body)
}

// IsSuccess checks if the response status code indicates success (2xx)
func (r *APIResponse) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError checks if the response status code indicates client error (4xx)
func (r *APIResponse) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError checks if the response status code indicates server error (5xx)
func (r *APIResponse) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// GetCookies returns all cookies from the client's cookie jar for the base URL
func (c *APIClient) GetCookies() []*http.Cookie {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil
	}
	return c.CookieJar.Cookies(baseURL)
}

// GetSessionID returns the session cookie value if it exists
func (c *APIClient) GetSessionID() string {
	cookies := c.GetCookies()
	for _, cookie := range cookies {
		if cookie.Name == "wandersphere_session" {
			return cookie.Value
		}
	}
	return ""
}

// HasValidSession checks if the client has a valid session cookie
func (c *APIClient) HasValidSession() bool {
	return c.GetSessionID() != ""
}

// ClearCookies clears all cookies from the client
func (c *APIClient) ClearCookies() error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return fmt.Errorf("failed to create new cookie jar: %w", err)
	}
	c.CookieJar = jar
	c.HTTPClient.Jar = jar
	return nil
}
