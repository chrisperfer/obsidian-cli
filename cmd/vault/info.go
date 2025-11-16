package vault

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display vault information",
	Long:  `Show information about your Obsidian vault including name, path, and file count`,
	Example: `  # Show vault info
  obsidian-cli vault info

  # Show vault info in JSON format
  obsidian-cli vault info --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		info, err := client.GetVaultInfo()
		if err != nil {
			return fmt.Errorf("failed to get vault info: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(info)
	},
}
