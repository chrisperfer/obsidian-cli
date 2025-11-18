package client

import (
	"fmt"
	"net/url"
)

// GetNote retrieves a note by path
func (c *Client) GetNote(path string) (*Note, error) {
	var result Note
	encodedPath := url.PathEscape(path)
	resp, err := c.client.R().
		SetHeader("Accept", "application/vnd.olrapi.note+json").
		SetResult(&result).
		SetError(&APIError{}).
		Get("/vault/" + encodedPath)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return &result, nil
}

// CreateNote creates a new note
func (c *Client) CreateNote(path, content string) (*Note, error) {
	encodedPath := url.PathEscape(path)

	resp, err := c.client.R().
		SetHeader("Content-Type", "text/markdown").
		SetBody(content).
		SetError(&APIError{}).
		Put("/vault/" + encodedPath)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	// PUT returns 204 No Content on success, so return a minimal Note with just the path
	return &Note{Path: path}, nil
}

// UpdateNote updates an existing note (uses PUT for full content replacement)
func (c *Client) UpdateNote(path, content string) (*Note, error) {
	encodedPath := url.PathEscape(path)

	resp, err := c.client.R().
		SetHeader("Content-Type", "text/markdown").
		SetBody(content).
		SetError(&APIError{}).
		Put("/vault/" + encodedPath)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	// PUT returns 204 No Content on success, so return a minimal Note with just the path
	return &Note{Path: path}, nil
}

// DeleteNote deletes a note
func (c *Client) DeleteNote(path string) error {
	encodedPath := url.PathEscape(path)

	resp, err := c.client.R().
		SetError(&APIError{}).
		Delete("/vault/" + encodedPath)

	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return apiErr
		}
		return fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return nil
}

// SearchNotes searches for notes by query using the simple search endpoint
func (c *Client) SearchNotes(query string) ([]SearchResult, error) {
	var results []SearchResult

	resp, err := c.client.R().
		SetQueryParam("query", query).
		SetResult(&results).
		SetError(&APIError{}).
		Post("/search/simple/")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return results, nil
}
