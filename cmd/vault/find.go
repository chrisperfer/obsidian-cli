package vault

import (
	"fmt"
	"path/filepath"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	findPattern string
)

var findCmd = &cobra.Command{
	Use:   "find [path]",
	Short: "Find all notes recursively",
	Long:  `Find all notes in the vault recursively. Always recursive. Optionally filter by path and pattern.`,
	Example: `  # Find all notes
  obsidian-cli vault find

  # Find all notes in a folder
  obsidian-cli vault find "Daily/"

  # Find with pattern matching
  obsidian-cli vault find --pattern "*.md"
  obsidian-cli vault find --pattern "2025-*.md"

  # Combine path and pattern
  obsidian-cli vault find "Daily/" --pattern "2025-11-*.md"

  # Find in JSON format
  obsidian-cli vault find --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		var basePath string
		if len(args) > 0 {
			basePath = args[0]
		}

		// Get all files recursively
		allFiles, err := client.ListAllFiles()
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		// Filter by base path and pattern
		filtered := []string{}
		for _, file := range allFiles {
			// Check base path
			if basePath != "" {
				if len(file) < len(basePath) || file[:len(basePath)] != basePath {
					continue
				}
			}

			// Check pattern
			if findPattern != "" {
				matched, err := filepath.Match(findPattern, filepath.Base(file))
				if err != nil {
					return fmt.Errorf("invalid pattern: %w", err)
				}
				if !matched {
					continue
				}
			}

			filtered = append(filtered, file)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(filtered)
	},
}

func init() {
	findCmd.Flags().StringVarP(&findPattern, "pattern", "p", "", "Filter files by glob pattern (e.g., '*.md', '2025-*.md')")
}
