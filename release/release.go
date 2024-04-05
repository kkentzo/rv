package release

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const (
	ReleaseFormat   = "20060102150405.000"
	CurrentLinkName = "current"
)

var ReleaseFormatRe = regexp.MustCompile(`\b\d{14}\.\d{3}\b`)

// Execute the release flow given a workspace directory and a zip file (bundle)
// If the username is empty, then the current user/group is used
// Steps:
// 1. create the workspace if necessary
// 2. resolve the uid and gid of the files to be created
// 3. create the release directory inside the workspace
// 4. decompress the bundle into the release directory
// 5. update the workspace's `current` link to point to the new release
// 6. apply the policy of how many releases to keep
//
// The function returns the ID of the release (directory name) and/or an error
// if the ID is not an empty string, then the release directory still exists (even on error) and can be used
func Install(workspaceDir, bundlePath string, keepN uint, username, groupname string, stdout io.Writer) (string, error) {
	// we should not accept this value because
	// it will leave us with no releases at all
	if keepN == 0 {
		return "", errors.New("can not accept keeping no releases in the workspace")
	}
	// we will work with absolute directories
	if !path.IsAbs(workspaceDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to determine the current working directory: %v", err)
		}
		workspaceDir = path.Join(cwd, workspaceDir)

	}
	fmt.Fprintf(stdout, "[info] workspace=%s\n", workspaceDir)

	// figure out file/directory ownership
	uid, gid, err := resolveUser(username)
	if err != nil {
		return "", fmt.Errorf("failed to resolve user: %v", err)
	}
	if groupname != "" {
		gid, err = resolveGroup(groupname)
		if err != nil {
			return "", fmt.Errorf("failed to resolve group: %v", err)
		}
	}

	// create release under workspace
	id := time.Now().Format(ReleaseFormat)
	releaseDir := path.Join(workspaceDir, id)
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create release: %v", err)
	}
	fmt.Fprintf(stdout, "[info] release=%s\n", id)
	// decompress bundle file
	fmt.Fprintf(stdout, "[release] unpacking bundle=%s to %s\n", bundlePath, releaseDir)
	if err := decompressArchive(bundlePath, releaseDir, uid, gid); err != nil {
		// cleanup release directory
		defer os.RemoveAll(releaseDir)
		return "", fmt.Errorf("failed to decompress archive: %v", err)
	}

	// update current link
	fmt.Fprintf(stdout, "[release] updating current to %s\n", id)
	if err := createOrUpdateLink(workspaceDir, id); err != nil {
		// cleanup release directory
		defer os.RemoveAll(releaseDir)
		return "", fmt.Errorf("failed to create/update link: %v", err)
	}
	// clean up excess releases
	if err := cleanupReleases(workspaceDir, keepN, stdout); err != nil {
		return id, fmt.Errorf("failed to clean up releases (keep=%d)", keepN)
	}
	return id, nil
}

func Rewind(workspaceDir, target string, stdout io.Writer) (string, error) {
	releases, err := getReleasesDesc(workspaceDir)
	if err != nil {
		return target, err
	}
	// do we need to figure out which release to rewind to?
	if target == "" {
		switch len(releases) {
		case 0:
			return "", errors.New("can not rewind: no releases found in workspace")
		case 1:
			return "", errors.New("can not rewind having only one release in workspace")
		default:
			target = releases[1]
		}
	}

	targetPath := path.Join(workspaceDir, target)
	// does the target exist? => noop
	if !fileExists(targetPath) {
		return "", fmt.Errorf("release %s not found", target)
	}
	// figure out current link
	current, err := GetCurrent(workspaceDir)
	if err != nil {
		return "", fmt.Errorf("could not determine current release: %v", err)
	}
	fmt.Fprintf(stdout, "[info] current=%s\n", current)

	// is the target the same as the current link? => noop
	if current == target {
		return "", fmt.Errorf("will not rewind: target %s is already current", target)
	}

	// set the current link to the target release
	fmt.Fprintf(stdout, "[rewind] setting current to %s\n", target)
	if err := createOrUpdateLink(workspaceDir, target); err != nil {
		return "", fmt.Errorf("current link: %v", err)
	}

	// delete the releases that were performed later than the target release
	for _, rel := range releases {
		if rel == target {
			break
		}
		releasePath := path.Join(workspaceDir, rel)
		fmt.Fprintf(stdout, "[cleanup] deleting %s\n", rel)
		if err := os.RemoveAll(releasePath); err != nil {
			return target, fmt.Errorf("failed to delete release %s: %v", releasePath, err)
		}
	}

	return target, nil
}

func List(workspaceDir string) ([]string, error) {
	releases, err := getReleasesDesc(workspaceDir)
	if err != nil {
		return releases, fmt.Errorf("failed to list releases: %v", err)
	}
	current, err := os.Readlink(path.Join(workspaceDir, CurrentLinkName))
	if err != nil {
		return releases, fmt.Errorf("failed to resolve current release: %v", err)
	}
	// mark current release
	for idx, rel := range releases {
		if rel == current {
			releases[idx] = rel + " <== current"
		}
	}

	return releases, nil
}

func cleanupReleases(workspaceDir string, keepN uint, stdout io.Writer) error {
	releases, err := getReleasesAsc(workspaceDir)
	if err != nil {
		return err
	}
	// assemble all release file names (and only those)
	obsoleteN := len(releases) - int(keepN)
	if obsoleteN > 0 {
		for idx, releaseName := range releases {
			if idx < obsoleteN {
				releasePath := path.Join(workspaceDir, releaseName)
				fmt.Fprintf(stdout, "[cleanup] deleting %s (keep=%d)\n", releaseName, keepN)
				if err := os.RemoveAll(releasePath); err != nil {
					return fmt.Errorf("failed to delete release %s: %v", releasePath, err)
				}
			}
		}
	}

	return nil
}

func getReleasesAsc(workspaceDir string) ([]string, error) {
	releases, err := getReleases(workspaceDir)
	if err != nil {
		return []string{}, err
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i] < releases[j]
	})
	return releases, nil
}

func getReleasesDesc(workspaceDir string) ([]string, error) {
	releases, err := getReleases(workspaceDir)
	if err != nil {
		return []string{}, err
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i] > releases[j]
	})
	return releases, nil
}

func getReleases(workspaceDir string) ([]string, error) {
	releases := []string{}

	entries, err := ioutil.ReadDir(workspaceDir)
	if err != nil {
		return releases, err
	}

	for _, e := range entries {
		if ReleaseFormatRe.Match([]byte(e.Name())) {
			releases = append(releases, e.Name())
		}
	}

	return releases, nil
}

func createOrUpdateLink(workspaceDir, target string) error {
	link := path.Join(workspaceDir, CurrentLinkName)
	// does the link already exist?
	_, err := os.Stat(link)
	if !os.IsNotExist(err) {
		if err := os.Remove(link); err != nil {
			return fmt.Errorf("update failed: %v", err)
		}
	}
	// TODO: what if the following fails? We are stuck with no `current` link
	if err := os.Symlink(target, link); err != nil {
		return fmt.Errorf("create failed: %v", err)
	}
	return nil
}

// return the uid and gid of the requested user
// if the username is empty,
// then the function returns the uid and gid of the current user
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

// return the gid of the requested group
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

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

func GetCurrent(workspace string) (string, error) {
	target, err := os.Readlink(path.Join(workspace, CurrentLinkName))
	if err != nil {
		return "", err
	}
	return target, nil
}
