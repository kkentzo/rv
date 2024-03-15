package cmd

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kkentzo/rv/release"
	"github.com/spf13/cobra"
)

// ================
// HELPER FUNCTIONS
// ================
func createRelease(workspacePath, includedFile string, keepN uint) (string, error) {
	// create bundle
	bundlePath := fmt.Sprintf("%s.zip", uuid.NewString())
	defer deleteBundle(bundlePath)
	if err := createBundle(bundlePath, includedFile); err != nil {
		return "", err
	}

	// prepare command
	cmdOutput := createOutputBuffer(RootCmd)
	RootCmd.SetArgs([]string{"release", "-w", workspacePath, "-k", fmt.Sprint(keepN), "-a", bundlePath})
	// FIRE!
	if err := RootCmd.Execute(); err != nil {
		return "", err
	}

	return cmdOutput.String(), nil
}

func createReleases(workspacePath string, n uint) ([]string, error) {
	releases := []string{}
	for i := 0; i < int(n); i++ {
		out, err := createRelease(workspacePath, "foo", n)
		if err != nil {
			return releases, err
		}
		releases = append(releases, parseReleaseFromOutput(out))
		time.Sleep(10 * time.Millisecond)
	}
	return releases, nil
}

func rewindRelease(workspacePath, target string) (string, error) {
	// prepare command
	cmdOutput := createOutputBuffer(RootCmd)
	args := []string{"rewind", "-w", workspacePath, "-t", target}
	RootCmd.SetArgs(args)
	// FIRE!
	if err := RootCmd.Execute(); err != nil {
		return "", err
	}
	return cmdOutput.String(), nil
}

func parseReleaseFromOutput(out string) string {
	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "[success]") {
			return release.ReleaseFormatRe.FindString(line)
		}
	}
	return ""
}

func listReleases(workspacePath string) (string, error) {
	// prepare command
	cmdOutput := createOutputBuffer(RootCmd)
	RootCmd.SetArgs([]string{"list", "-w", workspacePath})
	// FIRE!
	if err := RootCmd.Execute(); err != nil {
		return "", err
	}

	return cmdOutput.String(), nil
}

func createOutputBuffer(cmd *cobra.Command) *bytes.Buffer {
	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	return actual
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func createBundle(zipFileName, includedFileName string) error {
	outFile, err := os.Create(zipFileName)
	if err != nil {
		return err
	}

	w := zip.NewWriter(outFile)

	if _, err := w.Create(includedFileName); err != nil {
		return fmt.Errorf("failed to include %s to zip file", zipFileName)
	}

	if err := w.Close(); err != nil {
		_ = outFile.Close()
		return errors.New("Warning: closing zipfile writer failed: " + err.Error())
	}

	if err := outFile.Close(); err != nil {
		return errors.New("Warning: closing zipfile failed: " + err.Error())
	}

	return nil
}

func deleteBundle(fname string) error {
	return os.Remove(fname)
}
