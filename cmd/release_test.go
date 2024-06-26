package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kkentzo/rv/release"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Release_ShouldCreateResources_FromScratch(t *testing.T) {
	// do not create the workspace -- just specify the path
	workspacePath := uuid.NewString()
	// ensure that we'll clean up
	defer os.RemoveAll(workspacePath)

	// create and execute release
	out, err := createRelease(workspacePath, "foo.txt", 1)
	releaseId := parseReleaseFromOutput(out)
	require.NoError(t, err)

	// the workspace should now be present
	assert.DirExists(t, workspacePath)
	// the workspace should contain the release (extracted bundle under a versioned directory)
	assert.DirExists(t, path.Join(workspacePath, releaseId))
	// the release should have the correct contents
	assert.FileExists(t, path.Join(workspacePath, releaseId, "foo.txt"))
	// the workspace should contain the "current" release link
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName))
	// the "current" release link should point to the release folder
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "foo.txt"))
}

func Test_Release_ShouldCleanUp_WhenBundleDoesNotExist(t *testing.T) {
	workspacePath := uuid.NewString()
	// ensure that we'll clean up
	defer os.RemoveAll(workspacePath)

	// DO NOT create this bundle
	bundlePath := fmt.Sprintf("%s.zip", uuid.NewString())

	// prepare and execute command
	cmd := New()
	cmd.SetArgs([]string{"release", "-w", workspacePath, "-a", bundlePath})
	out := createOutputBuffer(cmd)
	// FIRE!
	require.NoError(t, cmd.Execute())

	// check that the output contains the relevant message
	assert.Contains(t, out.String(), "no such file or directory")

	// the workspace will not be cleared up
	assert.DirExists(t, workspacePath)
	// the workspace directory should be empty
	entries, err := ioutil.ReadDir(workspacePath)
	assert.NoError(t, err)
	assert.Empty(t, entries)
	// the workspace should NOT contain the "current" release link
	assert.NoFileExists(t, path.Join(workspacePath, release.CurrentLinkName))
}

func Test_Release_ShouldUpdateCurrent_WhenPreviousReleaseExists(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	// === create the first release ===
	// create and execute release
	out, err := createRelease(workspacePath, "foo.txt", 2)
	releaseId := parseReleaseFromOutput(out)

	require.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId, "foo.txt"))
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "foo.txt"))

	time.Sleep(10 * time.Millisecond)

	// === create the second release ===
	out, err = createRelease(workspacePath, "bar.txt", 2)
	releaseId = parseReleaseFromOutput(out)
	require.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId, "bar.txt"))
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "bar.txt"))
}

func Test_Release_ShouldKeepNMostRecentReleases(t *testing.T) {
	workspacePath := uuid.NewString()
	defer os.RemoveAll(workspacePath)

	// === create the first release ===
	// create and execute release
	out, err := createRelease(workspacePath, "foo.txt", 1)
	releaseId1 := parseReleaseFromOutput(out)

	require.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId1, "foo.txt"))
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "foo.txt"))

	time.Sleep(10 * time.Millisecond)

	// === create the second release ===
	out, err = createRelease(workspacePath, "bar.txt", 1)
	releaseId2 := parseReleaseFromOutput(out)
	require.NoError(t, err)
	assert.FileExists(t, path.Join(workspacePath, releaseId2, "bar.txt"))
	assert.FileExists(t, path.Join(workspacePath, release.CurrentLinkName, "bar.txt"))

	// the first release should be gone now
	assert.NoDirExists(t, path.Join(workspacePath, releaseId1))
}

func Test_Release_ShouldNotAcceptKeepZeroReleases(t *testing.T) {
	cmd := New()
	createOutputBuffer(cmd)
	cmd.SetArgs([]string{"release", "-w", "workspace", "-k", "0", "-a", "foo.zip"})
	// FIRE!
	err := cmd.Execute()
	assert.ErrorContains(t, err, "zero is not a valid value for --keep (-k) flag")
}
