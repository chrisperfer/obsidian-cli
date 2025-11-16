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
	createContent string
	createFromStdin bool
)

var createCmd = &cobra.Command{
	Use:   "create [path]",
	Short: "Create a new note",
	Long:  `Create a new note in your Obsidian vault`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Create a note with content
  obsidian-cli note create "Ideas.md" --content "# My Ideas"

  # Create a note from stdin
  echo "# Daily Notes" | obsidian-cli note create "Daily/2024-01-15.md" --stdin

  # Create an empty note
  obsidian-cli note create "Notes/Empty.md"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		content := createContent

		// Read from stdin if requested
		if createFromStdin {
			reader := bufio.NewReader(os.Stdin)
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(data)
		}

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		note, err := client.CreateNote(path, content)
		if err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Created note: %s", path))
		return formatter.Print(note)
	},
}

func init() {
	createCmd.Flags().StringVarP(&createContent, "content", "c", "", "note content")
	createCmd.Flags().BoolVar(&createFromStdin, "stdin", false, "read content from stdin")
}
