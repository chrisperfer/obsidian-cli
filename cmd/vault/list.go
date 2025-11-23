package vault

import (
	"fmt"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	listRecursive bool
)

var listCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List items in the vault",
	Long:  `List items in the vault. By default, lists top-level items only. Use --recursive to list all files recursively.`,
	Example: `  # List top-level items
  obsidian-cli vault list

  # List items in a folder
  obsidian-cli vault list "Food and Drink/"

  # List all files recursively
  obsidian-cli vault list --recursive

  # List all files under a path recursively
  obsidian-cli vault list "Daily/" --recursive

  # List in JSON format
  obsidian-cli vault list --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		var files []string
		var basePath string

		if len(args) > 0 {
			basePath = args[0]
		}

		if listRecursive || basePath != "" {
			// Use search API to get all files
			allFiles, err := client.ListAllFiles()
			if err != nil {
				return fmt.Errorf("failed to list files: %w", err)
			}

			// Filter by base path if specified
			if basePath != "" {
				filtered := []string{}
				for _, file := range allFiles {
					if len(file) >= len(basePath) && file[:len(basePath)] == basePath {
						if listRecursive {
							filtered = append(filtered, file)
						} else {
							// Non-recursive: only show immediate children
							remainder := file[len(basePath):]
							// Check if there's a slash in the remainder (indicating deeper nesting)
							hasSlash := false
							for _, ch := range remainder {
								if ch == '/' {
									hasSlash = true
									break
								}
							}
							if !hasSlash {
								filtered = append(filtered, file)
							}
						}
					}
				}
				files = filtered
			} else {
				if listRecursive {
					files = allFiles
				} else {
					// This shouldn't happen (recursive without path should use vault API)
					// but handle it anyway
					files = allFiles
				}
			}
		} else {
			// Use standard vault API for top-level listing
			result, err := client.ListFiles()
			if err != nil {
				return fmt.Errorf("failed to list files: %w", err)
			}
			files = result.Files
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(files)
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listRecursive, "recursive", "r", false, "List files recursively")
}
