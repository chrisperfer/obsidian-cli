package config

import (
	"fmt"
	"strings"

	"github.com/chrisperfer/obsidian-cli/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current obsidian-cli configuration`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		// Mask the API key for security
		maskedCfg := *cfg
		if maskedCfg.API.Key != "" {
			maskedCfg.API.Key = strings.Repeat("*", len(maskedCfg.API.Key))
		}

		data, err := yaml.Marshal(&maskedCfg)
		if err != nil {
			return fmt.Errorf("failed to format configuration: %w", err)
		}

		fmt.Println(string(data))
		return nil
	},
}
