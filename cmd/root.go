package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// persistent (global) command-line arguments
	rootDirectory string
	// command
	rootCmdDescription = "Manage local releases"
	rootCmd            = &cobra.Command{
		Use:   "rv",
		Short: rootCmdDescription,
		Long:  rootCmdDescription,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().
		StringVarP(&rootDirectory, "workspace", "w", "", "directory under which releases will be located and managed")
	rootCmd.MarkPersistentFlagRequired("workspace")
}
