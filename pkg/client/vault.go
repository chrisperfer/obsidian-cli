package client

import (
	"fmt"
)

// ListFiles lists all files in the vault
func (c *Client) ListFiles() (*FileList, error) {
	var result FileList
	resp, err := c.client.R().
		SetResult(&result).
		SetError(&APIError{}).
		Get("/vault/")

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

// GetVaultInfo retrieves vault information from the root endpoint
func (c *Client) GetVaultInfo() (*VaultInfo, error) {
	var result VaultInfo
	resp, err := c.client.R().
		SetResult(&result).
		SetError(&APIError{}).
		Get("/")

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

// ListAllFiles retrieves all files recursively using the search API
func (c *Client) ListAllFiles() ([]string, error) {
	// Use search with empty query to get all files
	results, err := c.SearchNotes("")
	if err != nil {
		return nil, err
	}

	files := make([]string, len(results))
	for i, result := range results {
		files[i] = result.Filename
	}

	return files, nil
}
