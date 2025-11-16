package client

import (
	"fmt"
)

// ListPeriodicNotes lists periodic notes
func (c *Client) ListPeriodicNotes(noteType string) ([]string, error) {
	var result struct {
		Files []string `json:"files"`
	}

	resp, err := c.client.R().
		SetResult(&result).
		SetError(&APIError{}).
		Get("/periodic/" + noteType)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return result.Files, nil
}

// GetPeriodicNote retrieves a specific periodic note
func (c *Client) GetPeriodicNote(noteType, date string) (*Note, error) {
	var result Note
	endpoint := "/periodic/" + noteType

	req := c.client.R().
		SetResult(&result).
		SetError(&APIError{})

	if date != "" {
		req.SetQueryParam("date", date)
	}

	resp, err := req.Get(endpoint)

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

// CreatePeriodicNote creates a periodic note
func (c *Client) CreatePeriodicNote(noteType, date, content string) (*Note, error) {
	var result Note
	endpoint := "/periodic/" + noteType

	body := map[string]interface{}{
		"content": content,
	}
	if date != "" {
		body["date"] = date
	}

	resp, err := c.client.R().
		SetBody(body).
		SetResult(&result).
		SetError(&APIError{}).
		Post(endpoint)

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
