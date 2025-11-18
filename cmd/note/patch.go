package note

import (
	"bufio"
	"fmt"
	"io"
	"os"

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

  # Append to a block reference
  obsidian-cli note patch "Notes.md" --target "abc123" --target-type block --operation append --content "Additional context"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		content := patchContent

		// Read from stdin if requested
		if patchFromStdin {
			reader := bufio.NewReader(os.Stdin)
			data, err := io.ReadAll(reader)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			content = string(data)
		}

		if content == "" && !patchFromStdin {
			return fmt.Errorf("content is required (use --content or --stdin)")
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

		apiClient, err := cmdutil.GetClient()
		if err != nil {
			return err
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
	patchCmd.Flags().StringVar(&patchOperation, "operation", "", "operation: append, prepend, or replace (required)")
	patchCmd.Flags().StringVarP(&patchTargetType, "target-type", "t", "", "target type: heading, block, or frontmatter (required)")
	patchCmd.Flags().StringVar(&patchTarget, "target", "", "target heading/block/field name (required)")
	patchCmd.Flags().StringVar(&patchTargetDelimiter, "delimiter", "", "delimiter for nested headings (default: ::)")
	patchCmd.Flags().BoolVar(&patchTrimWhitespace, "trim-whitespace", false, "trim whitespace from target")
	patchCmd.Flags().BoolVar(&patchCreateIfMissing, "create-if-missing", false, "create frontmatter field if it doesn't exist")
}
