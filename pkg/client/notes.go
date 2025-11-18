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

// PatchNote performs a targeted update on a note (append/prepend/replace at heading/block/frontmatter)
func (c *Client) PatchNote(req *PatchNoteRequest) error {
	encodedPath := url.PathEscape(req.Path)

	request := c.client.R().
		SetHeader("Content-Type", "text/markdown").
		SetHeader("Operation", string(req.Operation)).
		SetHeader("Target-Type", string(req.TargetType)).
		SetHeader("Target", req.Target).
		SetBody(req.Content).
		SetError(&APIError{})

	// Optional headers
	if req.TargetDelimiter != "" {
		request.SetHeader("Target-Delimiter", req.TargetDelimiter)
	}
	if req.TrimTargetWhitespace {
		request.SetHeader("Trim-Target-Whitespace", "true")
	}
	if req.CreateTargetIfMissing {
		request.SetHeader("Create-Target-If-Missing", "true")
	}

	resp, err := request.Patch("/vault/" + encodedPath)

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

// AppendNote appends content to the end of an existing note
func (c *Client) AppendNote(path, content string) error {
	encodedPath := url.PathEscape(path)

	resp, err := c.client.R().
		SetHeader("Content-Type", "text/markdown").
		SetBody(content).
		SetError(&APIError{}).
		Post("/vault/" + encodedPath)

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
