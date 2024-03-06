package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var (
	releaseCmdDescription = "Uncompress the specified bundle into the release folder and update the `current` link"
	releaseCmd            = &cobra.Command{
		Use:   "release",
		Short: releaseCmdDescription,
		Long:  releaseCmdDescription,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			bundlePath := args[0]
			fmt.Printf("release command called with bundle_path=%s [workspace=%s]", bundlePath, rootDirectory)
		},
	}
)

func init() {
	rootCmd.AddCommand(releaseCmd)
}
