package commands

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available Obsidian commands",
	Long:  `Retrieve and display a list of all available Obsidian commands`,
	Example: `  # List all commands
  obsidian-cli commands list

  # List commands in JSON format
  obsidian-cli commands list --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		commands, err := client.ListCommands()
		if err != nil {
			return fmt.Errorf("failed to list commands: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))

		// Format as table for text output
		if cmdutil.GetOutputFormat(cobraCmd) == "text" {
			headers := []string{"ID", "Name"}
			rows := make([][]string, len(commands))
			for i, c := range commands {
				rows[i] = []string{c.ID, c.Name}
			}
			formatter.PrintTable(headers, rows)
			return nil
		}

		return formatter.Print(commands)
	},
}
