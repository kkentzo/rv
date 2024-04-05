package cmd

import (
	"github.com/spf13/cobra"
)

type GlobalVariables struct {
	WorkspacePath string
}

func New() *cobra.Command {
	globals := &GlobalVariables{}

	descr := "Manage multiple release bundles locally"

	root := &cobra.Command{
		Use:   "rv",
		Short: descr,
		Long:  descr,
	}
	root.AddCommand(ReleaseCommand(globals))
	root.AddCommand(ListCommand(globals))
	root.AddCommand(RewindCommand(globals))
	root.AddCommand(VersionCommand())
	return root
}

func requireGlobalFlags(cmd *cobra.Command, globals *GlobalVariables) *cobra.Command {
	cmd.Flags().StringVarP(&globals.WorkspacePath, "workspace", "w", "", "directory that contains all available releases")
	cmd.MarkFlagRequired("workspace")
	return cmd
}
