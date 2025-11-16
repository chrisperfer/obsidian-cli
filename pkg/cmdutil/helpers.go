package cmdutil

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/client"
	"github.com/chrisperfer/obsidian-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetClient returns a configured API client
func GetClient() (*client.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return client.New(cfg.API.URL, cfg.API.Key, cfg.API.Timeout, cfg.API.Insecure), nil
}

// GetOutputFormat returns the output format from flags or config
func GetOutputFormat(cmd *cobra.Command) string {
	if cmd != nil {
		if outputFormat, err := cmd.Flags().GetString("output"); err == nil && outputFormat != "" {
			return outputFormat
		}
	}
	return viper.GetString("defaults.output")
}

// GetLLMMode returns whether LLM mode is enabled
func GetLLMMode(cmd *cobra.Command) bool {
	if cmd != nil {
		if llmMode, err := cmd.Flags().GetBool("llm-mode"); err == nil && llmMode {
			return true
		}
	}
	return viper.GetBool("llm.mode")
}
