package note

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/chrisperfer/obsidian-cli/pkg/client"
	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	patchContent             string
	patchFromStdin           bool
	patchOperation           string
	patchTargetType          string
	patchTarget              string
	patchTargetDelimiter     string
	patchTrimWhitespace      bool
	patchCreateIfMissing     bool
	patchArrayAdd            string
	patchArrayRemove         string
)

var patchCmd = &cobra.Command{
	Use:   "patch [path]",
	Short: "Perform targeted updates on a note",
	Long:  `Perform targeted updates on a note by appending, prepending, or replacing content at specific headings, blocks, or frontmatter fields`,
	Args:  cobra.ExactArgs(1),
	Example: `  # Append content under a heading
  obsidian-cli note patch "Notes.md" --target "TODO" --target-type heading --operation append --content "- New task"

  # Prepend to a nested heading
  obsidian-cli note patch "Project.md" --target "Planning::Tasks" --target-type heading --operation prepend --content "## High Priority"

  # Update a frontmatter field
  obsidian-cli note patch "Document.md" --target "status" --target-type frontmatter --operation replace --content "completed"

  # Add a value to a frontmatter array
  obsidian-cli note patch "Document.md" --target "tags" --array-add "new-tag"

  # Remove a value from a frontmatter array
  obsidian-cli note patch "Document.md" --target "tags" --array-remove "old-tag"

  # Append to a block reference
  obsidian-cli note patch "Notes.md" --target "abc123" --target-type block --operation append --content "Additional context"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		content := patchContent

		apiClient, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		// Handle array operations by updating the full note
		if patchArrayAdd != "" || patchArrayRemove != "" {
			if patchTarget == "" {
				return fmt.Errorf("target is required for array operations (use --target)")
			}

			// Get the current note
			note, err := apiClient.GetNote(path)
			if err != nil {
				return fmt.Errorf("failed to get note: %w", err)
			}

			// Get the current value of the field from frontmatter
			currentValue, exists := note.Frontmatter[patchTarget]
			var currentArray []string

			if exists {
				// Try to convert to string array
				switch v := currentValue.(type) {
				case []interface{}:
					for _, item := range v {
						if str, ok := item.(string); ok {
							// Strip leading # if present (Obsidian API adds these)
							str = strings.TrimPrefix(str, "#")
							currentArray = append(currentArray, str)
						}
					}
				case []string:
					for _, str := range v {
						// Strip leading # if present
						str = strings.TrimPrefix(str, "#")
						currentArray = append(currentArray, str)
					}
				case string:
					// Single value, convert to array
					v = strings.TrimPrefix(v, "#")
					currentArray = []string{v}
				default:
					return fmt.Errorf("field %s is not an array or string", patchTarget)
				}
			}

			// Perform array operation
			if patchArrayAdd != "" {
				// Check if value already exists
				found := false
				for _, item := range currentArray {
					if item == patchArrayAdd {
						found = true
						break
					}
				}
				if !found {
					currentArray = append(currentArray, patchArrayAdd)
				}
			} else if patchArrayRemove != "" {
				// Remove value from array
				newArray := []string{}
				for _, item := range currentArray {
					if item != patchArrayRemove {
						newArray = append(newArray, item)
					}
				}
				currentArray = newArray
			}

			// Parse the note content to update frontmatter
			noteContent := note.Content

			// Find and replace the frontmatter field
			frontmatterRegex := regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
			matches := frontmatterRegex.FindStringSubmatch(noteContent)

			if len(matches) < 2 {
				return fmt.Errorf("could not parse frontmatter")
			}

			frontmatter := matches[1]

			// Build the new array value
			var newValue string
			if len(currentArray) == 0 {
				newValue = "[]"
			} else {
				lines := []string{}
				for _, item := range currentArray {
					lines = append(lines, "  - "+item)
				}
				newValue = "\n" + strings.Join(lines, "\n")
			}

			// Replace the field in frontmatter
			fieldRegex := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(patchTarget) + `:.*?(?:\n(?:  |\t).*)*`)
			newFrontmatter := fieldRegex.ReplaceAllString(frontmatter, patchTarget+":"+newValue)

			// If field didn't exist, add it
			if !fieldRegex.MatchString(frontmatter) {
				newFrontmatter = frontmatter + "\n" + patchTarget + ":" + newValue
			}

			// Reconstruct the note
			newContent := "---\n" + newFrontmatter + "\n---" + strings.TrimPrefix(noteContent, matches[0])

			// Update the note
			_, err = apiClient.UpdateNote(path, newContent)
			if err != nil {
				return fmt.Errorf("failed to update note: %w", err)
			}

			formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
			if patchArrayAdd != "" {
				formatter.PrintSuccess(fmt.Sprintf("Added '%s' to %s in: %s", patchArrayAdd, patchTarget, path))
			} else {
				formatter.PrintSuccess(fmt.Sprintf("Removed '%s' from %s in: %s", patchArrayRemove, patchTarget, path))
			}
			return nil
		}

		// Read from stdin if requested
		if patchFromStdin {
			reader := bufio.NewReader(os.Stdin)
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(data)
		}

		if content == "" && !patchFromStdin && patchArrayAdd == "" && patchArrayRemove == "" {
			return fmt.Errorf("content is required (use --content, --stdin, --array-add, or --array-remove)")
		}

		if patchOperation == "" {
			return fmt.Errorf("operation is required (use --operation with append, prepend, or replace)")
		}

		if patchTargetType == "" {
			return fmt.Errorf("target-type is required (use --target-type with heading, block, or frontmatter)")
		}

		if patchTarget == "" {
			return fmt.Errorf("target is required (use --target)")
		}

		// Validate operation
		var operation client.PatchOperation
		switch patchOperation {
		case "append":
			operation = client.PatchOperationAppend
		case "prepend":
			operation = client.PatchOperationPrepend
		case "replace":
			operation = client.PatchOperationReplace
		default:
			return fmt.Errorf("invalid operation: %s (must be append, prepend, or replace)", patchOperation)
		}

		// Validate target type
		var targetType client.PatchTargetType
		switch patchTargetType {
		case "heading":
			targetType = client.PatchTargetHeading
		case "block":
			targetType = client.PatchTargetBlock
		case "frontmatter":
			targetType = client.PatchTargetFrontmatter
		default:
			return fmt.Errorf("invalid target-type: %s (must be heading, block, or frontmatter)", patchTargetType)
		}

		req := &client.PatchNoteRequest{
			Path:                  path,
			Content:               content,
			Operation:             operation,
			TargetType:            targetType,
			Target:                patchTarget,
			TargetDelimiter:       patchTargetDelimiter,
			TrimTargetWhitespace:  patchTrimWhitespace,
			CreateTargetIfMissing: patchCreateIfMissing,
		}

		err = apiClient.PatchNote(req)
		if err != nil {
			return fmt.Errorf("failed to patch note: %w", err)
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Patched note: %s (operation: %s, target: %s %s)", path, patchOperation, patchTargetType, patchTarget))
		return nil
	},
}

func init() {
	patchCmd.Flags().StringVarP(&patchContent, "content", "c", "", "content to append/prepend/replace")
	patchCmd.Flags().BoolVar(&patchFromStdin, "stdin", false, "read content from stdin")
	patchCmd.Flags().StringVar(&patchOperation, "operation", "", "operation: append, prepend, or replace (required unless using array operations)")
	patchCmd.Flags().StringVarP(&patchTargetType, "target-type", "t", "", "target type: heading, block, or frontmatter (required unless using array operations)")
	patchCmd.Flags().StringVar(&patchTarget, "target", "", "target heading/block/field name (required)")
	patchCmd.Flags().StringVar(&patchTargetDelimiter, "delimiter", "", "delimiter for nested headings (default: ::)")
	patchCmd.Flags().BoolVar(&patchTrimWhitespace, "trim-whitespace", false, "trim whitespace from target")
	patchCmd.Flags().BoolVar(&patchCreateIfMissing, "create-if-missing", false, "create frontmatter field if it doesn't exist")
	patchCmd.Flags().StringVar(&patchArrayAdd, "array-add", "", "add a value to a frontmatter array field")
	patchCmd.Flags().StringVar(&patchArrayRemove, "array-remove", "", "remove a value from a frontmatter array field")
}
