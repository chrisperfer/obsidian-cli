package note

import (
	"github.com/spf13/cobra"
)

// NoteCmd represents the note command
var NoteCmd = &cobra.Command{
	Use:   "note",
	Short: "Create, read, update, and delete notes",
	Long:  `Perform operations on individual notes in your Obsidian vault`,
}

func init() {
	NoteCmd.AddCommand(getCmd)
	NoteCmd.AddCommand(createCmd)
	NoteCmd.AddCommand(updateCmd)
	NoteCmd.AddCommand(patchCmd)
	NoteCmd.AddCommand(appendCmd)
	NoteCmd.AddCommand(deleteCmd)
	NoteCmd.AddCommand(searchCmd)
}
