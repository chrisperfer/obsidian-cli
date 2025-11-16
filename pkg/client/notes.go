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
	var result Note
	encodedPath := url.PathEscape(path)

	resp, err := c.client.R().
		SetBody(map[string]interface{}{
			"content": content,
		}).
		SetResult(&result).
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

	return &result, nil
}

// UpdateNote updates an existing note
func (c *Client) UpdateNote(path, content string) (*Note, error) {
	var result Note
	encodedPath := url.PathEscape(path)

	resp, err := c.client.R().
		SetBody(map[string]interface{}{
			"content": content,
		}).
		SetResult(&result).
		SetError(&APIError{}).
		Patch("/vault/" + encodedPath)

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

// SearchNotes searches for notes by query
func (c *Client) SearchNotes(query string) (*SearchResults, error) {
	var result SearchResults

	resp, err := c.client.R().
		SetQueryParam("query", query).
		SetResult(&result).
		SetError(&APIError{}).
		Get("/search/")

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
