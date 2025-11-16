package commands

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var executeCmd = &cobra.Command{
	Use:   "execute [command-id]",
	Short: "Execute an Obsidian command",
	Long:  `Execute a specific Obsidian command by its ID`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Execute a command by ID
  obsidian-cli commands execute editor:toggle-bold

  # List commands to find IDs
  obsidian-cli commands list`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		commandID := args[0]

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		if err := client.ExecuteCommand(commandID); err != nil {
			return fmt.Errorf("failed to execute command: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Executed command: %s", commandID))
		return nil
	},
}
