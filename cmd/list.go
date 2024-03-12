package cmd

import (
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

var (
	// command
	listCmdDescription = "list all the releases in the workspace"
	ListCmd            = &cobra.Command{
		Use:   "list",
		Short: listCmdDescription,
		Long:  listCmdDescription,
		Run: func(cmd *cobra.Command, args []string) {
			if releases, err := release.List(globalWorkspacePath); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				for _, rel := range releases {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", rel)
				}
			}
		},
	}
)

func init() {
	requireWorkspaceFlag(ListCmd)
	RootCmd.AddCommand(ListCmd)
}
