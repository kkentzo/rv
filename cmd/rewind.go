package cmd

import (
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

func RewindCommand(globals *GlobalVariables) *cobra.Command {
	var (
		// command-line arguments
		target string
		// command
		descr = "Reset the current release"
		cmd   = &cobra.Command{
			Use:   "rewind",
			Short: descr,
			Long:  descr,
			Run: func(cmd *cobra.Command, args []string) {
				// perform release
				releaseID, err := release.Rewind(globals.WorkspacePath, target, cmd.OutOrStdout())
				if err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "[success] active version is %s\n", releaseID)
				}
			},
		}
	)

	cmd.Flags().StringVarP(&target, "target", "t", "", "target release to reset the current link to")
	return requireGlobalFlags(cmd, globals)
}
