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

// GetVaultInfo retrieves vault information
func (c *Client) GetVaultInfo() (*VaultInfo, error) {
	var result VaultInfo
	resp, err := c.client.R().
		SetResult(&result).
		SetError(&APIError{}).
		Get("/vault/info")

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
