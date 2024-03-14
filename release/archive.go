package release

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func decompressArchive(archivePath, targetDir string) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return decompressZip(archivePath, targetDir)
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		return decompressTarGzip(archivePath, targetDir)
	} else {
		return errors.New("unsupported archive type (supported types: zip, tar.gz)")
	}
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
			if err := os.MkdirAll(targetFilePath, file.Mode()); err != nil {
				return err
			}
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

func decompressTarGzip(gzipFile, targetDir string) error {
	stream, err := os.Open(gzipFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer stream.Close()
	uncompressedStream, err := gzip.NewReader(stream)
	if err != nil {
		return fmt.Errorf("failed to read archive: %v", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("failed to extract file from archive: %v", err)
		}

		filePath := path.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(filePath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", filePath, err)
			}
		case tar.TypeReg:
			outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", filePath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to create file %s: %v", filePath, err)
			}
			outFile.Close()

		default:
			return fmt.Errorf("unsupported file type: file=%s type=%d", header.Name, header.Typeflag)
		}

	}

	return nil
}
