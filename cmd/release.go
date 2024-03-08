package cmd

import (
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var (
	releaseCmdDescription = "Uncompress the specified bundle into the workspace and update the `current` link"
	releaseCmd            = &cobra.Command{
		Use:   "release",
		Short: releaseCmdDescription,
		Long:  releaseCmdDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if releaseID, err := release.Install(workspacePath, args[0]); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), releaseID)
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(releaseCmd)
}
