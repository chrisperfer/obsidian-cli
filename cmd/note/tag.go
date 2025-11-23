package note

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chrisperfer/obsidian-cli/pkg/client"
	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage note tags",
	Long:  `Add, remove, or list tags on notes. Convenience wrapper around patch --array-add/--array-remove for the tags field.`,
}

var tagAddCmd = &cobra.Command{
	Use:   "add [path] [tags...]",
	Short: "Add tags to a note",
	Long:  `Add one or more tags to a note's frontmatter`,
	Args:  cobra.MinimumNArgs(2),
	Example: `  # Add a single tag
  obsidian-cli note tag add "Document.md" "important"

  # Add multiple tags
  obsidian-cli note tag add "Document.md" "important" "urgent" "review"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		tags := args[1:]

		apiClient, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		// Get the current note
		note, err := apiClient.GetNote(path)
		if err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		// Get current tags
		currentTags := []string{}
		if tagsValue, exists := note.Frontmatter["tags"]; exists {
			switch v := tagsValue.(type) {
			case []interface{}:
				for _, item := range v {
					if str, ok := item.(string); ok {
						str = strings.TrimPrefix(str, "#")
						currentTags = append(currentTags, str)
					}
				}
			case []string:
				for _, str := range v {
					str = strings.TrimPrefix(str, "#")
					currentTags = append(currentTags, str)
				}
			case string:
				v = strings.TrimPrefix(v, "#")
				currentTags = []string{v}
			}
		}

		// Add new tags (avoid duplicates)
		addedTags := []string{}
		for _, tag := range tags {
			found := false
			for _, existing := range currentTags {
				if existing == tag {
					found = true
					break
				}
			}
			if !found {
				currentTags = append(currentTags, tag)
				addedTags = append(addedTags, tag)
			}
		}

		if len(addedTags) == 0 {
			formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
			formatter.PrintSuccess("No new tags to add (all tags already exist)")
			return nil
		}

		// Update the note using the same frontmatter update logic from patch
		err = updateNoteTags(apiClient, path, note.Content, currentTags)
		if err != nil {
			return err
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Added tags: %s", strings.Join(addedTags, ", ")))
		return nil
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove [path] [tags...]",
	Short: "Remove tags from a note",
	Long:  `Remove one or more tags from a note's frontmatter`,
	Args:  cobra.MinimumNArgs(2),
	Example: `  # Remove a single tag
  obsidian-cli note tag remove "Document.md" "old-tag"

  # Remove multiple tags
  obsidian-cli note tag remove "Document.md" "old-tag" "deprecated" "archive"`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]
		tagsToRemove := args[1:]

		apiClient, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		// Get the current note
		note, err := apiClient.GetNote(path)
		if err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		// Get current tags
		currentTags := []string{}
		if tagsValue, exists := note.Frontmatter["tags"]; exists {
			switch v := tagsValue.(type) {
			case []interface{}:
				for _, item := range v {
					if str, ok := item.(string); ok {
						str = strings.TrimPrefix(str, "#")
						currentTags = append(currentTags, str)
					}
				}
			case []string:
				for _, str := range v {
					str = strings.TrimPrefix(str, "#")
					currentTags = append(currentTags, str)
				}
			case string:
				v = strings.TrimPrefix(v, "#")
				currentTags = []string{v}
			}
		}

		// Remove tags
		newTags := []string{}
		removedTags := []string{}
		for _, tag := range currentTags {
			shouldRemove := false
			for _, removeTag := range tagsToRemove {
				if tag == removeTag {
					shouldRemove = true
					removedTags = append(removedTags, tag)
					break
				}
			}
			if !shouldRemove {
				newTags = append(newTags, tag)
			}
		}

		if len(removedTags) == 0 {
			formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
			formatter.PrintSuccess("No tags removed (tags not found)")
			return nil
		}

		// Update the note
		err = updateNoteTags(apiClient, path, note.Content, newTags)
		if err != nil {
			return err
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		formatter.PrintSuccess(fmt.Sprintf("Removed tags: %s", strings.Join(removedTags, ", ")))
		return nil
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List tags on a note",
	Long:  `Display all tags from a note's frontmatter`,
	Args:  cobra.ExactArgs(1),
	Example: `  # List tags
  obsidian-cli note tag list "Document.md"

  # List tags in JSON format
  obsidian-cli note tag list "Document.md" --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		path := args[0]

		apiClient, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		// Get the note
		note, err := apiClient.GetNote(path)
		if err != nil {
			return fmt.Errorf("failed to get note: %w", err)
		}

		// Get tags
		tags := []string{}
		if tagsValue, exists := note.Frontmatter["tags"]; exists {
			switch v := tagsValue.(type) {
			case []interface{}:
				for _, item := range v {
					if str, ok := item.(string); ok {
						str = strings.TrimPrefix(str, "#")
						tags = append(tags, str)
					}
				}
			case []string:
				for _, str := range v {
					str = strings.TrimPrefix(str, "#")
					tags = append(tags, str)
				}
			case string:
				v = strings.TrimPrefix(v, "#")
				tags = []string{v}
			}
		}

		formatter := output.NewFormatter(cmdutil.GetOutputFormat(cobraCmd), cmdutil.GetLLMMode(cobraCmd))
		return formatter.Print(tags)
	},
}

func init() {
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	tagCmd.AddCommand(tagListCmd)
}

// updateNoteTags updates the tags field in a note's frontmatter
func updateNoteTags(apiClient *client.Client, path, noteContent string, tags []string) error {
	// Find and replace the frontmatter field
	frontmatterRegex := regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
	matches := frontmatterRegex.FindStringSubmatch(noteContent)

	if len(matches) < 2 {
		return fmt.Errorf("could not parse frontmatter")
	}

	frontmatter := matches[1]

	// Build the new array value
	var newValue string
	if len(tags) == 0 {
		newValue = " []"
	} else {
		lines := []string{}
		for _, tag := range tags {
			lines = append(lines, "  - " + tag)
		}
		newValue = "\n" + strings.Join(lines, "\n")
	}

	// Replace the field in frontmatter
	fieldRegex := regexp.MustCompile(`(?m)^tags:.*?(?:\n(?:  |\t).*)*`)
	newFrontmatter := fieldRegex.ReplaceAllString(frontmatter, "tags:"+newValue)

	// If field didn't exist, add it
	if !fieldRegex.MatchString(frontmatter) {
		newFrontmatter = frontmatter + "\ntags:" + newValue
	}

	// Reconstruct the note
	newContent := "---\n" + newFrontmatter + "\n---" + strings.TrimPrefix(noteContent, matches[0])

	// Update the note
	_, err := apiClient.UpdateNote(path, newContent)
	return err
}
