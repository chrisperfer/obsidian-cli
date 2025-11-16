package cmd

import (
	"fmt"
	"os"

	"github.com/chrisperfer/obsidian-cli/cmd/commands"
	"github.com/chrisperfer/obsidian-cli/cmd/config"
	"github.com/chrisperfer/obsidian-cli/cmd/note"
	"github.com/chrisperfer/obsidian-cli/cmd/periodic"
	"github.com/chrisperfer/obsidian-cli/cmd/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile      string
	outputFormat string
	llmMode      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "obsidian-cli",
	Short: "A CLI tool to interact with Obsidian via REST API",
	Long: `obsidian-cli is a command-line interface for interacting with your Obsidian vault
through the Local REST API plugin. It provides progressive disclosure of capabilities
and is optimized for use with LLMs.`,
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.obsidian-cli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "output format (text|json|yaml)")
	rootCmd.PersistentFlags().BoolVar(&llmMode, "llm-mode", false, "enable LLM-optimized output")

	// Add subcommands
	rootCmd.AddCommand(vault.VaultCmd)
	rootCmd.AddCommand(note.NoteCmd)
	rootCmd.AddCommand(periodic.PeriodicCmd)
	rootCmd.AddCommand(commands.CommandsCmd)
	rootCmd.AddCommand(config.ConfigCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".obsidian-cli")
	}

	viper.SetEnvPrefix("OBSIDIAN")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("api.url", "http://localhost:27124")
	viper.SetDefault("api.timeout", "30s")
	viper.SetDefault("api.insecure", false)
	viper.SetDefault("defaults.output", "text")
	viper.SetDefault("defaults.confirm_deletes", true)
	viper.SetDefault("llm.mode", false)
	viper.SetDefault("llm.max_content_length", 10000)

	if err := viper.ReadInConfig(); err == nil {
		// Config file found and loaded successfully
	}
}
