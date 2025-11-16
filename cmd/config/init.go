package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chrisperfer/obsidian-cli/pkg/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  `Interactive setup wizard to configure obsidian-cli`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Obsidian CLI Configuration")
		fmt.Println("===========================")
		fmt.Println()

		cfg := config.GetDefaultConfig()

		// API URL
		fmt.Printf("API URL [%s]: ", cfg.API.URL)
		url, _ := reader.ReadString('\n')
		url = strings.TrimSpace(url)
		if url != "" {
			cfg.API.URL = url
		}

		// API Key
		fmt.Print("API Key (required): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("API key is required")
		}
		cfg.API.Key = key

		// Insecure mode
		fmt.Print("Skip TLS verification? [y/N]: ")
		insecure, _ := reader.ReadString('\n')
		insecure = strings.TrimSpace(strings.ToLower(insecure))
		cfg.API.Insecure = insecure == "y" || insecure == "yes"

		// Save configuration
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		home, _ := os.UserHomeDir()
		fmt.Printf("\n✓ Configuration saved to %s/.obsidian-cli.yaml\n", home)
		return nil
	},
}
