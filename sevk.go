// Package sevk provides the official Go SDK for the sevk email platform.
package sevk

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// Version is the current version of the sevk Go SDK.
	Version = "1.1.0"

	defaultBaseURL = "https://api.sevk.io"
	defaultTimeout = 30 * time.Second
)

// Client is the main entry point for the sevk SDK.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client

	Contacts      *ContactsResource
	Audiences     *AudiencesResource
	Templates     *TemplatesResource
	Broadcasts    *BroadcastsResource
	Domains       *DomainsResource
	Topics        *TopicsResource
	Segments      *SegmentsResource
	Subscriptions *SubscriptionsResource
	Emails        *EmailsResource
	Webhooks      *WebhooksResource
	Activity      *ActivityResource
	Inbound       *InboundResource
	Outbound      *OutboundResource
	Usage         *UsageResource
}

// Options configures a Client.
type Options struct {
	BaseURL    string
	HTTPClient *http.Client
	Timeout    time.Duration
}

// New creates a client with the given API key.
func New(apiKey string) *Client {
	return NewWithOptions(apiKey, Options{})
}

// NewWithOptions creates a client with custom options.
func NewWithOptions(apiKey string, opts Options) *Client {
	baseURL := defaultBaseURL
	if opts.BaseURL != "" {
		baseURL = opts.BaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")

	httpClient := opts.HTTPClient
	if httpClient == nil {
		timeout := defaultTimeout
		if opts.Timeout > 0 {
			timeout = opts.Timeout
		}
		httpClient = &http.Client{Timeout: timeout}
	} else {
		copy := *httpClient
		if opts.Timeout > 0 {
			copy.Timeout = opts.Timeout
		}
		httpClient = &copy
	}
	httpClient.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	c := &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	c.Contacts = &ContactsResource{client: c}
	c.Audiences = &AudiencesResource{client: c}
	c.Templates = &TemplatesResource{client: c}
	c.Broadcasts = &BroadcastsResource{client: c}
	c.Domains = &DomainsResource{client: c}
	c.Topics = &TopicsResource{client: c}
	c.Segments = &SegmentsResource{client: c}
	c.Subscriptions = &SubscriptionsResource{client: c}
	c.Emails = &EmailsResource{client: c}
	c.Webhooks = &WebhooksResource{client: c}
	c.Activity = &ActivityResource{client: c}
	c.Inbound = &InboundResource{client: c}
	c.Outbound = &OutboundResource{client: c}
	c.Usage = &UsageResource{client: c}

	return c
}

func (c *Client) request(method, path string, body any, result any) error {
	return c.requestWithHeaders(method, path, body, nil, result)
}

func (c *Client) requestWithHeaders(method, path string, body any, extraHeaders map[string]string, result any) error {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return newRequestError(err)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return &Error{Message: "failed to marshal request body: " + err.Error(), Code: "SERIALIZATION_ERROR", cause: err}
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return newRequestError(err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sevk-go/"+Version)
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return newRequestError(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return newRequestError(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return parseAPIError(resp.StatusCode, resp.Header, respBody)
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}

	if result != nil {
		var root any
		if len(bytes.TrimSpace(respBody)) == 0 || json.Unmarshal(respBody, &root) != nil {
			return invalidResponseError(resp.StatusCode, respBody)
		}
		if _, ok := root.(map[string]any); !ok {
			return invalidResponseError(resp.StatusCode, respBody)
		}
		if err := json.Unmarshal(respBody, result); err != nil {
			decodeError := invalidResponseError(resp.StatusCode, respBody)
			decodeError.cause = err
			return decodeError
		}
	}

	return nil
}

func invalidResponseError(statusCode int, body []byte) *Error {
	return &Error{
		StatusCode:   statusCode,
		Message:      "expected a JSON object response",
		Code:         "INVALID_RESPONSE",
		ResponseBody: string(body),
	}
}

func (c *Client) get(path string, result any) error {
	return c.request(http.MethodGet, path, nil, result)
}

func (c *Client) post(path string, body any, result any) error {
	return c.request(http.MethodPost, path, body, result)
}

func (c *Client) postWithHeaders(path string, body any, headers map[string]string, result any) error {
	return c.requestWithHeaders(http.MethodPost, path, body, headers, result)
}

func (c *Client) patch(path string, body any, result any) error {
	return c.request(http.MethodPatch, path, body, result)
}

func (c *Client) put(path string, body any, result any) error {
	return c.request(http.MethodPut, path, body, result)
}

func (c *Client) delete(path string, result any) error {
	return c.request(http.MethodDelete, path, nil, result)
}
