package cmd

import (
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var (
	// command-line arguments
	keepN uint
	// command
	releaseCmdDescription = "Uncompress the specified bundle into the workspace and update the `current` link"
	ReleaseCmd            = &cobra.Command{
		Use:   "release",
		Short: releaseCmdDescription,
		Long:  releaseCmdDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if releaseID, err := release.Install(workspacePath, args[0], keepN); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), releaseID)
			}
		},
	}
)

func init() {
	ReleaseCmd.Flags().UintVarP(&keepN, "keep", "k", 3, "maximum number of releases to keep in workspace at all times")
	RootCmd.AddCommand(ReleaseCmd)
}
