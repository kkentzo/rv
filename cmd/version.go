package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// these will be populated at build time using an ldflag
var (
	GitCommit  string
	AppVersion string
)

func VersionCommand() *cobra.Command {
	// command
	descr := "list all the releases in the workspace"
	return &cobra.Command{
		Use:   "version",
		Short: descr,
		Long:  descr,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]\n", AppVersion, GitCommit)
		},
	}
}
