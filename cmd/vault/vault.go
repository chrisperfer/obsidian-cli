package vault

import (
	"github.com/spf13/cobra"
)

// VaultCmd represents the vault command
var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage vault and list notes",
	Long:  `Perform vault-level operations such as listing notes and viewing vault information`,
}

func init() {
	VaultCmd.AddCommand(listCmd)
	VaultCmd.AddCommand(infoCmd)
}
