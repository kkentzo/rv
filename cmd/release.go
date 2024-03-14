package cmd

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"

	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var (
	// command-line arguments
	archivePath         string
	keepN               uint
	username, groupname string
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
			// figure out file ownership
			uid, gid, err := resolveUser(username)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
				return
			}
			if groupname != "" {
				gid, err = resolveGroup(groupname)
				if err != nil {
					fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
					return
				}
			}
			// perform release
			if releaseID, err := release.Install(globalWorkspacePath, archivePath, keepN, uid, gid, cmd.OutOrStdout()); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "error: %v\n", err)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "[success] active version is %s\n", releaseID)
			}
		},
	}
)

func resolveUser(username string) (uid int, gid int, err error) {
	var u *user.User
	if username == "" {
		u, err = user.Current()
	} else {
		u, err = user.Lookup(username)
	}
	if err != nil {
		return
	}
	uid, err = strconv.Atoi(u.Uid)
	if err != nil {
		return
	}
	gid, err = strconv.Atoi(u.Gid)
	if err != nil {
		return
	}
	return
}

func resolveGroup(groupname string) (gid int, err error) {
	var g *user.Group
	g, err = user.LookupGroup(groupname)
	if err != nil {
		return
	}
	gid, err = strconv.Atoi(g.Gid)
	if err != nil {
		return
	}
	return
}

func init() {
	requireWorkspaceFlag(ReleaseCmd)

	ReleaseCmd.Flags().StringVarP(&archivePath, "archive", "a", "", "path to archive file containing the release")
	ReleaseCmd.Flags().UintVarP(&keepN, "keep", "k", 3, "maximum number of releases to keep in workspace at all times")
	ReleaseCmd.Flags().StringVarP(&username, "user", "u", "", "user to whom all extracted archive files will belong to")
	ReleaseCmd.Flags().StringVarP(&groupname, "group", "g", "", "group to whom all extracted archive files will belong to")
	ReleaseCmd.MarkFlagRequired("archive")

	RootCmd.AddCommand(ReleaseCmd)
}
