package release

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

const (
	ReleaseFormat   = "20060102150405.000"
	CurrentLinkName = "current"
)

var ReleaseFormatRe = regexp.MustCompile(`\b\d{14}\.\d{3}\b`)

// Execute the release flow given a workspace directory and a zip file (bundle):
// 1. creates the workspace if necessary
// 2. creates the release directory inside the workspace
// 3. decompresses the bundle into the release directory
// 4. updates the workspace's `current` link to point to the new release
// 5. applies the policy of how many releases to keep
// The function returns the ID of the release (directory name) and/or an error
// if the ID is not an empty string, then the release directory still exists (even on error) and can be used
func Install(workspaceDir, bundlePath string, keepN uint) (string, error) {
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
	// create release under workspace
	id := time.Now().Format(ReleaseFormat)
	releaseDir := path.Join(workspaceDir, id)
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create release: %v", err)
	}
	// decompress bundle file
	if err := decompressZip(bundlePath, releaseDir); err != nil {
		// cleanup release directory
		defer os.RemoveAll(releaseDir)
		return "", fmt.Errorf("failed to decompress archive: %v", err)
	}
	// update current link
	if err := createOrUpdateLink(workspaceDir, id); err != nil {
		// cleanup release directory
		defer os.RemoveAll(releaseDir)
		return "", fmt.Errorf("failed to create/update link: %v", err)
	}
	// clean up excess releases
	if err := cleanupReleases(workspaceDir, keepN); err != nil {
		return id, fmt.Errorf("failed to clean up releases (keep=%d)", keepN)
	}
	return id, nil
}

func List(workspaceDir string) ([]string, error) {
	releases, err := getReleases(workspaceDir)
	if err != nil {
		return releases, fmt.Errorf("failed to list releases: %v", err)
	}
	// sort releases in descending order (oldest to newest)
	sort.Slice(releases, func(i, j int) bool {
		return releases[i] > releases[j]
	})

	return releases, nil
}

func cleanupReleases(workspaceDir string, keepN uint) error {
	releases, err := getReleases(workspaceDir)
	if err != nil {
		return err
	}
	// assemble all release file names (and only those)
	obsoleteN := len(releases) - int(keepN)
	if obsoleteN > 0 {
		// sort releases in ascending order (oldest to newest)
		sort.Slice(releases, func(i, j int) bool {
			return releases[i] < releases[j]
		})
		for idx, releaseName := range releases {
			if idx < obsoleteN {
				releasePath := path.Join(workspaceDir, releaseName)
				if err := os.RemoveAll(releasePath); err != nil {
					return fmt.Errorf("failed to delete release %s: %v", releasePath, err)
				}
			}
		}
	}

	return nil
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
		os.Remove(link)
	}
	// TODO: what if the following fails? We are stuck with no `current` link
	return os.Symlink(target, link)
}

func decompressZip(zipFile, targetDir string) error {
	// Open the zip archive for reading
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// Iterate through each file in the archive
	for _, file := range r.File {
		// Open the file inside the zip archive
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		// Create the corresponding file in the target directory
		targetFilePath := filepath.Join(targetDir, file.Name)
		if file.FileInfo().IsDir() {
			// Create directories if file is a directory
			os.MkdirAll(targetFilePath, file.Mode())
		} else {
			// Create the file if it doesn't exist
			targetFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			if err != nil {
				return err
			}
			defer targetFile.Close()

			// Copy contents from the file inside the zip archive to the target file
			if _, err := io.Copy(targetFile, rc); err != nil {
				return err
			}
		}
	}

	return nil
}
