package note

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	frontmatterOnly bool
	contentOnly     bool
)

var getCmd = &cobra.Command{
	Use:   "get [path]",
	Short: "Read note content",
	Long:  `Read the content of a note from your Obsidian vault`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Get a note by path
  obsidian-cli note get "Daily/2024-01-15.md"

  # Output as JSON
  obsidian-cli note get "Ideas.md" --output json

  # Get frontmatter only
  obsidian-cli note get "Project.md" --frontmatter-only`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]

		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		note, err := client.GetNote(path)
		if err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))

		// Handle special display modes
		if frontmatterOnly {
			return formatter.Print(note.Frontmatter)
		}

		if contentOnly {
			// Truncate if in LLM mode
			maxLen := viper.GetInt("llm.max_content_length")
			content := formatter.TruncateContent(note.Content, maxLen)
			fmt.Println(content)
			return nil
		}

		return formatter.Print(note)
	},
}

func init() {
	getCmd.Flags().BoolVar(&frontmatterOnly, "frontmatter-only", false, "show only frontmatter")
	getCmd.Flags().BoolVar(&contentOnly, "content-only", false, "show only content")
}
