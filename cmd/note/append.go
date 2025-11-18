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
	appendContent    string
	appendFromStdin  bool
)

var appendCmd = &cobra.Command{
	Use:   "append [path]",
	Short: "Append content to the end of a note",
	Long:  `Append content to the end of an existing note. Creates the note if it doesn't exist.`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Append text to a note
  obsidian-cli note append "Journal.md" --content "\n## New Entry\nToday's notes..."

  # Append from stdin
  echo "Additional content" | obsidian-cli note append "Notes.md" --stdin`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		content := appendContent

		// Read from stdin if requested
		if appendFromStdin {
			reader := bufio.NewReader(os.Stdin)
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(data)
		}

		if content == "" && !appendFromStdin {
			return fmt.Errorf("content is required (use --content or --stdin)")
		}

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		err = client.AppendNote(path, content)
		if err != nil {
			return fmt.Errorf("failed to append to note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Appended content to: %s", path))
		return nil
	},
}

func init() {
	appendCmd.Flags().StringVarP(&appendContent, "content", "c", "", "content to append")
	appendCmd.Flags().BoolVar(&appendFromStdin, "stdin", false, "read content from stdin")
}
