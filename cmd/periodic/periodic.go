package periodic

import (
	"github.com/spf13/cobra"
)

// PeriodicCmd represents the periodic command
var PeriodicCmd = &cobra.Command{
	Use:   "periodic",
	Short: "Work with periodic notes",
	Long:  `Create, read, and list periodic notes (daily, weekly, monthly)`,
}

func init() {
	PeriodicCmd.AddCommand(listCmd)
	PeriodicCmd.AddCommand(getCmd)
	PeriodicCmd.AddCommand(createCmd)
}
