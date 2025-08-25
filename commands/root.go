package commands

import (
	"github.com/spf13/cobra"

	"clusterix-code/internal/commands"
)

var RootCmd = &cobra.Command{
	Use:   "clx-code",
	Short: "Clusterix Code API CLI",
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.AddCommand(commands.MigrateUpCmd)
	RootCmd.AddCommand(commands.MigrateDownCmd)
	RootCmd.AddCommand(commands.ImportMachineConfigsCmd)
	RootCmd.AddCommand(commands.SyncUsersCmd)
}
