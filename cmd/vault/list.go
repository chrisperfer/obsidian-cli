package vault

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes in the vault",
	Long:  `Retrieve and display a list of all notes stored in your Obsidian vault`,
	Example: `  # List all notes
  obsidian-cli vault list

  # List notes in JSON format
  obsidian-cli vault list --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		files, err := client.ListFiles()
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(files.Files)
	},
}
