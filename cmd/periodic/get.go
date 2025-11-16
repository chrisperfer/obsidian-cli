package periodic

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var getDate string

var getCmd = &cobra.Command{
	Use:   "get [type]",
	Short: "Get a specific periodic note",
	Long:  `Retrieve a periodic note for a specific date (daily, weekly, monthly)`,
	Args:  cobra.ExactArgs(1),
	ValidArgs: []string{"daily", "weekly", "monthly"},
	Example: `  # Get today's daily note
  obsidian-cli periodic get daily

  # Get daily note for a specific date
  obsidian-cli periodic get daily --date 2024-01-15

  # Get weekly note
  obsidian-cli periodic get weekly --date 2024-01-15`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		noteType := args[0]

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		note, err := client.GetPeriodicNote(noteType, getDate)
		if err != nil {
			return fmt.Errorf("failed to get periodic note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(note)
	},
}

func init() {
	getCmd.Flags().StringVarP(&getDate, "date", "d", "", "specific date (YYYY-MM-DD)")
}
