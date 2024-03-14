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

func decompressArchive(archivePath, targetDir string, uid, gid int) error {
	if strings.HasSuffix(archivePath, ".zip") {
		return decompressZip(archivePath, targetDir, uid, gid)
	} else if strings.HasSuffix(archivePath, ".tar.gz") {
		return decompressTarGzip(archivePath, targetDir, uid, gid)
	} else {
		return errors.New("unsupported archive type (supported types: zip, tar.gz)")
	}
}

func decompressZip(zipFile, targetDir string, uid, gid int) error {
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

		// Create the corresponding file in the target directory
		targetFilePath := filepath.Join(targetDir, file.Name)
		if file.FileInfo().IsDir() {
			// Create directories if file is a directory
			if err := createDirectory(targetFilePath, file.Mode(), uid, gid); err != nil {
				// close file
				rc.Close()
				return err
			}
		} else {
			// Create the file if it doesn't exist
			targetFile, err := createFile(targetFilePath, file.Mode(), uid, gid)
			if err != nil {
				rc.Close()
				return err
			}

			// Copy contents from the file inside the zip archive to the target file
			_, err = io.Copy(targetFile, rc)
			// close files
			targetFile.Close()
			if err != nil {
				rc.Close()
				return err
			}
		}

		// close the file in the archive
		rc.Close()
	}

	return nil
}

func decompressTarGzip(gzipFile, targetDir string, uid, gid int) error {
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
			if err := createDirectory(filePath, os.FileMode(header.Mode), uid, gid); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", filePath, err)
			}
		case tar.TypeReg:
			outFile, err := createFile(filePath, os.FileMode(header.Mode), uid, gid)
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

func createFile(path string, mode os.FileMode, uid, gid int) (f *os.File, err error) {
	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return
	}
	if err = os.Chown(path, uid, gid); err != nil {
		f.Close()
		f = nil
		return
	}
	return
}

func createDirectory(path string, mode os.FileMode, uid, gid int) error {
	if err := os.MkdirAll(path, mode); err != nil {
		return err
	}
	return os.Chown(path, uid, gid)
}
