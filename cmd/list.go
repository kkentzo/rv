package cmd

import (
	"fmt"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

func ListCommand(globals *GlobalVariables) *cobra.Command {
	descr := "list all the releases in the workspace"
	cmd := &cobra.Command{
		Use:   "list",
		Short: descr,
		Long:  descr,
		Run: func(cmd *cobra.Command, args []string) {
			if releases, err := release.List(globals.WorkspacePath); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				for _, rel := range releases {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", rel)
				}
			}
		},
	}

	return requireGlobalFlags(cmd, globals)
}
