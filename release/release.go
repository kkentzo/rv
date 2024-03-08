package release

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
)

func Install(workspaceDir, bundlePath string) (string, error) {
	// create release under workspace
	id := time.Now().Format("20060102150405.000")
	releaseDir := path.Join(workspaceDir, id)
	if err := os.MkdirAll(releaseDir, 0755); err != nil {
		return id, fmt.Errorf("failed to create release: %v", err)
	}
	// decompress bundle file
	if err := decompressZip(bundlePath, releaseDir); err != nil {
		// cleanup release directory
		defer os.RemoveAll(releaseDir)
		return id, fmt.Errorf("failed to decompress archive: %v", err)
	}
	// update current link
	if err := createOrUpdateLink(workspaceDir, id); err != nil {
		// cleanup release directory
		defer os.RemoveAll(releaseDir)
		return id, fmt.Errorf("failed to create/update link: %v", err)
	}
	return id, nil
}

func createOrUpdateLink(workspaceDir, target string) error {
	link := path.Join(workspaceDir, "current")
	// does the link already exist?
	_, err := os.Stat(link)
	if !os.IsNotExist(err) {
		os.Remove(link)
	}
	return os.Symlink(target, link)
}

func decompressZip(zipFile, targetDir string) error {
	// Open the zip archive for reading
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create the target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

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
