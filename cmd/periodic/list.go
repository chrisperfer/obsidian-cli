package periodic

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List periodic notes",
	Long:  `List periodic notes of a specific type (daily, weekly, monthly)`,
	Args:  cobra.ExactArgs(1),
	ValidArgs: []string{"daily", "weekly", "monthly"},
	Example: `  # List daily notes
  obsidian-cli periodic list daily

  # List weekly notes in JSON format
  obsidian-cli periodic list weekly --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		noteType := args[0]

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		files, err := client.ListPeriodicNotes(noteType)
		if err != nil {
			return fmt.Errorf("failed to list periodic notes: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(files)
	},
}
