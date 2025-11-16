package periodic

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	createDate    string
	createContent string
)

var createCmd = &cobra.Command{
	Use:   "create [type]",
	Short: "Create a periodic note",
	Long:  `Create a new periodic note (daily, weekly, monthly)`,
	Args:  cobra.ExactArgs(1),
	ValidArgs: []string{"daily", "weekly", "monthly"},
	Example: `  # Create today's daily note
  obsidian-cli periodic create daily

  # Create daily note for a specific date
  obsidian-cli periodic create daily --date 2024-01-15 --content "# Daily Note"

  # Create weekly note
  obsidian-cli periodic create weekly --content "# Week 3"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		noteType := args[0]

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		note, err := client.CreatePeriodicNote(noteType, createDate, createContent)
		if err != nil {
			return fmt.Errorf("failed to create periodic note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Created %s note", noteType))
		return formatter.Print(note)
	},
}

func init() {
	createCmd.Flags().StringVarP(&createDate, "date", "d", "", "specific date (YYYY-MM-DD)")
	createCmd.Flags().StringVarP(&createContent, "content", "c", "", "note content")
}
