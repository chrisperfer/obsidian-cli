package note

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var forceDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Short: "Delete a note",
	Long:  `Remove a note from your Obsidian vault`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Delete a note (with confirmation)
  obsidian-cli note delete "Old.md"

  # Force delete without confirmation
  obsidian-cli note delete "Old.md" --force`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]

		// Confirm deletion unless --force is used
		if !forceDelete && viper.GetBool("defaults.confirm_deletes") {
			reader := bufio.NewReader(os.Stdin)
			fmt.Printf("Are you sure you want to delete '%s'? [y/N]: ", path)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))
			if response != "y" && response != "yes" {
				fmt.Println("Deletion cancelled")
				return nil
			}
		}

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		if err := client.DeleteNote(path); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Deleted note: %s", path))
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "skip confirmation prompt")
}
