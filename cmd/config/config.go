package config

import (
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure CLI settings",
	Long:  `Initialize and manage obsidian-cli configuration`,
}

func init() {
	ConfigCmd.AddCommand(initCmd)
	ConfigCmd.AddCommand(showCmd)
}
