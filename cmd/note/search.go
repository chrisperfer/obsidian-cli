package note

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search notes by content",
	Long:  `Search for notes in your Obsidian vault by content`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Search for notes containing "project"
  obsidian-cli note search "project"

  # Search with JSON output
  obsidian-cli note search "meeting" --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		query := args[0]

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		results, err := client.SearchNotes(query)
		if err != nil {
			return fmt.Errorf("failed to search notes: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(results)
	},
}
