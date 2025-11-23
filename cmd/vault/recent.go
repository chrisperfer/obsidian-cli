package vault

import (
	"fmt"
	"sort"
	"time"

	"github.com/chrisperfer/obsidian-cli/pkg/cmdutil"
	"github.com/chrisperfer/obsidian-cli/pkg/output"
	"github.com/spf13/cobra"
)

var (
	recentDays  int
	recentHours int
	recentSort  string
)

type FileWithStat struct {
	Path     string    `json:"path"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Size     int64     `json:"size"`
}

var recentCmd = &cobra.Command{
	Use:   "recent [path]",
	Short: "Find recently created or modified notes",
	Long:  `Find notes that were created or modified within a specified time period. Defaults to last 3 days.`,
	Example: `  # Notes modified in last 3 days (default)
  obsidian-cli vault recent

  # Custom time periods
  obsidian-cli vault recent --days 7
  obsidian-cli vault recent --hours 24

  # Start from specific path
  obsidian-cli vault recent "Daily/"
  obsidian-cli vault recent "Food and Drink/" --days 30

  # Sort by creation time instead of modification time
  obsidian-cli vault recent --sort created

  # Output in JSON format
  obsidian-cli vault recent --output json`,
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		client, err := cmdutil.GetClient()
		if err != nil {
			return err
		}

		var basePath string
		if len(args) > 0 {
			basePath = args[0]
		}

		// Calculate cutoff time
		var cutoff time.Time
		if recentHours > 0 {
			cutoff = time.Now().Add(-time.Duration(recentHours) * time.Hour)
		} else {
			cutoff = time.Now().AddDate(0, 0, -recentDays)
		}

		// Get all files
		allFiles, err := client.ListAllFiles()
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		// Filter by base path
		if basePath != "" {
			filtered := []string{}
			for _, file := range allFiles {
				if len(file) >= len(basePath) && file[:len(basePath)] == basePath {
					filtered = append(filtered, file)
				}
			}
			allFiles = filtered
		}

		// Get metadata for each file and filter by time
		recentFiles := []FileWithStat{}
		for _, file := range allFiles {
			note, err := client.GetNote(file)
			if err != nil {
				// Skip files we can't read
				continue
			}

			if note.Stat == nil {
				continue
			}

			created := time.Unix(note.Stat.CTime/1000, 0)
			modified := time.Unix(note.Stat.MTime/1000, 0)

			// Check if file meets criteria
			isRecent := false
			if recentSort == "created" {
				isRecent = created.After(cutoff)
			} else {
				isRecent = modified.After(cutoff)
			}

			if isRecent {
				recentFiles = append(recentFiles, FileWithStat{
					Path:     file,
					Created:  created,
					Modified: modified,
					Size:     note.Stat.Size,
				})
			}
		}

		// Sort results
		sort.Slice(recentFiles, func(i, j int) bool {
			if recentSort == "created" {
				return recentFiles[i].Created.After(recentFiles[j].Created)
			}
			return recentFiles[i].Modified.After(recentFiles[j].Modified)
		})

		// Format output based on format flag
		outputFormat := cmdutil.GetOutputFormat(cobraCmd)
		if outputFormat == "json" || outputFormat == "yaml" {
			formatter := output.NewFormatter(outputFormat, cmdutil.GetLLMMode(cobraCmd))
			return formatter.Print(recentFiles)
		}

		// Text format - show file paths with timestamps
		if len(recentFiles) == 0 {
			fmt.Println("No recent files found")
			return nil
		}

		for _, f := range recentFiles {
			if recentSort == "created" {
				fmt.Printf("%s  (created: %s)\n", f.Path, f.Created.Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("%s  (modified: %s)\n", f.Path, f.Modified.Format("2006-01-02 15:04:05"))
			}
		}

		return nil
	},
}

func init() {
	recentCmd.Flags().IntVarP(&recentDays, "days", "d", 3, "Number of days to look back")
	recentCmd.Flags().IntVar(&recentHours, "hours", 0, "Number of hours to look back (overrides --days)")
	recentCmd.Flags().StringVarP(&recentSort, "sort", "s", "modified", "Sort by 'created' or 'modified' time")
}
