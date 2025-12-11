// Package sevk provides the official Go SDK for Sevk - Email Marketing Platform
package sevk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.sevk.email/api"
	defaultTimeout = 30 * time.Second
)

// Client is the main entry point for the Sevk SDK
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	// Resources
	Contacts      *ContactsResource
	Audiences     *AudiencesResource
	Templates     *TemplatesResource
	Broadcasts    *BroadcastsResource
	Domains       *DomainsResource
	Topics        *TopicsResource
	Segments      *SegmentsResource
	Subscriptions *SubscriptionsResource
	Emails        *EmailsResource
}

// Options for configuring the Sevk client
type Options struct {
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// New creates a new Sevk client with the given API key
func New(apiKey string) *Client {
	return NewWithOptions(apiKey, Options{})
}

// NewWithOptions creates a new Sevk client with custom options
func NewWithOptions(apiKey string, opts Options) *Client {
	baseURL := defaultBaseURL
	if opts.BaseURL != "" {
		baseURL = opts.BaseURL
	}

	httpClient := opts.HTTPClient
	if httpClient == nil {
		timeout := defaultTimeout
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	c := &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	// Initialize resources
	c.Contacts = &ContactsResource{client: c}
	c.Audiences = &AudiencesResource{client: c}
	c.Templates = &TemplatesResource{client: c}
	c.Broadcasts = &BroadcastsResource{client: c}
	c.Domains = &DomainsResource{client: c}
	c.Topics = &TopicsResource{client: c}
	c.Segments = &SegmentsResource{client: c}
	c.Subscriptions = &SubscriptionsResource{client: c}
	c.Emails = &EmailsResource{client: c}

	return c
}

// request makes an HTTP request to the Sevk API
func (c *Client) request(method, path string, body interface{}, result interface{}) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return &Error{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		return &Error{
			StatusCode: resp.StatusCode,
			Message:    apiErr.Message,
			Code:       apiErr.Code,
		}
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

// get makes a GET request
func (c *Client) get(path string, result interface{}) error {
	return c.request(http.MethodGet, path, nil, result)
}

// post makes a POST request
func (c *Client) post(path string, body interface{}, result interface{}) error {
	return c.request(http.MethodPost, path, body, result)
}

// patch makes a PATCH request
func (c *Client) patch(path string, body interface{}, result interface{}) error {
	return c.request(http.MethodPatch, path, body, result)
}

// put makes a PUT request
func (c *Client) put(path string, body interface{}, result interface{}) error {
	return c.request(http.MethodPut, path, body, result)
}

// delete makes a DELETE request
func (c *Client) delete(path string) error {
	return c.request(http.MethodDelete, path, nil, nil)
}
