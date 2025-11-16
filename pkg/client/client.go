package client

import (
	"crypto/tls"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client represents the REST API client
type Client struct {
	baseURL  string
	apiKey   string
	client   *resty.Client
	insecure bool
}

// New creates a new API client
func New(baseURL, apiKey string, timeout time.Duration, insecure bool) *Client {
	restyClient := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Authorization", "Bearer "+apiKey).
		SetHeader("Content-Type", "application/json").
		SetTimeout(timeout)

	if insecure {
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}

	return &Client{
		baseURL:  baseURL,
		apiKey:   apiKey,
		client:   restyClient,
		insecure: insecure,
	}
}

// APIError represents an API error response
type APIError struct {
	Code        string   `json:"code"`
	Message     string   `json:"message"`
	Suggestions []string `json:"suggestions,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}
