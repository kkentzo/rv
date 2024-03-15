package cmd

import (
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

// rewindCmd represents the rewind command
var (
	// command-line arguments
	target string
	// command
	rewindCmdDescription = "Reset the current release"
	RewindCmd            = &cobra.Command{
		Use:   "rewind",
		Short: rewindCmdDescription,
		Long:  rewindCmdDescription,
		Run: func(cmd *cobra.Command, args []string) {
			// perform release
			releaseID, err := release.Rewind(globalWorkspacePath, target, cmd.OutOrStdout())
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "[success] active version is %s\n", releaseID)
			}
		},
	}
)

func init() {
	requireWorkspaceFlag(RewindCmd)
	RewindCmd.Flags().StringVarP(&target, "target", "t", "", "target release to reset the current link to")
	RootCmd.AddCommand(RewindCmd)
}
