package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	API struct {
		URL      string        `yaml:"url" mapstructure:"url"`
		Key      string        `yaml:"key" mapstructure:"key"`
		Insecure bool          `yaml:"insecure" mapstructure:"insecure"`
		Timeout  time.Duration `yaml:"timeout" mapstructure:"timeout"`
	} `yaml:"api" mapstructure:"api"`

	Defaults struct {
		Output         string `yaml:"output" mapstructure:"output"`
		VaultPath      string `yaml:"vault_path" mapstructure:"vault_path"`
		ConfirmDeletes bool   `yaml:"confirm_deletes" mapstructure:"confirm_deletes"`
	} `yaml:"defaults" mapstructure:"defaults"`

	LLM struct {
		Mode             bool `yaml:"mode" mapstructure:"mode"`
		MaxContentLength int  `yaml:"max_content_length" mapstructure:"max_content_length"`
	} `yaml:"llm" mapstructure:"llm"`
}

// Load reads the configuration from viper
func Load() (*Config, error) {
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.API.Key == "" {
		return nil, fmt.Errorf("API key not configured. Run 'obsidian-cli config init' to set up")
	}

	return &cfg, nil
}

// Save writes the configuration to the config file
func Save(cfg *Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := home + "/.obsidian-cli.yaml"
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetDefaultConfig returns a default configuration
func GetDefaultConfig() *Config {
	cfg := &Config{}
	cfg.API.URL = "http://localhost:27124"
	cfg.API.Timeout = 30 * time.Second
	cfg.API.Insecure = false
	cfg.Defaults.Output = "text"
	cfg.Defaults.ConfirmDeletes = true
	cfg.LLM.Mode = false
	cfg.LLM.MaxContentLength = 10000
	return cfg
}
