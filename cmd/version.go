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

var (
	// command
	versionCmdDescription = "list all the releases in the workspace"
	VersionCmd            = &cobra.Command{
		Use:   "version",
		Short: versionCmdDescription,
		Long:  versionCmdDescription,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "%s [%s]\n", AppVersion, GitCommit)
		},
	}
)

func init() {
	RootCmd.AddCommand(VersionCmd)
}
