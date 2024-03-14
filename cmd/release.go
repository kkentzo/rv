package cmd

import (
	"errors"
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var (
	// command-line arguments
	archivePath string
	keepN       uint
	// command
	releaseCmdDescription = "Uncompress the specified archive into the workspace and update the `current` link"
	ReleaseCmd            = &cobra.Command{
		Use:   "release",
		Short: releaseCmdDescription,
		Long:  releaseCmdDescription,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if keepN == 0 {
				return errors.New("zero is not a valid value for --keep (-k) flag")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if releaseID, err := release.Install(globalWorkspacePath, archivePath, keepN, cmd.OutOrStdout()); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "[success] active version is %s\n", releaseID)
			}
		},
	}
)

func init() {
	requireWorkspaceFlag(ReleaseCmd)

	ReleaseCmd.Flags().StringVarP(&archivePath, "archive", "a", "", "path to archive file containing the release")
	ReleaseCmd.Flags().UintVarP(&keepN, "keep", "k", 3, "maximum number of releases to keep in workspace at all times")
	ReleaseCmd.MarkFlagRequired("archive")

	RootCmd.AddCommand(ReleaseCmd)
}
