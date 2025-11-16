package note

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	updateContent string
	updateFromStdin bool
)

var updateCmd = &cobra.Command{
	Use:   "update [path]",
	Short: "Update an existing note",
	Long:  `Update the content of an existing note in your Obsidian vault`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Update a note with new content
  obsidian-cli note update "Ideas.md" --content "# Updated Ideas"

  # Update a note from stdin
  echo "# New Content" | obsidian-cli note update "Notes.md" --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		content := updateContent

		// Read from stdin if requested
		if updateFromStdin {
			reader := bufio.NewReader(os.Stdin)
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(data)
		}

		if content == "" && !updateFromStdin {
			return fmt.Errorf("content is required (use --content or --stdin)")
		}

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		note, err := client.UpdateNote(path, content)
		if err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Updated note: %s", path))
		return formatter.Print(note)
	},
}

func init() {
	updateCmd.Flags().StringVarP(&updateContent, "content", "c", "", "new note content")
	updateCmd.Flags().BoolVar(&updateFromStdin, "stdin", false, "read content from stdin")
}
