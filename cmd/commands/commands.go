package commands

import (
	"github.com/spf13/cobra"
)

// CommandsCmd represents the commands command
var CommandsCmd = &cobra.Command{
	Use:   "commands",
	Short: "Execute Obsidian commands",
	Long:  `List and execute Obsidian commands programmatically`,
}

func init() {
	CommandsCmd.AddCommand(listCmd)
	CommandsCmd.AddCommand(executeCmd)
}
