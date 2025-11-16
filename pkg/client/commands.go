package client

import (
	"fmt"
)

// ListCommands retrieves the list of available Obsidian commands
func (c *Client) ListCommands() ([]Command, error) {
	var result CommandList

	resp, err := c.client.R().
		SetResult(&result).
		SetError(&APIError{}).
		Get("/commands/")

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if !resp.IsSuccess() {
		if apiErr, ok := resp.Error().(*APIError); ok {
			return nil, apiErr
		}
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	return result.Commands, nil
}

// ExecuteCommand executes an Obsidian command
func (c *Client) ExecuteCommand(commandID string) error {
	req := CommandExecuteRequest{
		CommandID: commandID,
	}

	resp, err := c.client.R().
		SetBody(req).
		SetError(&APIError{}).
		Post("/commands/" + commandID)

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
