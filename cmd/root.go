package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// persistent (global) command-line arguments
	// each command that needs it, sets the flag and uses this variable
	globalWorkspacePath string
	// command
	rootCmdDescription = "Manage multiple release bundles locally"
	RootCmd            = &cobra.Command{
		Use:   "rv",
		Short: rootCmdDescription,
		Long:  rootCmdDescription,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func requireWorkspaceFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&globalWorkspacePath, "workspace", "w", "", "directory that contains all available releases")
	cmd.MarkFlagRequired("workspace")
}
